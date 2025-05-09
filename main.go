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

	fmt.Printf("%+v\n", cfg)

	if err != nil {
		log.Fatalln("Failed to load config", err)
	}

	prov := cfg.GetStringWithDefault("provider", "wallhaven")
	provider, exists := providers.GetProvider(prov)
	if !exists {
		log.Fatalf("issue with provider: %v", err)
	}

	_, err = provider.ParseArgs(cfg)
	if err != nil {
		fmt.Printf("unable to run args: %s", err)
	}
}
