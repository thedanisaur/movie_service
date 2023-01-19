package types

type Tracker struct {
	ID []byte `json:"tracker_id"`
	Text string `json:"tracker_text"`
	Count int `json:"tracker_count"`
	// rank is just the popularity ordering from the sql view for
	// simpler sorting on the front end
	Rank int `json:"tracker_rank"`
	CreatedOn string `json:"tracker_created_on"`
	UpdatedOn string `json:"tracker_updated_on"`
	CreatedBy string `json:"tracker_created_by"`
}