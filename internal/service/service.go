package service

import (
	"M4ksR/VK-Golang-Homework.git/internal/movie"
	"M4ksR/VK-Golang-Homework.git/internal/stand"
	"context"
	"errors"
	"fmt"
)

var ErrNoMovies = errors.New("На стенде отсутствуют фильмы")

type Service struct {
	client *stand.Client
}

func NewService(client *stand.Client) *Service {
	return &Service{client: client}
}

type LoadFailure struct {
	ID      int
	Message string
}

type LoadResult struct {
	Movies   []movie.Movie
	Failures []LoadFailure
}

func (s *Service) findBounds(ctx context.Context) (int, int, error) {
	currID, prevID := 1, 0
	for {
		_, err := s.client.FetchMovie(ctx, currID)
		if err != nil {
			if errors.Is(err, stand.ErrMovieNotFound) {
				return prevID, currID, nil
			}
			return 0, 0, err
		}
		prevID = currID
		currID *= 2
	}
}

func (s *Service) findLastID(ctx context.Context, left, right int) (int, error) {
	for right-left > 1 {
		mid := left + (right-left)/2
		_, err := s.client.FetchMovie(ctx, mid)
		if err != nil {
			if errors.Is(err, stand.ErrMovieNotFound) {
				right = mid
			} else {
				return 0, err
			}
		} else {
			left = mid
		}
	}
	return left, nil
}

func (s *Service) GetLastMovie(ctx context.Context) (movie.Movie, error) {
	low, high, err := s.findBounds(ctx)
	if err != nil {
		return movie.Movie{}, err
	}

	lastID, err := s.findLastID(ctx, low, high)
	if err != nil {
		return movie.Movie{}, err
	}

	if lastID == 0 {
		return movie.Movie{}, ErrNoMovies
	}

	m, err := s.client.FetchMovie(ctx, lastID)
	if err != nil {
		return movie.Movie{}, err
	}

	return m, nil
}

func (s *Service) LoadMovies(ctx context.Context, from, to int) (LoadResult, error) {
	if from < 1 || from > to {
		return LoadResult{}, fmt.Errorf("Неверный диапазон: from=%d, to=%d", from, to)
	}

	var results LoadResult

	for id := from; id <= to; id++ {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return results, ctxErr
		}

		movie, err := s.client.FetchMovie(ctx, id)

		if ctxErr := ctx.Err(); ctxErr != nil {
			return results, ctxErr
		}

		if err != nil {
			failure := LoadFailure{id, err.Error()}
			results.Failures = append(results.Failures, failure)
			continue
		}
		results.Movies = append(results.Movies, movie)
	}

	return results, nil
}
