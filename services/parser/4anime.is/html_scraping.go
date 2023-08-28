package fouranimeisparsing

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Info struct {
	Poster string
	Title  string
}

func getDataFromHtml(html *string) (*Info, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(*html))
	if err != nil {
		return nil, err
	}

	var info Info

	document.Find("h1.anime_name").Each(func(index int, element *goquery.Selection) {
		animeName := element.Text()
		info.Title = animeName
	})

	image := document.Find(".anime_poster-img")
	imgSrc, _ := image.Attr("src")

	info.Poster = imgSrc

	return &info, nil
}
