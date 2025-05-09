package main

import (
	"fmt"
	"log"
	"os"

	"github.com/davenicholson-xyz/wallmancer/config"
	"github.com/davenicholson-xyz/wallmancer/providers"
)

func main() {
	result, err := runApp()
	if err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
	fmt.Println(result)
}

func runApp() (string, error) {
	flg := config.NewFlagSet()

	flg.DefineString("provider", "", "wallpaper provider")
	flg.DefineString("username", "", "wallhaven.cc username")
	flg.DefineString("apikey", "", "wallhaven.cc api key")
	flg.DefineBool("nsfw", false, "Fetch NSFW images")

	flg.DefineString("random", "", "query for random wallpaper")
	flg.DefineBool("hot", false, "hot")
	flg.DefineBool("top", false, "toplist")
	flg.DefineString("seed", "", "random seed for search")

	flgValues := flg.Collect()

	cfg, err := config.New("config.yml")
	cfg.FlagOverride(flgValues)

	if err != nil {
		return "", fmt.Errorf("Failed to load config: %w", err)
	}

	prov := cfg.GetStringWithDefault("provider", "wallhaven")
	provider, exists := providers.GetProvider(prov)
	if !exists {
		return "", fmt.Errorf("Provider error: %w", err)
	}

	result, err := provider.ParseArgs(cfg)
	if err != nil {
		return "", fmt.Errorf("Unable to process args: %w", err)
	}

	return result, nil
}
