package config

import "os"

type Config struct {
	MovieBaseURL string
	HTTPAddr     string
}

func Load() Config {
	MOVIE_BASE_URL := os.Getenv("MOVIE_BASE_URL")
	HTTP_ADDR := os.Getenv("HTTP_ADDR")

	config := Config{
		MovieBaseURL: "https://homeworksite.site",
		HTTPAddr:     ":8080",
	}

	if MOVIE_BASE_URL != "" {
		config.MovieBaseURL = MOVIE_BASE_URL
	}

	if HTTP_ADDR != "" {
		config.HTTPAddr = HTTP_ADDR
	}

	return config
}
