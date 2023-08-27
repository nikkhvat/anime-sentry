package main

import (
	animegoorgparsing "anime-bot-schedule/parsing/animevost.org"
	"log"
)

func main() {
	data, err := animegoorgparsing.Fetch("https://animevost.org/tip/tv/3003-masamune-kun-no-revenge-r.html")

	log.Println("======================================")
	log.Println("AddedEpisode: ( ", data.AddedEpisode, " )")
	log.Println("NextEpisode: ( ", data.NextEpisode, " )")
	log.Println("NextEpisodeDate: ( ", data.NextEpisodeDate, " )")
	log.Println("Poster: ( ", data.Poster, " )")
	log.Println("Title: ( ", data.Title, " )")
	log.Println("======================================")
	log.Println(err)
}
