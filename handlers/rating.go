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

func GetRatings(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(GetRatings), txid.String())
	err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
	database := db.GetInstance()
	ratings_query := `
		SELECT series_title
			, chosen_by
			, movies_in_series
			, good_votes
			, bad_votes
			, total_votes
			, rating
		FROM rating_vw
	`
	rows, err := database.Query(ratings_query)
	if err != nil {
		log.Printf("Failed to query rating_vw:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}

	var ratings []types.Rating
	for rows.Next() {
		var rating types.Rating
		err = rows.Scan(&rating.Series,
			&rating.ChosenBy,
			&rating.MoviesInSeries,
			&rating.Good,
			&rating.Bad,
			&rating.Total,
			&rating.Rating)
		if err != nil {
			log.Printf("Failed to scan row:\n%s\n", err.Error())
			return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
		}
		ratings = append(ratings, rating)
	}

	err = rows.Err()
	if err != nil {
		log.Printf("Failed after row scan:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}

	return c.Status(fiber.StatusOK).JSON(ratings)
}
