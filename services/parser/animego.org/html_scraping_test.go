package animegoorg_parsing

import (
	"anime-bot-schedule/pkg/file"
	"testing"
)

func TestParseHtml(t *testing.T) {
	raw, err := file.Read("../../../tests/dump/html/animego.org")

	if err != nil {
		t.Fatalf("error read dump file")
	}

	data, err := getDataFromHtml(raw)

	if err != nil {
		t.Fatalf("error parse")
	}

	if *data.Title != "Реинкарнация безработного: История о приключениях в другом мире 2" {
		t.Fatalf("title doesn't match")
	}

	if *data.Image != "https://animego.org/media/cache/thumbs_500x700/upload/anime/images/64d4b199e9863172700580.jpg" {
		t.Fatalf("title doesn't match")
	}

	if len(data.Episods) != 3 {
		t.Fatal("episods lens doesn't match")
	}

	if data.Episods[0].Title != "TBA" && data.Episods[1].Title != "TBA" && data.Episods[2].Title != "TBA" {
		t.Fatalf("episod title doesn't match")
	}

	if data.Episods[0].Date != "10 сентября 2023" && data.Episods[1].Date != "3 сентября 2023" && data.Episods[2].Date != "27 августа 2023" {
		t.Fatalf("episod date doesn't match")
	}

	if data.Episods[0].Number != "10 серия" && data.Episods[1].Number != "9 серия" && data.Episods[2].Number != "8 серия" {
		t.Fatalf("episod number doesn't match")
	}

	if data.Episods[0].Relized != false && data.Episods[1].Relized != false && data.Episods[2].Relized != true {
		t.Fatalf("episod relised doesn't match")
	}
}
