package main

import (
	"log"
	"movie_service/db"
	"movie_service/handlers"
	"net/http"
)

func main() {
	defer db.GetInstance().Close()
	http.HandleFunc(handlers.RatingRoute(), handlers.RatingHandler)
	http.HandleFunc(handlers.TimelineRoute(), handlers.TimelineHandler)
	http.HandleFunc(handlers.TrackerRoute(), handlers.TrackerHandler)
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		log.Fatal("Error starting http server:", err)
		return
	}
}
