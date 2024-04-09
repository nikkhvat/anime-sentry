package amediaonline_parsing

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getDataFromHtml(html string) (*AnimediaOnlineResp, error) {

	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var result AnimediaOnlineResp
	src, _ := document.Find(".pmovie__img").Children().Eq(0).Attr("src")

	result.Title = document.Find("h1").Text()
	result.Poster = "https://amedia.site" + src

	document.Find(".seriianime").Each(func(i int, s *goquery.Selection) {
		text := s.Children().Eq(0).Text()
		textSplit := strings.Split(text, "-")[0]

		result.AddedEpisode = strings.TrimSpace(textSplit)
	})

	return &result, nil
}
