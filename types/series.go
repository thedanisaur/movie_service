package types

import "github.com/google/uuid"

type Series struct {
	Order     string    `json:"series_order"`
	Name      string    `json:"series_name"`
	Title     string    `json:"series_title"`
	CreatedOn string    `json:"series_created_on"`
	PersonID  uuid.UUID `json:"series_person_id"`
	ChosenBy  string    `json:"series_chosen_by"`
}

type SeriesImage struct {
	SeriesName string `json:"series_name"`
	MovieName  string `json:"movie_name"`
	ImagePath  string `json:"image_path"`
}
