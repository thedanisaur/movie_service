package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"movie_service/db"
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
	database := db.GetInstance()
	// Get all the movies
	movies_votes_query := `
		SELECT series_name
			, movie_name
			, movie_title
			, dan_vote
			, nick_vote
		FROM dn_movies_votes_vw
	`
	movie_votes_rows, err := database.Query(movies_votes_query)
	if err != nil {
		log.Printf("Failed to query dn_movies_votes_vw:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}

	var movies []types.TimelineMovie
	for movie_votes_rows.Next() {
		var movie types.TimelineMovie
		err = movie_votes_rows.Scan(&movie.SeriesName,
			&movie.MovieName,
			&movie.MovieTitle,
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

	return c.Status(fiber.StatusOK).JSON(movies)
}

func PostMovie(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(PostMovie), txid.String())
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
		values = append(values, util.FormatName(movie.Title), series_name, movie.Title)
	}
	query = query[0 : len(query)-2]
	err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
	database := db.GetInstance()
	result, err := database.Exec(query, values...)
	if err != nil {
		log.Printf("Failed to insert record into movies:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed retrieve inserted id\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}

	return c.Status(fiber.StatusOK).JSON(id)
}
