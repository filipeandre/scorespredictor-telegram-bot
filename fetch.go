package main

import (
	"github.com/gocolly/colly"
)

type htmlTable struct {
	Html string`json:"html"`
}

//Fetch pages
func collect(url string, sport string, date string) *colly.Collector{

	c := colly.NewCollector(colly.Async(true))
	err := c.Post(url, map[string]string{
		"current_sport": sport,
		"current_date": date,
		"force_get_json": "1",
		"bet_": "0",
	})

	if err != nil {
		stderr(err)
	}

	return c
}