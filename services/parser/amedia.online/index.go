package amediaonlineparsing

import (
	"log"
	"net/http"
	urlutil "net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

	var result AnimeGoResp

	// Извлечение данных о сериях
	document.Find(".info").Each(func(i int, s *goquery.Selection) {
		added := s.Find("span").First().Text()
		total := s.Find("span").Eq(1).Text()
		title := strings.TrimSpace(s.Text())
		next := strings.Split(title, " ")[29] + " серия"

		nextDate, _ := s.Find("span > a").Attr("href")
		nextDate = strings.TrimPrefix(nextDate, "https://amedia.online/dat/")
		nextDate = strings.TrimSuffix(nextDate, "/")
		decodedNextDate, err := urlutil.QueryUnescape(nextDate)
		if err != nil {
			log.Println("Error decoding NextEpisodeDate:", err)
		} else {
			nextDate = decodedNextDate
		}

		result.AddedEpisode = added
		result.TotalEpisodes = total
		result.NextEpisode = next
		result.NextEpisodeDate = nextDate
	})

	// Извлечение постера
	document.Find(".film-poster img").Each(func(i int, s *goquery.Selection) {
		imgSrc, _ := s.Attr("data-src")
		result.Poster = "https://amedia.online" + imgSrc
	})

	// Извлечение названия
	document.Find(".titleor").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Text())
		result.Title = title
	})

	return &result, nil
}
