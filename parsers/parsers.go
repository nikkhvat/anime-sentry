package parsers

type Episode struct {
	Title    string `json:"title"`
	Date     string `json:"date"`
	Released bool   `json:"released"`
	Number   string `json:"number"`
}

type Dubbing struct {
	Title string `json:"title"`
}

type AnimeResponse struct {
	Episodes []Episode
	Dubbings []string
	Image    *string
	Title    *string
}
