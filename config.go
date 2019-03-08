package main

import (
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strings"
)

type telegram struct {
	Token string `json:"token"`
	Channel string `json:"channel"`
	Channel2 string `json:"channel2"`
}

type configuration struct {
	Url      string    `json:"url"`
	Sports   []string  `json:"sports"`
	Telegram *telegram `json:"telegram"`
	Temp     string    `json:"temp"`
	Filter   string    `json:"filter"`
}

//Load configs
func loadConfiguration() *configuration {

	path:= os.Getenv("SCORESPREDICTOR_HOME") + ".env.yaml"
	if _, err := os.Stat(path); err == nil {
		err := godotenv.Load(path)
		if err != nil {
			onError(err)
		}
	}

	configuration := &configuration{
		Url:os.Getenv("SCORESPREDICTOR_URL"),
		Telegram: &telegram {
			Token: os.Getenv("TELEGRAM_BOT_TOKEN"),
			Channel: os.Getenv("TELEGRAM_BOT_CHANNEL"),
			Channel2: os.Getenv("TELEGRAM_BOT_CHANNEL2"),
		},
		Sports: strings.Split(os.Getenv("SCORESPREDICTOR_SPORTS"), ";"),
		Temp:   os.TempDir(),
		Filter: os.Getenv("FILTER"),
	}

	return configuration
}
