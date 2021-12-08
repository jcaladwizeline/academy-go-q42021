package model

type Anime struct {
	AnimeID  int    `json:"anime_id"`
	Title    string `json:"title"`
	Synopsis string `json:"synopsis"`
	Studio   string `json:"studio"`
}
