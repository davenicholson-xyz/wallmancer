package main

import (
	"fmt"
	"log"

	"github.com/davenicholson-xyz/wallmancer/config"
)

func main() {
	cfg, err := config.New("config.yml")
	if err != nil {
		log.Fatalln("Failed to laod config", err)
	}

	username := cfg.GetString("wallhaven_username")
	fmt.Println(username)
}
