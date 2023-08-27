package animegoorgparsing

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type AnimeVostResponse struct {
	AddedEpisode    string
	NextEpisode     string
	NextEpisodeDate string
	Poster          string
	Title           string
}

func Fetch(url string) (*AnimeVostResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var result AnimeVostResponse

	result.Title, _ = document.Find("div[style*='width: 240px;'] img").Attr("title")
	result.Poster, _ = document.Find("meta[property='og:image']").Attr("content")
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

				result.AddedEpisode = strconv.Itoa(addedEpisode)
			}
		}

		parts = strings.Split(secondBracketContent, " серия - ")
		if len(parts) == 2 {
			result.NextEpisode = parts[0]

			date := strings.Split(parts[1], "]")
			result.NextEpisodeDate = date[0]
		}
	} else {
		fmt.Println("Could not extract episode info.")
	}

	result.AddedEpisode = result.AddedEpisode + " серия"
	result.NextEpisode = result.NextEpisode + " серия"
	result.Poster = "https://animevost.org" + result.Poster

	return &result, nil
}
