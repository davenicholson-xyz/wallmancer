package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/davenicholson-xyz/go-cachetools/cachetools"
	"github.com/davenicholson-xyz/wallmancer/appcontext"
	"github.com/davenicholson-xyz/wallmancer/config"
	"github.com/davenicholson-xyz/wallmancer/files"
	"github.com/davenicholson-xyz/wallmancer/providers"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelInfo)
	result, err := runApp()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println(result)
}

func runApp() (string, error) {
	app := appcontext.NewAppContext()

	flg := config.NewFlagSet()

	flg.DefineString("provider", "", "wallpaper provider")
	flg.DefineString("username", "", "wallhaven.cc username")
	flg.DefineString("apikey", "", "wallhaven.cc api key")
	flg.DefineBool("nsfw", false, "Fetch NSFW images")
	flg.DefineInt("expiry", 0, "cache expiry in seconds")

	flg.DefineBool("clear", false, "clear the wallmancer cache")

	flg.DefineString("random", "", "query for random wallpaper")
	flg.DefineBool("hot", false, "hot")
	flg.DefineBool("top", false, "toplist")
	flg.DefineString("seed", "", "random seed for search")

	flgValues := flg.Collect()

	default_cfg_path, exists := files.DefaultConfigFilepath()
	cfg, err := config.New(default_cfg_path)
	cfg.FlagOverride(flgValues)

	app.AddConfig(cfg)

	flagstring := fmt.Sprintf("%v+", app.Config)
	slog.Info(flagstring)

	if err != nil {
		return "", fmt.Errorf("Failed to load config: %w", err)
	}

	ct, err := cachetools.New("wallmancer")
	if err != nil {
		return "", fmt.Errorf("Error creating cache: %w", err)
	}

	app.AddCacheTools(ct)

	if app.Config.GetBool("clear") {
		slog.Info("Clearing the cache")
		err := app.CacheTools.Clear()
		if err != nil {
			return "", fmt.Errorf("Error deleting cache: %w", err)
		}
		return "Cache deleted", nil
	}

	prov := app.Config.GetStringWithDefault("provider", "wallhaven")
	provider, exists := providers.GetProvider(prov)
	if !exists {
		return "", fmt.Errorf("Provider error: %w", err)
	}

	result, err := provider.ParseArgs(app)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return result, nil
}
