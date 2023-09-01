package fouranimeis_parsing

import (
	"anime-bot-schedule/pkg/fetch"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type AnimeGoResp struct {
	LastEpisode     string
	TotalEpisodes   int64
	Poster          string
	Title           string
	URL             string
	LastEpisodeLink string
}

type AnimeFetchResponse struct {
	Status     bool   `json:"status"`
	HTML       string `json:"html"`
	TotalItems int    `json:"totalItems"`
}

func getIdFromUrl(url string) string {
	parts := strings.Split(url, "-")
	lastPart := parts[len(parts)-1]
	idParts := strings.Split(lastPart, "?")
	return idParts[0]
}

type AnimeEpisod struct {
	EpisodNumber string `json:"episod"`
	EpisodLink   string `json:"link"`
}

func currentUrl(rawUrl string) (string, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	pathParts := strings.Split(parsedUrl.Path, "/")

	isEdited := false

	for i, part := range pathParts {

		if part == "watch" && i < len(pathParts)-1 {
			isEdited = true
			newPathParts := append(pathParts[:i], pathParts[i+2:]...)
			parsedUrl.Path = strings.Join(newPathParts, "/")

			query := parsedUrl.Query()
			query.Del("ep")
			parsedUrl.RawQuery = query.Encode()
			break
		}
	}

	if isEdited {
		return parsedUrl.String() + "/" + pathParts[2], nil
	}

	return parsedUrl.String(), nil
}

func getLastEpisod(id string) (*AnimeEpisod, error) {
	client := &http.Client{}

	var url2fetch = "https://4anime.is/ajax/episode/list/" + id

	req, err := http.NewRequest("GET", url2fetch, nil)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var animeResp AnimeFetchResponse
	err = json.NewDecoder(resp.Body).Decode(&animeResp)
	if err != nil {
		return nil, err
	}

	episodsDocument, err := goquery.NewDocumentFromReader(strings.NewReader(animeResp.HTML))

	if err != nil {
		fmt.Println("Error creating the document:", err)
		return nil, errors.New("Error creating the document")
	}

	var lastEpisodeNumber string
	var lastEpisodeLink string

	episodsDocument.Find("li.ep-item").Each(func(index int, element *goquery.Selection) {
		lastEpisodeNumber = element.Find("a").Text()

		link, exist := element.Find("a").Attr("href")

		if exist {
			lastEpisodeLink = link
		}
	})

	fmt.Println("Last episode number is:", lastEpisodeNumber)

	episod := AnimeEpisod{
		EpisodNumber: lastEpisodeNumber,
		EpisodLink:   lastEpisodeLink,
	}

	return &episod, nil
}

func Fetch(rawurl string) (*AnimeGoResp, error) {

	url, err := currentUrl(rawurl)

	if err != nil {
		return nil, err
	}

	id := getIdFromUrl(url)

	var data AnimeGoResp

	lastEpisod, err := getLastEpisod(id)

	data.LastEpisode = lastEpisod.EpisodNumber + " episod"
	data.LastEpisodeLink = "https://4anime.is" + lastEpisod.EpisodLink

	i, err := strconv.ParseInt(lastEpisod.EpisodNumber, 10, 64)
	if err == nil {
		data.TotalEpisodes = i
	}

	if err != nil {
		return nil, err
	}

	data.URL = url
	body, err := fetch.GET(url)

	info, err := getDataFromHtml(body)

	data.Poster = info.Poster
	data.Title = info.Title

	return &data, nil
}
