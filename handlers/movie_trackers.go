package handlers

import (
	"fmt"
	"log"
	"movie_service/db"
	"movie_service/types"
	"movie_service/util"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetMovieTrackers(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(GetMovieTrackers), txid.String())
	err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
	movie_name := c.Params("movie_name")
	username := c.Params("username")
	database := db.GetInstance()
	movie_trackers_query := `
		SELECT movie_name
			, BIN_TO_UUID(tracker_id) AS tracker_id
			, tracker_count
		FROM movie_trackers
			, people
		WHERE movie_tracker_created_by = person_id
		AND movie_name = ?
		AND person_username = ?
	`
	movie_trackers_rows, err := database.Query(movie_trackers_query, movie_name, username)
	if err != nil {
		log.Printf("Failed to query for movie trackers:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}

	var movie_trackers []types.MovieTracker
	for movie_trackers_rows.Next() {
		var movie_tracker types.MovieTracker
		err = movie_trackers_rows.Scan(&movie_tracker.MovieName,
			&movie_tracker.TrackerID,
			&movie_tracker.TrackerCount)
		if err != nil {
			log.Printf("Failed to scan movie_trackers_rows:\n%s\n", err.Error())
			return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
		}
		movie_trackers = append(movie_trackers, movie_tracker)
	}

	return c.Status(fiber.StatusOK).JSON(movie_trackers)
}

// func PostMovieTrackers(c *fiber.Ctx) error {
// 	txid := uuid.New()
// 	log.Printf("%s | %s\n", util.GetFunctionName(PostMovie), txid.String())
// 	series_name := c.Params("series")
// 	var movies []types.Movie
// 	err := json.Unmarshal(c.Body(), &movies)
// 	if err != nil {
// 		log.Printf("Failed to parse movie data\n%s\n", err.Error())
// 		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Bad Request: %s\n", txid.String()))
// 	}
// 	query := `INSERT INTO movies (movie_name, series_name, movie_title, movie_created_on) VALUES `
// 	values := []interface{}{}
// 	for _, movie := range movies {
// 		query += `(?, ?, ?, CURDATE()), `
// 		values = append(values, util.FormatName(movie.Title), series_name, movie.Title)
// 	}
// 	query = query[0 : len(query)-2]
// 	err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
// 	database := db.GetInstance()
// 	result, err := database.Exec(query, values...)
// 	if err != nil {
// 		log.Printf("Failed to insert record into movies:\n%s\n", err.Error())
// 		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
// 	}
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		log.Printf("Failed retrieve inserted id\n%s\n", err.Error())
// 		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
// 	}

// 	json := &fiber.Map{
// 		"id": id,
// 	}

// 	return c.Status(fiber.StatusOK).JSON(json)
// }
