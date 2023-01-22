package types

type Movie struct {
	Name       string `json:"movie_name"`
	SeriesName string `json:"series_name"`
	Title      string `json:"movie_title"`
	CreatedOn  string `json:"movie_created_on"`
}

type Movie2 struct {
	SeriesName  string         `json:"series_name"`
	SeriesTitle string         `json:"series_title"`
	MovieName   string         `json:"movie_name"`
	MovieTitle  string         `json:"movie_title"`
	DanVote     string         `json:"dan_vote"`
	NickVote    string         `json:"nick_vote"`
	Trackers    []MovieTracker `json:"movie_trackers"`
}
