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

func GetMovieTrackers(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(GetMovieTrackers), txid.String())
	err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
	movie_name := c.Params("movie_name")
	username := c.Params("username")
	database, err := db.GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
		return c.Status(fiber.StatusInternalServerError).SendString(err_string)
	}
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

func GetMovieTrackersByID(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(GetMovieTrackersByID), txid.String())
	tracker_id := c.Params("tracker_id")
	// Don't worry about converting the uuid back,
	// the mt_vw already has tracker_id in uuid format
	tracker_movies_query := `
        SELECT m.movie_title
            , mt_vw.tracker_count
        FROM movies_trackers_vw mt_vw
            , movies m
        WHERE mt_vw.movie_name = m.movie_name
        AND mt_vw.tracker_id = ?
        ORDER BY mt_vw.tracker_count DESC
	`
	err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
	database, err := db.GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
		return c.Status(fiber.StatusInternalServerError).SendString(err_string)
	}
	tracker_movies_rows, err := database.Query(tracker_movies_query, tracker_id)
	if err != nil {
		log.Printf("Failed to query for tracker's movies:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}

	var tracker_movies []types.MovieTracker2
	for tracker_movies_rows.Next() {
		var tracker_movie types.MovieTracker2
		err = tracker_movies_rows.Scan(&tracker_movie.MovieTitle,
			&tracker_movie.TrackerCount)
		if err != nil {
			log.Printf("Failed to scan tracker_movies_rows:\n%s\n", err.Error())
			return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
		}
		tracker_movies = append(tracker_movies, tracker_movie)
	}

	return c.Status(fiber.StatusOK).JSON(tracker_movies)
}

func PostMovieTrackers(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(PostMovieTrackers), txid.String())
	if security.ValidateJWT(c) != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprintf("Unauthorized: %s\n", txid.String()))
	}
	username := c.Params("username")
	var movie_trackers []types.MovieTracker
	err := json.Unmarshal(c.Body(), &movie_trackers)
	if err != nil {
		log.Printf("Failed to parse movie data\n%s\n", err.Error())
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Bad Request: %s\n", txid.String()))
	}
	inserts := 0
	updates := 0
	for i := 0; i < len(movie_trackers); i++ {
		movie_tracker := movie_trackers[i]
		query := `
			SELECT COUNT(*)
			FROM movie_trackers
			WHERE tracker_id = UUID_TO_BIN(?)
			AND movie_name = ?
		`
		database, err := db.GetInstance()
		if err != nil {
			log.Printf("Failed to connect to DB\n%s\n", err.Error())
			err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
			return c.Status(fiber.StatusInternalServerError).SendString(err_string)
		}
		row := database.QueryRow(query, movie_tracker.TrackerID, movie_tracker.MovieName)
		var val int
		err = row.Scan(&val)
		if err != nil {
			log.Printf("Database Error:\n%s\n", err.Error())
			err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
			return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
		}
		// Update!
		if val != 0 {
			query := `
                UPDATE movie_trackers
                INNER JOIN people ON movie_trackers.movie_tracker_created_by = people.person_id
                SET movie_trackers.tracker_count = ?
                WHERE movie_trackers.movie_name = ?
                AND movie_trackers.tracker_id = UUID_TO_BIN(?)
                AND movie_trackers.movie_tracker_created_by = people.person_id
                AND people.person_username = ?
            `
			result, err := database.Exec(query, &movie_tracker.TrackerCount, &movie_tracker.MovieName, &movie_tracker.TrackerID, username)
			err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
			if err != nil {
				log.Printf("Failed to update record in movie_trackers:\n%s\n", err.Error())
				return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
			}
			_, err = result.LastInsertId()
			if err != nil {
				log.Printf("Failed retrieve inserted id\n%s\n", err.Error())
				return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
			}

		} else {
			// Insert
			query := `
                INSERT INTO movie_trackers (
                    movie_name
                    , tracker_id
                    , tracker_count
                    , movie_tracker_created_by
                )
                SELECT ?
                    , UUID_TO_BIN(?)
                    , ?
                    , p.person_id
                FROM people p
                WHERE p.person_username = ?
            `
			result, err := database.Exec(query, &movie_tracker.MovieName, &movie_tracker.TrackerID, &movie_tracker.TrackerCount, username)
			err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
			if err != nil {
				log.Printf("Failed to insert record into movie_trackers:\n%s\n", err.Error())
				return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
			}
			_, err = result.LastInsertId()
			if err != nil {
				log.Printf("Failed retrieve inserted id\n%s\n", err.Error())
				return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
			}
		}
	}

	json := &fiber.Map{
		"id":      txid.String(),
		"inserts": inserts,
		"updates": updates,
	}

	return c.Status(fiber.StatusOK).JSON(json)
}
