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

func GetTrackers(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(GetTrackers), txid.String())
	err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
	database, err := db.GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return c.Status(fiber.StatusInternalServerError).SendString(err_string)
	}
	query := `
		SELECT tracker_id
			, tracker_text
			, tracker_count
			, tracker_created_on
			, tracker_updated_on
			, tracker_created_by
		FROM trackers_vw
	`
	rows, err := database.Query(query)
	if err != nil {
		log.Printf("Failed to query trackers_vw:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}

	var trackers []types.Tracker
	i := 0
	for rows.Next() {
		var tracker types.Tracker
		err = rows.Scan(&tracker.ID,
			&tracker.Text,
			&tracker.Count,
			&tracker.CreatedOn,
			&tracker.UpdatedOn,
			&tracker.CreatedBy)
		if err != nil {
			log.Printf("Failed to scan row:\n%s\n", err.Error())
			return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
		}
		tracker.Rank = i
		trackers = append(trackers, tracker)
		i = i + 1
	}

	err = rows.Err()
	if err != nil {
		log.Printf("Failed after row scan:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}

	return c.Status(fiber.StatusOK).JSON(trackers)
}

func PostTrackers(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(PostTrackers), txid.String())
	if security.ValidateJWT(c) != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprintf("Unauthorized: %s\n", txid.String()))
	}
	var tracker types.Tracker
	err := c.BodyParser(&tracker)
	if err != nil {
		log.Printf("Failed to parse tracker data\n%s\n", err.Error())
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Bad Request: %s\n", txid.String()))
	}
	err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
	database, err := db.GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return c.Status(fiber.StatusInternalServerError).SendString(err_string)
	}
	query := `
		INSERT INTO trackers
		(
			tracker_id
			, tracker_text
			, tracker_created_on
			, tracker_updated_on
			, person_id
		)
		SELECT  UUID_TO_BIN(UUID())
				, ?
				, CURDATE()
				, CURDATE()
				, person_id
		FROM people
		WHERE person_username = LOWER(?);
	`
	result, err := database.Exec(query, tracker.Text, tracker.CreatedBy)
	if err != nil {
		log.Printf("Failed to insert record into trackers:\n%s\n", err.Error())
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
