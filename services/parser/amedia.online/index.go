package amediaonline_parsing

import (
	"anime-bot-schedule/pkg/fetch"
)

type AnimediaOnlineResp struct {
	AddedEpisode string
	Poster       string
	Title        string
}

func Fetch(url string) (*AnimediaOnlineResp, error) {
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
