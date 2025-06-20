package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"movie_service/db"
	"movie_service/types"
	"movie_service/util"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func PostVote(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(PostVote), txid.String())
	var vote_data types.Vote
	err := json.Unmarshal(c.Body(), &vote_data)
	if err != nil {
		log.Printf("Failed to parse vote data\n%s\n", err.Error())
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Bad Request: %s\n", txid.String()))
	}
	if len(vote_data.MovieName) == 0 || len(vote_data.PersonUsername) == 0 || len(vote_data.Value) == 0 {
		log.Printf("Missing vote data\n%s\n", vote_data)
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Bad Request: %s\n", txid.String()))
	}
	database, err := db.GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
		return c.Status(fiber.StatusInternalServerError).SendString(err_string)
	}
	var vote_value types.Vote
	vote_query := `
		SELECT COUNT(vote_value)
		FROM votes v
			, people p
		WHERE v.person_id = p.person_id
		AND v.movie_name = ?
		AND person_username = ?
	`
	// First check to see if a vote already exists
	row := database.QueryRow(vote_query, vote_data.MovieName, &vote_data.PersonUsername)
	err = row.Scan(&vote_value.Value)
	if err != nil {
		log.Printf("Database Error:\n%s\n", err.Error())
		err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}
	val, err := strconv.Atoi(vote_value.Value)
	if err != nil {
		log.Printf("Value Error:\n%s\n", err.Error())
		err_string := fmt.Sprintf("Vote value error: %s\n", txid.String())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}
	if val != 0 {
		log.Printf("Database Error: Vote already cast\n")
		err_string := fmt.Sprintf("Vote already cast: %s\n", txid.String())
		return c.Status(fiber.StatusConflict).SendString(err_string)
	}

	query := `
		INSERT INTO votes (movie_name
			, vote_value
			, person_id
		)
		SELECT ?
			, ?
			, person_id
		FROM people
		WHERE person_username = ?
	`
	result, err := database.Exec(query, &vote_data.MovieName, &vote_data.Value, &vote_data.PersonUsername)
	err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
	if err != nil {
		log.Printf("Failed to insert record into votes:\n%s\n", err.Error())
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
