package stand

import (
	"M4ksR/VK-Golang-Homework.git/internal/movie"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrMovieNotFound = errors.New("Фильм не найден")

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *Client) FetchMovie(ctx context.Context, id int) (movie.Movie, error) {
	link := fmt.Sprintf("%s/%d/info.0.json", c.baseURL, id)

	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, "GET", link, nil)
	if err != nil {
		return movie.Movie{}, fmt.Errorf("Для фильма %d, создание запроса: %w", id, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return movie.Movie{}, fmt.Errorf("Для фильма %d, ошибка запроса: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return movie.Movie{}, fmt.Errorf(
				"Для фильма %d, HTTP статус: %d, ошибка: %w",
				id,
				resp.StatusCode,
				ErrMovieNotFound,
			)
		default:
			return movie.Movie{}, fmt.Errorf(
				"Для фильма %d, HTTP статус: %d",
				id,
				resp.StatusCode,
			)
		}
	}

	var m movie.Movie
	jsonDecoder := json.NewDecoder(resp.Body)
	err = jsonDecoder.Decode(&m)

	if err != nil {
		return movie.Movie{}, fmt.Errorf("Для фильма %d, ошибка json: %w", id, err)
	}

	return m, nil
}
