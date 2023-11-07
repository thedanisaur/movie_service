package handlers

import (
	"fmt"
	"log"
	"movie_service/db"
	"movie_service/types"
	"movie_service/util"
	"sort"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetTimeline(c *fiber.Ctx) error {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(GetTimeline), txid.String())
	err_string := fmt.Sprintf("Database Error: %s\n", txid.String())
	database, err := db.GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return c.Status(fiber.StatusInternalServerError).SendString(err_string)
	}
	// Start by getting the series rating in order.
	rating_vw_query := `
		SELECT series_name
			, series_order
			, series_title
			, series_created_on
			, good_votes
			, bad_votes
			, rating
			, chosen_by
		FROM rating_vw
	`
	series_rating_rows, err := database.Query(rating_vw_query)
	if err != nil {
		log.Printf("Failed to query rating_vw:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}
	var series_rating []types.SeriesRating
	for series_rating_rows.Next() {
		var rating types.SeriesRating
		err = series_rating_rows.Scan(&rating.SeriesName,
			&rating.SeriesOrder,
			&rating.SeriesTitle,
			&rating.SeriesCreatedOn,
			&rating.SeriesGoodVotes,
			&rating.SeriesBadVotes,
			&rating.SeriesRating,
			&rating.SeriesChosenBy)
		if err != nil {
			log.Printf("Failed to scan series_rating_rows:\n%s\n", err.Error())
			return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
		}
		series_rating = append(series_rating, rating)
	}
	err = series_rating_rows.Err()
	if err != nil {
		log.Printf("Failed after series_rating_rows scan:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}

	// Get all the movie's votes to integrate with the series ratings
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

	// Get the series images
	series_images_query := `
		SELECT series_name
		, image_path
		FROM series_image_vw
	`
	series_images_rows, err := database.Query(series_images_query)
	if err != nil {
		log.Printf("Failed to query series_images_query:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}
	images := make(map[string]string)
	for series_images_rows.Next() {
		var image types.SeriesImage
		err = series_images_rows.Scan(&image.SeriesName, &image.ImagePath)
		if err != nil {
			log.Printf("Failed to scan series_images_rows:\n%s\n", err.Error())
			return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
		}
		images[image.SeriesName] = image.ImagePath
	}
	err = series_images_rows.Err()
	if err != nil {
		log.Printf("Failed after series_images_rows scan:\n%s\n", err.Error())
		return c.Status(fiber.StatusServiceUnavailable).SendString(err_string)
	}

	// Build the timeline
	var timeline []types.Timeline
	for i := 0; i < len(series_rating); i++ {
		var series_movies []types.TimelineMovie
		for j := 0; j < len(movies); j++ {
			if movies[j].SeriesName == series_rating[i].SeriesName {
				series_movies = append(series_movies, movies[j])
			}
		}
		if series_movies == nil {
			series_movies = make([]types.TimelineMovie, 0)
		}
		timeline = append(timeline, types.Timeline{
			SeriesOrder:     series_rating[i].SeriesOrder,
			SeriesTitle:     series_rating[i].SeriesTitle,
			SeriesRank:      i,
			SeriesGoodVotes: series_rating[i].SeriesGoodVotes,
			SeriesBadVotes:  series_rating[i].SeriesBadVotes,
			SeriesRating:    series_rating[i].SeriesRating,
			SeriesChosenBy:  series_rating[i].SeriesChosenBy,
			SeriesCreatedOn: series_rating[i].SeriesCreatedOn,
			SeriesImage:     images[series_rating[i].SeriesName],
			SeriesMovies:    series_movies})
	}

	sort.Slice(timeline, func(i, j int) bool {
		return timeline[i].SeriesOrder > timeline[j].SeriesOrder
	})

	return c.Status(fiber.StatusOK).JSON(timeline)
}
