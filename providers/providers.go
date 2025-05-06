package providers

import (
	"github.com/davenicholson-xyz/wallmancer/config"
)

type Provider interface {
	Name() string
	BuildURL(cfg *config.Config) (string, error)
}
