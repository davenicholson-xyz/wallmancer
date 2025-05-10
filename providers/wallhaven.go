package providers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/davenicholson-xyz/wallmancer/appcontext"
	"github.com/davenicholson-xyz/wallmancer/download"
	"github.com/davenicholson-xyz/wallmancer/files"
)

type WallhavenProvider struct{}

type Wallpaper struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

type WallhavenData struct {
	Wallpapers []Wallpaper   `json:"data"`
	Meta       WallhavenMeta `json:"meta"`
}

type WallhavenMeta struct {
	LastPage int `json:"last_page"`
	Total    int `json:"total"`
}

func (w *WallhavenProvider) Name() string {
	return "wallhaven"
}

func (w *WallhavenProvider) ParseArgs(app *appcontext.AppContext) (string, error) {

	if app.Config.GetString("random") != "" || app.Config.GetBool("top") || app.Config.GetBool("hot") {
		wp, err := w.fetchRandom(app)
		if err != nil {
			return "", err
		}
		return wp, nil
	}

	return "", nil
}

func (w *WallhavenProvider) fetchRandom(app *appcontext.AppContext) (string, error) {
	url := download.NewURL("https://wallhaven.cc/api/v1/search")
	app.AddURLBuilder(url)

	lm := download.NewLinkManager()
	app.AddLinkManager(lm)

	var outfile string
	var selected string

	seed := app.Config.GetStringWithDefault("seed", download.GenerateSeed(6))
	app.URLBuilder.AddString("seed", seed)

	apikey := app.Config.GetString("apikey")
	if apikey != "" {
		app.URLBuilder.AddString("apikey", apikey)
	}

	url.SetString("purity", "100")
	if app.Config.GetBool("nsfw") {
		app.URLBuilder.SetString("purity", "111")
	}

	random := app.Config.GetString("random")
	if random != "" {
		app.URLBuilder.SetString("sorting", "random")
		app.URLBuilder.AddString("q", random)
		outfile = filepath.Join("wallhaven", "random")
	}

	if app.Config.GetBool("hot") {
		url.SetString("sorting", "hot")
		outfile = filepath.Join("wallhaven", "hot")
	}

	if app.Config.GetBool("top") {
		app.URLBuilder.SetString("sorting", "toplist")
		outfile = filepath.Join("wallhaven", "top")
	}

	selected, err := checkCacheForQuery(app, outfile)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	if selected != "" {
		output, err := files.ApplyWallpaper(selected, w.Name())
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}
		current_string := fmt.Sprintf("%s\n%s", selected, output)
		err = files.WriteStringToCache("wallhaven/current", current_string)
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}
		return selected, nil
	}

	selected, err = fetchQuery(app, outfile)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	if selected != "" {
		output, err := files.ApplyWallpaper(selected, w.Name())
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}
		current_string := fmt.Sprintf("%s\n%s", selected, output)
		err = files.WriteStringToCache("wallhaven/current", current_string)
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}
		return selected, nil
	}

	return "", nil

}

func processPage(app *appcontext.AppContext) (int, int, error) {
	request := app.URLBuilder.Build()

	resp, err := download.FetchJson(request)
	if err != nil {
		return 0, 0, fmt.Errorf("Could not fetch page: %w", err)
	}

	var wd WallhavenData
	if err := json.Unmarshal(resp, &wd); err != nil {
		return 0, 0, fmt.Errorf("Could not process JSON data: %w", err)
	}

	var links []string
	for _, link := range wd.Wallpapers {
		links = append(links, link.Path)
	}

	app.LinkManager.AddLinks(links)

	return wd.Meta.Total, wd.Meta.LastPage, nil
}

func checkCacheForQuery(app *appcontext.AppContext, outfile string) (string, error) {

	if outfile == "wallhaven/random" {
		last_query, err := app.CacheTools.ReadLineFromFile(outfile, 1)
		if err != nil {
			return "", nil
		}

		cleanUrl := app.URLBuilder.Without("apikey").Without("seed")
		query_url := cleanUrl.Build()

		if last_query == query_url {
			if files.IsFileFresh(app.CacheTools.Join(outfile), app.Config.GetIntWithDefault("expiry", 600)) {
				slog.Info("Using cached results")
				selected, err := files.GetRandomLine(app.CacheTools.Join(outfile))
				if err != nil {
					return "", fmt.Errorf("%w", err)
				}
				return selected, nil
			}
		} else {
			return "", nil
		}
	}

	if files.IsFileFresh(app.CacheTools.Join(outfile), app.Config.GetIntWithDefault("expiry", 600)) {
		selected, err := files.GetRandomLine(app.CacheTools.Join(outfile))
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}
		return selected, nil
	}

	return "", nil
}

func fetchQuery(app *appcontext.AppContext, outfile string) (string, error) {
	slog.Info("Using new query results")

	if outfile == "wallhaven/random" {
		cleanUrl := app.URLBuilder.Without("apikey").Without("seed")
		query_url := cleanUrl.Build()
		app.CacheTools.WriteStringToFile("wallhaven/last_query", query_url)
	}

	_, last, err := processPage(app)
	if err != nil {
		return "", fmt.Errorf("Unable to process page: %v -- %w", app.URLBuilder.Build(), err)
	}

	if app.LinkManager.Count() == 0 {
		files.WriteStringToCache(filepath.Join("wallhaven", "last_query"), "")
		return "", fmt.Errorf("No wallpapers found")
	}

	//TODO: Dont forget to make this concurrent
	if last > 1 {
		last_page := min(last, app.Config.GetIntWithDefault("max_pages", 5))
		for page := 2; page <= last_page; page++ {
			app.URLBuilder.SetInt("page", page)
			_, _, err = processPage(app)
			if err != nil {
				return "", fmt.Errorf("Unable to process page: %v -- %w", app.URLBuilder.Build(), err)
			}
		}
	}

	all_links := app.LinkManager.GetLinks()
	files.WriteSliceToCache(outfile, all_links)

	selected, err := files.GetRandomLine(app.CacheTools.Join(outfile))
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	return selected, nil
}
