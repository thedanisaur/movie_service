package handlers

import (
	"encoding/json"
	"fmt"
	"movie_service/db"
	"movie_service/types"
	"net/http"
	"sort"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func TimelineRoute() string {
	return "/timeline"
}

func TimelineHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	txid := uuid.New()
	fmt.Printf("TimelineHandler | %s\n", txid.String())
	switch request.Method {
	case "GET":
		result := timelineGet()
		if result == nil {
			msg := fmt.Sprintf("%s %s failed: %s", request.Method, TimelineRoute(), txid.String())
			err := types.Error{Msg: msg}
			json.NewEncoder(writer).Encode(err)
		} else {
			json.NewEncoder(writer).Encode(result)
		}
	default:
		msg := fmt.Sprintf("%s %s unavailable: %s", request.Method, TimelineRoute(), txid.String())
		result := types.Error{Msg: msg}
		json.NewEncoder(writer).Encode(result)
	}
}

func timelineGet() []types.Timeline {
	fmt.Println("timelineGet")
	database := db.GetInstance()
	// Execute the query
	series_rating_rows, err := database.Query("SELECT series_name, series_order, series_title, series_created_on, good_votes, bad_votes, rating, chosen_by FROM rating_vw r_vw")
	if err != nil {
		fmt.Printf("Failed to query rating_vw\n%s\n", err.Error())
		return nil
	}

	var series_rating []types.SeriesRating
	for series_rating_rows.Next() {
		var rating types.SeriesRating
		err = series_rating_rows.Scan(&rating.SeriesName,
			&rating.SeriesOrder,
			&rating.SeriesTitle,
			&rating.SeriesCreatedOn,
			&rating.SeriesGoodVotes,
			&rating.SeriesBadVotes,
			&rating.SeriesRating,
			&rating.SeriesChosenBy)
		if err != nil {
			fmt.Printf("Failed to scan series_rating_rows\n%s\n", err.Error())
			return nil
		}
		series_rating = append(series_rating, rating)
	}

	err = series_rating_rows.Err()
	if err != nil {
		fmt.Printf("Failed after series_rating_rows scan\n%s\n", err.Error())
		return nil
	}

	// Get all the movies
	movie_votes_rows, err := database.Query("SELECT series_name, movie_title, dan_vote, nick_vote FROM dn_movies_votes_vw")
	if err != nil {
		fmt.Printf("Failed to query dn_movies_votes_vw\n%s\n", err.Error())
		return nil
	}

	var movies []types.Movie
	for movie_votes_rows.Next() {
		var movie types.Movie
		err = movie_votes_rows.Scan(&movie.SeriesName,
			&movie.MovieTitle,
			&movie.DanVote,
			&movie.NickVote)
		if err != nil {
			fmt.Printf("Failed to scan movie_votes_rows\n%s\n", err.Error())
			return nil
		}
		movies = append(movies, movie)
	}

	err = movie_votes_rows.Err()
	if err != nil {
		fmt.Printf("Failed after movie_votes_rows scan\n%s\n", err.Error())
		return nil
	}

	var timeline []types.Timeline
	for i := 0; i < len(series_rating); i++ {
		var series_movies []types.Movie
		for j := 0; j < len(movies); j++ {
			if movies[j].SeriesName == series_rating[i].SeriesName {
				series_movies = append(series_movies, movies[j])
			}
		}
		timeline = append(timeline, types.Timeline{
			SeriesOrder:     series_rating[i].SeriesOrder,
			SeriesTitle:     series_rating[i].SeriesTitle,
			SeriesRank:      i,
			SeriesGoodVotes: series_rating[i].SeriesGoodVotes,
			SeriesBadVotes:  series_rating[i].SeriesBadVotes,
			SeriesRating:    series_rating[i].SeriesRating,
			SeriesChosenBy:  series_rating[i].SeriesChosenBy,
			SeriesCreatedOn: series_rating[i].SeriesCreatedOn,
			SeriesMovies:    series_movies})
	}

	sort.Slice(timeline, func(i, j int) bool {
		return timeline[i].SeriesOrder > timeline[j].SeriesOrder
	})

	return timeline
}
