package api

type TMDbMovie struct {
	Adult         bool    `json:"adult"`
	BackdropPath  string  `json:"backdrop_path"`
	Id            int     `json:"id"`
	OriginalLang  string  `json:"original_language"`
	OriginalTitle string  `json:"original_title"`
	Overview      string  `json:"overview"`
	Popularity    float64 `json:"popularity"`
	PosterPath    string  `json:"poster_path"`
	ReleaseDate   string  `json:"release_date"`
	Title         string  `json:"title"`
	AvgVote       float32 `json:"vote_average"`
	VoteCount     int     `json:"vote_count"`
}

type TMDbMovieResp struct {
	Results []TMDbMovie `json:"results"`
}
