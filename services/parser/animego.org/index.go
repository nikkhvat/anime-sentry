package animegoorgparsing

import (
	"anime-bot-schedule/pkg/fetch"
)

type Episod struct {
	Title   string `json:"title"`
	Date    string `json:"date"`
	Relized bool   `json:"relized"`
	Number  string `json:"number"`
}

type AnimeGoResp struct {
	Episods []Episod
	Image   *string
	Title   *string
}

func Fetch(url string) (*AnimeGoResp, error) {
	body, err := fetch.GET(url)

	if err != nil {
		return nil, err
	}

	data, err := getDataFromHtml(*body)

	if err != nil {
		return nil, err
	}

	return data, nil
}
