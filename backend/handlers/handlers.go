package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/milinddethe15/link-shortener/utils"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	baseURL     string
	ctx         = context.Background()
)

func Initialize(client *redis.Client, base string) {
	redisClient = client
	baseURL = base
}

// ShortenHandler shortens a URL and stores it in Redis
func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	shortKey := utils.GenerateShortKey()
	err := redisClient.Set(ctx, shortKey, originalURL, 0).Err()
	if err != nil {
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s/%s\n", baseURL, shortKey)
}

// RedirectHandler redirects the short URL to the original URL
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[1:]

	originalURL, err := redisClient.Get(ctx, shortKey).Result()
	if err == redis.Nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}
