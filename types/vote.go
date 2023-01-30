package types

type Vote struct {
	ID             []byte `json:"vote_id"`
	Value          string `json:"vote_value"`
	CreatedOn      string `json:"vote_created_on"`
	MovieName      string `json:"movie_name"`
	PersonID       []byte `json:"person_id"`
	PersonUsername string `json:"person_username"`
}
