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

func GetTrackers(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(GetTrackers), txid.String())
	err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
	database := db.GetInstance()
	rows, err := database.Query("SELECT tracker_id, tracker_text, tracker_count, tracker_created_on, tracker_updated_on, tracker_created_by FROM trackers_vw")
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
