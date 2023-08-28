package amediaonlineparsing

import (
	"anime-bot-schedule/pkg/fetch"
)

type AnimeGoResp struct {
	AddedEpisode    string
	TotalEpisodes   string
	NextEpisode     string
	NextEpisodeDate string
	Poster          string
	Title           string
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
