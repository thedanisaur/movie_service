package handlers

import (
	"encoding/json"
	"fmt"
	"movie_service/db"
	"movie_service/types"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func TrackerRoute() string {
	return "/trackers"
}

func TrackerHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	txid := uuid.New()
	fmt.Printf("TrackerHandler | %s\n", txid.String())
	switch request.Method {
	case "GET":
		result := trackerGet()
		if result == nil {
			msg := fmt.Sprintf("%s %s failed: %s", request.Method, TrackerRoute(), txid.String())
			err := types.Error{Msg: msg}
			json.NewEncoder(writer).Encode(err)
		} else {
			json.NewEncoder(writer).Encode(result)
		}
	default:
		msg := fmt.Sprintf("%s %s unavailable: %s", request.Method, TrackerRoute(), txid.String())
		result := types.Error{Msg: msg}
		json.NewEncoder(writer).Encode(result)
	}
}

func trackerGet() []types.Tracker {
	fmt.Println("trackerGet")
	database := db.GetInstance()
	// Execute the query
	rows, err := database.Query("SELECT tracker_id, tracker_text, tracker_count, tracker_created_on, tracker_updated_on, tracker_created_by FROM trackers_vw")
	if err != nil {
		fmt.Printf("Failed to query databse\n%s\n", err.Error())
		return nil
	}

	var trackers []types.Tracker
	i := 0
	for rows.Next() {
		var tracker types.Tracker
		err = rows.Scan(&tracker.ID,
			&tracker.Text,
			&tracker.Count,
			&tracker.CreatedOn,
			&tracker.UpdatedOn,
			&tracker.CreatedBy)
		if err != nil {
			fmt.Printf("Failed to scan row\n%s\n", err.Error())
			return nil
		}
		tracker.Rank = i
		trackers = append(trackers, tracker)
		i = i + 1
	}

	err = rows.Err()
	if err != nil {
		fmt.Printf("Failed after row scan\n%s\n", err.Error())
		return nil
	}

	return trackers
}
