package main

import (
	"fmt"
	"log"

	"github.com/davenicholson-xyz/wallmancer/config"
	"github.com/davenicholson-xyz/wallmancer/providers"
)

func main() {
	flg := config.NewFlagSet()

	flg.DefineString("provider", "", "wallpaper provider")
	flg.DefineString("wh_username", "", "wallhaven.cc username")
	flg.DefineString("wh_apikey", "", "wallhaven.cc api key")
	flg.DefineString("random", "", "query for random wallpaper")
	flg.DefineBool("hot", false, "wallhaven.cc hot list")
	flg.DefineBool("nsfw", false, "Fetch NSFW images")

	flgValues := flg.Collect()

	cfg, err := config.New("config.yml")
	cfg.FlagOverride(flgValues)

	if err != nil {
		log.Fatalln("Failed to load config", err)
	}

	prov := cfg.GetStringWithDefault("provider", "wallhaven")
	provider, exists := providers.GetProvider(prov)
	if !exists {
		log.Fatalf("issue with provider: %v", err)
	}

	url, err := provider.BuildURL(cfg)
	fmt.Println(url)
}
