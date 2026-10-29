package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func fetchMovie(ctx context.Context, id int, duration time.Duration) (Movie, error) {
	link := fmt.Sprintf("https://homeworksite.site/%d/info.0.json", id)

	// создадим дочерный контекст с таймаутом для запроса
	reqCtx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	// создадим запрос
	req, err := http.NewRequestWithContext(reqCtx, "GET", link, nil)
	if err != nil {
		return Movie{}, fmt.Errorf("Для фильма %d, создание запроса: %w", id, err)
	}

	// выполним запрос
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Movie{}, fmt.Errorf("Для фильма %d выполнение запроса: %w", id, err)
	}

	// сделаем отложенное закрытие тела HTTP ответа после окончания работы функции
	defer resp.Body.Close()

	// проверим статус
	if resp.StatusCode != http.StatusOK {
		return Movie{}, fmt.Errorf("Для фильма %d, HTTP статус: %d", id, resp.StatusCode)
	}

	var movie Movie

	// создадим json_decoder и извлечем содержимое в movie через .Decode()
	json_decoder := json.NewDecoder(resp.Body)
	err = json_decoder.Decode(&movie)

	// проверим правильность декодирования json
	if err != nil {
		return Movie{}, fmt.Errorf("Для фильма %d, ошибка чтения json: %w", id, err)
	}

	// в случае если все проверки пройдены запишем наш фильм в results
	return movie, nil
}

func getMovies(ctx context.Context, jobs <-chan int, results chan<- Result, duration time.Duration) {
	for id := range jobs {
		// если программа уже отменена то не начинаем новый HTTP запрос
		if ctx.Err() != nil {
			return
		}

		movie, err := fetchMovie(ctx, id, duration)

		// если отмена произошла во время HTTP запроса
		if ctx.Err() != nil {
			return
		}

		results <- Result{
			ID:    id,
			Movie: movie,
			Err:   err,
		}
	}
}
