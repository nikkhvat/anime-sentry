package animegoorg_parsing

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type IData struct {
	AddedEpisode    string
	NextEpisode     string
	NextEpisodeDate string
	Poster          string
	Title           string
}

func getDataFromHtml(html string) (*IData, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		return nil, err
	}

	var data IData

	data.Title, _ = document.Find("div[style*='width: 240px;'] img").Attr("title")
	data.Poster, _ = document.Find("meta[property='og:image']").Attr("content")
	episodeInfo := document.Find("div.shortstoryHead h1").Text()

	bracketsInfo := strings.Split(episodeInfo, "] [")
	if len(bracketsInfo) == 2 {
		firstBracketContent := strings.TrimSuffix(strings.TrimPrefix(bracketsInfo[0], "["), "]")
		secondBracketContent := strings.TrimSuffix(strings.TrimPrefix(bracketsInfo[1], "["), "]")

		parts := strings.Split(firstBracketContent, " из ")

		if len(parts) == 2 {
			episodeRange := strings.Split(parts[0], "-")

			if len(episodeRange) == 3 {
				addedEpisode, _ := strconv.Atoi(episodeRange[2])

				data.AddedEpisode = strconv.Itoa(addedEpisode)
			}
		}

		parts = strings.Split(secondBracketContent, " серия - ")
		if len(parts) == 2 {
			data.NextEpisode = parts[0]

			date := strings.Split(parts[1], "]")
			data.NextEpisodeDate = date[0]
		}
	} else {
		fmt.Println("Could not extract episode info.")
	}

	data.AddedEpisode = data.AddedEpisode + " серия"
	data.NextEpisode = data.NextEpisode + " серия"
	data.Poster = "https://animevost.org" + data.Poster

	log.Println(data)

	return &data, nil
}
