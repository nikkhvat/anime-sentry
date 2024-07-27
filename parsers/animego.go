package parsers

import (
	dateUtils "anime-sentry/pkg/date"
	"anime-sentry/pkg/fetch"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type GetDubbingsResponse struct {
	HTML string `json:"content"`
}

func getLastAndNextEpisode(episodes []Episode) (Episode, Episode) {
	var lastReleased Episode
	var nextEpisode Episode

	foundReleased := false

	for i := len(episodes) - 1; i >= 0; i-- {
		if episodes[i].Released && !foundReleased {
			lastReleased = episodes[i]
			foundReleased = true
		} else if foundReleased && !episodes[i].Released {
			nextEpisode = episodes[i]
			break
		}
	}

	return lastReleased, nextEpisode
}

func getDubbings(episode, id string) ([]string, error) {
	url := fmt.Sprintf("https://animego.org/anime/series?dubbing=&provider=&episode=%s&id=%s", episode, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("accept-language", "ru,en-GB;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Add("x-requested-with", "XMLHttpRequest")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var animeResp GetDubbingsResponse
	err = json.NewDecoder(resp.Body).Decode(&animeResp)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(animeResp.HTML))

	if err != nil {
		return nil, err
	}

	tabContent := doc.Find("#video-dubbing")

	var dubbings []string

	tabContent.Find(".video-player-toggle-item").Each(func(_ int, s *goquery.Selection) {
		dubbingName := s.Find(".video-player-toggle-item-name").Text()
		if name := strings.TrimSpace(dubbingName); name != "" {
			dubbings = append(dubbings, name)
		}
	})

	return dubbings, nil

}

func getDataFromHtmlAnimeGo(html string) (*AnimeResponse, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var episodes []Episode
	document.Find(".released-episodes-container .col-12 .released-episodes-item .row.m-0").Each(func(i int, s *goquery.Selection) {
		isReleased := s.Find("div:nth-child(4) span.cursor-pointer").Length() > 0

		number := strings.TrimSpace(s.Find("div:nth-child(1)").Text())
		title := strings.TrimSpace(s.Find("div:nth-child(2)").Text())
		date := strings.TrimSpace(s.Find("div:nth-child(3)").Text())

		numberArray := strings.Split(number, " ")

		parsedDate, _ := dateUtils.ConvertDate(date)

		episodes = append(episodes, Episode{
			Title:    title,
			Released: isReleased,
			// * For example: 19/07 (dd/mm)
			Date: parsedDate,
			// * For example: 19
			Number: numberArray[0],
		})
	})

	spanWithId := document.Find(".released-episodes-watch")

	lastReleased, _ := getLastAndNextEpisode(episodes)

	dataWatchedId, exists := spanWithId.Attr("data-watched-id")

	var dubbings []string

	if exists {
		dubbings, _ = getDubbings(lastReleased.Number, dataWatchedId)
	}

	image := document.Find("img")

	srcset, _ := image.Attr("srcset")

	srcsetArray := strings.Split(srcset, " ")

	title := document.Find("h1").Text()

	return &AnimeResponse{
		Episodes: episodes,
		Image:    &srcsetArray[0],
		Title:    &title,
		Dubbings: dubbings,
	}, nil
}

func FetchAnimeGo(url string) (*AnimeResponse, error) {
	body, err := fetch.GET(url)

	if err != nil {
		return nil, err
	}

	return getDataFromHtmlAnimeGo(*body)
}
