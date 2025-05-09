package providers

import (
	"encoding/json"
	"fmt"

	"github.com/davenicholson-xyz/wallmancer/config"
	"github.com/davenicholson-xyz/wallmancer/download"
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

	seed := cfg.GetStringWithDefault("seed", download.GenerateSeed(6))
	url.AddString("seed", seed)

	apikey := cfg.GetString("apikey")
	if apikey != "" {
		url.AddString("apikey", apikey)
	}

	if cfg.GetBool("nsfw") {
		url.AddString("purity", "111")
	}

	random := cfg.GetString("random")
	if random != "" {
		url.AddString("sorting", "random")
		url.AddString("q", random)
	}

	if cfg.GetBool("hot") {
		url.AddString("sorting", "hot")
	}

	if cfg.GetBool("top") {
		url.AddString("sorting", "toplist")
	}

	// query_url := url.Build()

	// resp, err := download.FetchJson(query_url)
	// if err != nil {
	// 	return "", fmt.Errorf("%w", err)
	// }
	//
	// var wd WallhavenData
	// if err := json.Unmarshal(resp, &wd); err != nil {
	// 	return "", fmt.Errorf("%w", err)
	// }
	//
	// var links []string
	// for _, link := range wd.Wallpapers {
	// 	links = append(links, link.Path)
	// }
	links := processPage(url)
	lm.AddLinks(links)

	fmt.Println(len(links))

	return "", nil
}

func processPage(url *download.URLBuilder) []string {
	request := url.Build()
	fmt.Println(request)

	resp, err := download.FetchJson(request)
	if err != nil {
		fmt.Println(err)
	}

	var wd WallhavenData
	if err := json.Unmarshal(resp, &wd); err != nil {
		fmt.Println(err)
	}

	var links []string
	for _, link := range wd.Wallpapers {
		links = append(links, link.Path)
	}
	return links
}
