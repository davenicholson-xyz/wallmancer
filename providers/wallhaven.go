package providers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/davenicholson-xyz/wallmancer/config"
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

func (w *WallhavenProvider) ParseArgs(cfg *config.Config) (string, error) {

	if cfg.GetString("random") != "" || cfg.GetBool("top") || cfg.GetBool("hot") {
		wp, err := w.fetchRandom(cfg)
		if err != nil {
			return "", err
		}
		return wp, nil
	}

	return "", nil
}

func (w *WallhavenProvider) fetchRandom(cfg *config.Config) (string, error) {
	url := download.NewURL("https://wallhaven.cc/api/v1/search")
	lm := download.NewLinkManager()

	var outfile string
	var selected string

	seed := cfg.GetStringWithDefault("seed", download.GenerateSeed(6))
	url.AddString("seed", seed)

	apikey := cfg.GetString("apikey")
	if apikey != "" {
		url.AddString("apikey", apikey)
	}

	url.SetString("purity", "100")
	if cfg.GetBool("nsfw") {
		url.SetString("purity", "111")
	}

	random := cfg.GetString("random")
	if random != "" {
		url.SetString("sorting", "random")
		url.AddString("q", random)
		outfile = filepath.Join("wallhaven", "random")
	}

	if cfg.GetBool("hot") {
		url.SetString("sorting", "hot")
		outfile = filepath.Join("wallhaven", "hot")
	}

	if cfg.GetBool("top") {
		url.SetString("sorting", "toplist")
		outfile = filepath.Join("wallhaven", "top")
	}

	selected, err := checkCacheForQuery(cfg, outfile, url)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	if selected != "" {
		files.ApplyWallpaper(selected, w.Name())
		return selected, nil
	}

	selected, err = fetchQuery(url, lm, cfg, outfile)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	if selected != "" {
		// ApplyWallpaper should return the filepath it is saved to
		files.ApplyWallpaper(selected, w.Name())
		return selected, nil
	}

	return "", nil

}

func processPage(url *download.URLBuilder, lm *download.LinkManager) (int, int, error) {
	request := url.Build()

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

	lm.AddLinks(links)

	return wd.Meta.Total, wd.Meta.LastPage, nil
}

func checkCacheForQuery(cfg *config.Config, outfile string, url *download.URLBuilder) (string, error) {

	cache_dir, err := files.GetCacheDir()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	if outfile == "wallhaven/random" {
		last_query, err := files.ReadFromCache("wallhaven/last_query")
		if err != nil {
			return "", nil
		}

		cleanUrl := url.Without("apikey").Without("seed")
		query_url := cleanUrl.Build()

		if last_query == query_url {
			if files.IsFileFresh(filepath.Join(cache_dir, outfile), cfg.GetIntWithDefault("expiry", 600)) {
				slog.Info("Using cached results")
				selected, err := files.GetRandomLine(filepath.Join(cache_dir, outfile))
				if err != nil {
					return "", fmt.Errorf("%w", err)
				}
				return selected, nil
			}
		} else {
			return "", nil
		}
	}

	if files.IsFileFresh(filepath.Join(cache_dir, outfile), cfg.GetIntWithDefault("expiry", 600)) {
		selected, err := files.GetRandomLine(filepath.Join(cache_dir, outfile))
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}
		return selected, nil
	}

	return "", nil
}

func fetchQuery(url *download.URLBuilder, lm *download.LinkManager, cfg *config.Config, outfile string) (string, error) {
	slog.Info("Using new query results")
	cache_dir, err := files.GetCacheDir()

	if outfile == "wallhaven/random" {
		cleanUrl := url.Without("apikey").Without("seed")
		query_url := cleanUrl.Build()
		files.WriteStringToCache(filepath.Join("wallhaven", "last_query"), query_url)
	}

	_, last, err := processPage(url, lm)
	if err != nil {
		return "", fmt.Errorf("Unable to process page: %v -- %w", url.Build(), err)
	}

	if lm.Count() == 0 {
		files.WriteStringToCache(filepath.Join("wallhaven", "last_query"), "")
		return "", fmt.Errorf("No wallpapers found")
	}

	if last > 1 {
		last_page := min(last, cfg.GetIntWithDefault("max_pages", 5))
		for page := 2; page <= last_page; page++ {
			url.SetInt("page", page)
			_, _, err = processPage(url, lm)
			if err != nil {
				return "", fmt.Errorf("Unable to process page: %v -- %w", url.Build(), err)
			}
		}
	}

	all_links := lm.GetLinks()
	files.WriteSliceToCache(outfile, all_links)

	selected, err := files.GetRandomLine(filepath.Join(cache_dir, outfile))
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	return selected, nil
}
