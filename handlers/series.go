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

func GetSeries(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(GetSeries), txid.String())
	database, err := db.GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
		return c.Status(fiber.StatusInternalServerError).SendString(err_string)
	}
	query := `
		SELECT series_order
			, series_name
			, series_title
			, series_created_on
			, BIN_TO_UUID(person_id) person_id
		FROM series
	`
	series_rows, err := database.Query(query)
	err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
	if err != nil {
		log.Printf("Failed to query series:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}

	var series_array []types.Series
	for series_rows.Next() {
		var series types.Series
		err = series_rows.Scan(&series.Order,
			&series.Name,
			&series.Title,
			&series.CreatedOn,
			&series.PersonID)
		if err != nil {
			log.Printf("Failed to scan series_rows:\n%s\n", err.Error())
			return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
		}
		series_array = append(series_array, series)
	}

	err = series_rows.Err()
	if err != nil {
		log.Printf("Failed after series_rows scan:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}

	return c.Status(fiber.StatusOK).JSON(series_array)
}

func PostSeries(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(PostSeries), txid.String())
	var series types.Series
	err := c.BodyParser(&series)
	if err != nil {
		log.Printf("Failed to parse series data\n%s\n", err.Error())
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Bad Request: %s\n", txid.String()))
	}
	series.Name = util.FormatName(series.Title)
	username := c.Get("Username")
	database, err := db.GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
		return c.Status(fiber.StatusInternalServerError).SendString(err_string)
	}
	query := `
		INSERT INTO series
		(
			series_title
			, series_name
			, series_created_on
			, person_id
		)
		SELECT ?
			, ?
			, CURDATE()
			, person_id
		FROM people
		WHERE person_username = ?;
	`
	result, err := database.Exec(query, series.Title, series.Name, username)
	err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
	if err != nil {
		log.Printf("Failed to insert record into series:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed retrieve inserted id\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}

	json := &fiber.Map{
		"id":          id,
		"series_name": series.Name,
	}

	return c.Status(fiber.StatusOK).JSON(json)
}
