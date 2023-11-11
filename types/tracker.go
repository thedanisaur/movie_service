package types

import "github.com/google/uuid"

type MovieTracker struct {
	MovieName    string    `json:"movie_name"`
	TrackerID    uuid.UUID `json:"tracker_id"`
	TrackerText  string    `json:"tracker_text"`
	TrackerCount int       `json:"tracker_count"`
}

type MovieTracker2 struct {
	MovieTitle   string `json:"movie_title"`
	TrackerCount int    `json:"tracker_count"`
}

type Tracker struct {
	ID    uuid.UUID `json:"tracker_id"`
	Text  string    `json:"tracker_text"`
	Count int       `json:"tracker_count"`
	// rank is just the popularity ordering from the sql view for
	// simpler sorting on the front end
	Rank      int    `json:"tracker_rank"`
	Image     string `json:"tracker_image"`
	CreatedOn string `json:"tracker_created_on"`
	UpdatedOn string `json:"tracker_updated_on"`
	CreatedBy string `json:"tracker_created_by"`
}
