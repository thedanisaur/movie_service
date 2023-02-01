package types

type Rating struct {
	Series         string  `json:"series_title"`
	ChosenBy       string  `json:"series_chosen_by"`
	MoviesInSeries int     `json:"movies_in_series"`
	Good           string  `json:"good_votes"`
	Bad            string  `json:"bad_votes"`
	Total          string  `json:"total_votes"`
	Rating         float64 `json:"series_rating"`
}
