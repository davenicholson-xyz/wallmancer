package providers

import (
	"github.com/davenicholson-xyz/wallmancer/config"
	"github.com/davenicholson-xyz/wallmancer/download"
)

type WallhavenProvider struct{}

func (w *WallhavenProvider) Name() string {
	return "wallhaven"
}

func (w *WallhavenProvider) BuildURL(cfg *config.Config) (string, error) {
	url := download.NewURL("https://wallhaven.cc/api/v1/search")

	apikey := cfg.GetString("wh_apikey")
	if apikey != "" {
		url.AddString("apikey", apikey)
	}

	query_url := url.Build()
	return query_url, nil
}
