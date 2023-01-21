package handlers

import (
	"fmt"
	"log"
	"movie_service/db"
	"movie_service/types"
	"movie_service/util"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func formatSeriesName(title string) string {
	name := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(title, "")
	name = strings.ReplaceAll(name, " ", "_")
	return strings.ToLower(name)
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
	series.Name = formatSeriesName(series.Title)
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
	if err != nil {
		log.Printf("Failed retrieve inserted id\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}

	return c.Status(fiber.StatusOK).JSON(id)
}
