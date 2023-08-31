package animevosttvarsing

import (
	"anime-bot-schedule/pkg/file"
	"testing"
)

func TestParseHtml(t *testing.T) {
	raw, err := file.Read("../../../tests/dump/html/animevost.tv")

	if err != nil {
		t.Fatalf("error read dump file")
	}

	data, err := getDataFromHtml(raw)

	if err != nil {
		t.Fatalf("error parse")
	}

	if data.AddedEpisode != "8 серия" {
		t.Fatalf("added edpisod doesn't match")
	}

	if data.NextEpisode != "9 серия" {
		t.Fatalf("next edpisod doesn't match")
	}

	if data.NextEpisodeDate != "28 августа" {
		t.Fatalf("next edpisod date doesn't match")
	}

	if data.Poster != "https://animevost.org/uploads/posts/2023-06/1687522224_1.jpg" {
		t.Fatalf("poster doesn't match")
	}

	if data.Title != "Месть Масамунэ! (второй сезон) / Masamune-kun no Revenge R" {
		t.Fatalf("title doesn't match")
	}
}
