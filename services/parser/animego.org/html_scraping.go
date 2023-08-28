package animegoorgparsing

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getDataFromHtml(html string) (*AnimeGoResp, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var episods []Episod
	document.Find(".released-episodes-container .col-12 .released-episodes-item .row.m-0").Each(func(i int, s *goquery.Selection) {
		relized := s.Find("div:nth-child(4) span.cursor-pointer").Length() > 0

		number := strings.TrimSpace(s.Find("div:nth-child(1)").Text())
		title := strings.TrimSpace(s.Find("div:nth-child(2)").Text())
		date := strings.TrimSpace(s.Find("div:nth-child(3)").Text())

		episods = append(episods, Episod{
			Title:   title,
			Date:    date,
			Relized: relized,
			Number:  number,
		})
	})

	image := document.Find("img").Eq(2)

	srcset, _ := image.Attr("srcset")

	srcsetArray := strings.Split(srcset, " ")

	title, _ := image.Attr("title")

	return &AnimeGoResp{
		Episods: episods,
		Image:   &srcsetArray[0],
		Title:   &title,
	}, nil
}
