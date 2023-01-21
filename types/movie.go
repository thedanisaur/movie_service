package types

type Movie struct {
	Name       string `json:"movie_name"`
	SeriesName string `json:"series_name"`
	Title      string `json:"movie_title"`
	CreatedOn  string `json:"movie_created_on"`
}
