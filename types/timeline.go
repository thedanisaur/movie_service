package types

type Timeline struct {
	SeriesOrder     int             `json:"series_order"`
	SeriesTitle     string          `json:"series_title"`
	SeriesRank      int             `json:"series_rank"`
	SeriesRating    string          `json:"series_rating"`
	SeriesGoodVotes string          `json:"series_good_votes"`
	SeriesBadVotes  string          `json:"series_bad_votes"`
	SeriesChosenBy  string          `json:"series_chosen_by"`
	SeriesCreatedOn string          `json:"series_created_on"`
	SeriesMovies    []TimelineMovie `json:"series_movies"`
}

type SeriesRating struct {
	SeriesName      string `json:"series_name"`
	SeriesOrder     int    `json:"series_order"`
	SeriesTitle     string `json:"series_title"`
	SeriesCreatedOn string `json:"series_created_on"`
	SeriesGoodVotes string `json:"series_good_votes"`
	SeriesBadVotes  string `json:"series_bad_votes"`
	SeriesRating    string `json:"series_rating"`
	SeriesChosenBy  string `json:"series_chosen_by"`
}

type TimelineMovie struct {
	SeriesName string `json:"series_name"`
	MovieTitle string `json:"movie_title"`
	DanVote    string `json:"dan_vote"`
	NickVote   string `json:"nick_vote"`
}
