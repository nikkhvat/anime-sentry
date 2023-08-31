package main

import (
	animegotvparsing "anime-bot-schedule/services/parser/animevost.tv"
	"log"
)

func main() {
	data, err := animegotvparsing.Fetch("https://animevost.tv/mest-masamunje-sezon-2-2023-1080-hd")

	log.Println("======================================")
	log.Println("AddedEpisode: ( ", data.AddedEpisode, " )")
	log.Println("NextEpisode: ( ", data.NextEpisode, " )")
	log.Println("NextEpisodeDate: ( ", data.NextEpisodeDate, " )")
	log.Println("Poster: ( ", data.Poster, " )")
	log.Println("Title: ( ", data.Title, " )")
	log.Println("======================================")
	log.Println(err)
}
