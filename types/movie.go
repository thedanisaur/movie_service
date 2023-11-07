package types

type Movie struct {
	SeriesName     string         `json:"series_name"`
	SeriesTitle    string         `json:"series_title"`
	MovieName      string         `json:"movie_name"`
	MovieTitle     string         `json:"movie_title"`
	MovieCreatedOn string         `json:"movie_created_on"`
	MovieImage     string         `json:"movie_image"`
	DanVote        string         `json:"dan_vote"`
	NickVote       string         `json:"nick_vote"`
	Trackers       []MovieTracker `json:"movie_trackers"`
}
