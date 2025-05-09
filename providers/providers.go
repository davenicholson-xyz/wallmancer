package providers

import (
	"github.com/davenicholson-xyz/wallmancer/config"
)

type Provider interface {
	Name() string
	ParseArgs(cfg *config.Config) (string, error)
}
