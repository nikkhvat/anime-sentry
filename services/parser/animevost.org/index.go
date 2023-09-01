package animegoorg_parsing

import (
	"anime-bot-schedule/pkg/fetch"
	"log"
)

func Fetch(url string) (*IData, error) {
	// * Fetch html
	body, err := fetch.GET(url)

	if err != nil {
		log.Panicln(err)
	}

	// * Parse html
	data, err := getDataFromHtml(*body)

	if err != nil {
		log.Panicln(err)
	}

	return data, nil
}
