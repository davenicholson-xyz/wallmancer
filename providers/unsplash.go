package providers

import (
	"github.com/davenicholson-xyz/wallmancer/config"
	"github.com/davenicholson-xyz/wallmancer/download"
)

type UnsplashProvider struct{}

func (u *UnsplashProvider) Name() string {
	return "unsplash"
}

func (u *UnsplashProvider) BuildURL(cfg *config.Config) (string, error) {
	url := download.NewURL("https://api.unsplash.cc/search")
	query_url := url.Build()
	return query_url, nil
}
