package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"movie_service/db"
	"movie_service/security"
	"movie_service/types"
	"movie_service/util"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetMovies(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(GetMovies), txid.String())
	err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
	database, err := db.GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return c.Status(fiber.StatusInternalServerError).SendString(err_string)
	}
	// Get all the movies
	movies_votes_query := `
		SELECT series_name
			, series_title
			, movie_name
			, movie_title
			, movie_created_on
			, movie_image
			, dan_vote
			, nick_vote
		FROM dn_movies_votes_vw
	`
	movie_votes_rows, err := database.Query(movies_votes_query)
	if err != nil {
		log.Printf("Failed to query dn_movies_votes_vw:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}

	var movies []types.Movie
	for movie_votes_rows.Next() {
		var movie types.Movie
		err = movie_votes_rows.Scan(&movie.SeriesName,
			&movie.SeriesTitle,
			&movie.MovieName,
			&movie.MovieTitle,
			&movie.MovieCreatedOn,
			&movie.MovieImage,
			&movie.DanVote,
			&movie.NickVote)
		if err != nil {
			log.Printf("Failed to scan movie_votes_rows:\n%s\n", err.Error())
			return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
		}
		movies = append(movies, movie)
	}

	err = movie_votes_rows.Err()
	if err != nil {
		log.Printf("Failed after movie_votes_rows scan:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}

	// Get top 3 trackers
	movie_trackers_query := `
		SELECT  mt.movie_name
			, BIN_TO_UUID(t.tracker_id) tracker_id
			, t.tracker_text
			, SUM(IFNULL(mt.tracker_count, 0)) tracker_count
		FROM    movie_trackers mt
			, trackers t
		WHERE   mt.tracker_id = t.tracker_id
		GROUP BY t.tracker_id
			, mt.movie_name
		ORDER BY  mt.movie_name
			, tracker_count DESC
	`
	movie_trackers_rows, err := database.Query(movie_trackers_query)
	if err != nil {
		log.Printf("Failed to query movies_trackers:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}

	var movie_trackers []types.MovieTracker
	for movie_trackers_rows.Next() {
		var movie_tracker types.MovieTracker
		err = movie_trackers_rows.Scan(&movie_tracker.MovieName,
			&movie_tracker.TrackerID,
			&movie_tracker.TrackerText,
			&movie_tracker.TrackerCount)
		if err != nil {
			log.Printf("Failed to scan movie_trackers_rows:\n%s\n", err.Error())
			return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
		}
		movie_trackers = append(movie_trackers, movie_tracker)
	}

	err = movie_trackers_rows.Err()
	if err != nil {
		log.Printf("Failed after movie_trackers_rows scan:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}

	for i := 0; i < len(movies); i++ {
		index := 0
		for _, movie_tracker := range movie_trackers {
			if movies[i].MovieName == movie_tracker.MovieName {
				movies[i].Trackers = append(movies[i].Trackers, movie_tracker)
				index = index + 1
			}
			// Only get the first 5 trackers
			if index == 5 {
				break
			}
		}
		if movies[i].Trackers == nil {
			movies[i].Trackers = make([]types.MovieTracker, 0)
		}
	}

	return c.Status(fiber.StatusOK).JSON(movies)
}

func PostMovie(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(PostMovie), txid.String())
	if security.ValidateJWT(c) != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprintf("Unauthorized: %s\n", txid.String()))
	}
	series_name := c.Params("series")
	var movies []types.Movie
	err := json.Unmarshal(c.Body(), &movies)
	if err != nil {
		log.Printf("Failed to parse movie data\n%s\n", err.Error())
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Bad Request: %s\n", txid.String()))
	}
	query := `INSERT INTO movies (movie_name, series_name, movie_title, movie_created_on) VALUES `
	values := []interface{}{}
	for _, movie := range movies {
		query += `(?, ?, ?, CURDATE()), `
		values = append(values, util.FormatName(movie.MovieTitle), series_name, movie.MovieTitle)
	}
	query = query[0 : len(query)-2]
	err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
	database, err := db.GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return c.Status(fiber.StatusInternalServerError).SendString(err_string)
	}
	result, err := database.Exec(query, values...)
	if err != nil {
		log.Printf("Failed to insert record into movies:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed retrieve inserted id\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}

	json := &fiber.Map{
		"id": id,
	}

	return c.Status(fiber.StatusOK).JSON(json)
}
