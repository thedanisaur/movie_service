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
