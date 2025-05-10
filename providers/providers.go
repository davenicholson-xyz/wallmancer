package providers

import "github.com/davenicholson-xyz/wallmancer/appcontext"

type Provider interface {
	Name() string
	ParseArgs(app *appcontext.AppContext) (string, error)
}
