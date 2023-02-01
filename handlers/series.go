package handlers

import (
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

func GetSeries(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(GetSeries), txid.String())
	database := db.GetInstance()
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
	if security.ValidateJWT(c) != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprintf("Unauthorized: %s\n", txid.String()))
	}
	var series types.Series
	err := c.BodyParser(&series)
	if err != nil {
		log.Printf("Failed to parse series data\n%s\n", err.Error())
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Bad Request: %s\n", txid.String()))
	}
	series.Name = util.FormatName(series.Title)
	database := db.GetInstance()
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
	result, err := database.Exec(query, series.Title, series.Name, series.ChosenBy)
	err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
	if err != nil {
		log.Printf("Failed to insert record into series:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}
	id, err := result.LastInsertId()
	// TODO || id == 0 is there because no error was being thrown
	// for not including a username, suspicious
	// This fix probably only works because of the auto increment
	if err != nil || id == 0 {
		log.Printf("Failed retrieve inserted id\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}

	json := &fiber.Map{
		"id":          id,
		"series_name": series.Name,
	}

	return c.Status(fiber.StatusOK).JSON(json)
}
