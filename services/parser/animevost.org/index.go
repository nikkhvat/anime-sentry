package animegoorg_parsing

import (
	"anime-bot-schedule/pkg/fetch"
	"log"
)

func Fetch(url string) (*IData, error) {
	body, err := fetch.GET(url)

	if err != nil {
		log.Panicln(err)
	}

	data, err := getDataFromHtml(*body)

	if err != nil {
		log.Panicln(err)
	}

	return data, nil
}
