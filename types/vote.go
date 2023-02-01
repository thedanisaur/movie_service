package types

import "github.com/google/uuid"

type Vote struct {
	ID             uuid.UUID `json:"vote_id"`
	Value          string    `json:"vote_value"`
	CreatedOn      string    `json:"vote_created_on"`
	MovieName      string    `json:"movie_name"`
	PersonID       uuid.UUID `json:"person_id"`
	PersonUsername string    `json:"person_username"`
}
