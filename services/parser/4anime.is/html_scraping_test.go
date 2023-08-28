package fouranimeisparsing

import (
	"anime-bot-schedule/pkg/file"
	"testing"
)

func TestParseHtml(t *testing.T) {
	raw, err := file.Read("../../../tests/dump/html/4anime.is")

	if err != nil {
		t.Fatalf("error read dump file")
	}

	data, err := getDataFromHtml(&raw)

	if err != nil {
		t.Fatalf("error parse")
	}

	if data.Poster != "https://img.bunnyccdn.co/_r/300x400/100/17/4a/174a13011301d3d0a2135f40162b8dee/174a13011301d3d0a2135f40162b8dee.jpg" {
		t.Fatalf("image doesn't match")
	}

	if data.Title != "Masamune-kun's Revenge R" {
		t.Fatalf("title doesn't match")
	}
}
