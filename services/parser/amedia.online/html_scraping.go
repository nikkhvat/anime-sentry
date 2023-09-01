package amediaonline_parsing

import (
	"log"
	urlutil "net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getDataFromHtml(html string) (*AnimeGoResp, error) {

	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
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

	words := strings.Fields(result.Title)
	cleanedTitle := strings.Join(words, " ")

	result.Title = cleanedTitle

	return &result, nil
}
