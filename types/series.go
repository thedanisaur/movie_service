package types

type Series struct {
	Order     string `json:"series_order"`
	Name      string `json:"series_name"`
	Title     string `json:"series_title"`
	CreatedOn string `json:"series_created_on"`
	PersonId  int    `json:"series_person_id"`
	ChosenBy  string `json:"series_chosen_by"`
}
