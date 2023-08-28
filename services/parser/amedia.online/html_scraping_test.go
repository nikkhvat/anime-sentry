package amediaonlineparsing

import (
	"anime-bot-schedule/pkg/file"
	"testing"
)

func TestParseHtml(t *testing.T) {
	raw, err := file.Read("../../../tests/dump/html/amedia.online")

	if err != nil {
		t.Fatalf("error read dump file")
	}

	data, err := getDataFromHtml(raw)

	if err != nil {
		t.Fatalf("error parse")
	}

	if data.Title != "Mushoku Tensei II: Isekai Ittara Honki Dasu / Реинкарнация безработного: История о приключениях в другом мире 2" {
		t.Fatalf("title doesn't match")
	}
}
