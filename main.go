package main

import (
	"fmt"
	"log"

	"github.com/davenicholson-xyz/wallmancer/config"
)

func main() {
	flg := config.NewFlagSet()

	flg.DefineString("wallhaven_username", "", "wallhaven.cc username")
	flg.DefineString("wallhaven_apikey", "", "wallhaven.cc api key")
	flg.DefineBool("nsfw", false, "Fetch NSFW images")

	flgValues := flg.Collect()
	fmt.Println(flgValues)

	cfg, err := config.New("config.yml")
	cfg.FlagOverride(flgValues)

	if err != nil {
		log.Fatalln("Failed to load config", err)
	}

	nsfw := cfg.GetBoolWithDefault("nsfw", false)

	fmt.Println(nsfw)

}
