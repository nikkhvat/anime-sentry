package parsing

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Episod struct {
	Title   string `json:"title"`
	Date    string `json:"date"`
	Relized bool   `json:"relized"`
	Number  string `json:"number"`
}

type AnimeGoResp struct {
	Episods []Episod
	Image   *string
	Title   *string
}

func AnimeGOFetch(url string) (*AnimeGoResp, error) {
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
