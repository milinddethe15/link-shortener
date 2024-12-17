package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/milinddethe15/link-shortener/db"
	"github.com/milinddethe15/link-shortener/handlers"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	redisURL := os.Getenv("REDIS_URL")
	baseURL := os.Getenv("BASE_URL")
	redisClient := db.Connect(redisURL)
	defer redisClient.Close()

	handlers.Initialize(redisClient, baseURL)

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow all origins
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// Serve endpoints
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", handlers.ShortenHandler)
	mux.HandleFunc("/", handlers.RedirectHandler)

	handler := c.Handler(mux)

	fmt.Println("Server is running on :8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
