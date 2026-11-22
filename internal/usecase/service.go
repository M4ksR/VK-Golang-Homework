package usecase

import (
	"M4ksR/VK-Golang-Homework.git/internal/domain"
	"context"
	"errors"
	"fmt"
)

type Service struct {
	provider MovieProvider
}

func NewService(provider MovieProvider) *Service {
	return &Service{provider: provider}
}

func (s *Service) findBounds(ctx context.Context) (int, int, error) {
	currID, prevID := 1, 0
	for {
		_, err := s.provider.GetMovie(ctx, currID)
		if err != nil {
			if errors.Is(err, domain.ErrMovieNotFound) {
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
		_, err := s.provider.GetMovie(ctx, mid)
		if err != nil {
			if errors.Is(err, domain.ErrMovieNotFound) {
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

func (s *Service) GetLastMovie(ctx context.Context) (domain.Movie, error) {
	low, high, err := s.findBounds(ctx)
	if err != nil {
		return domain.Movie{}, err
	}

	lastID, err := s.findLastID(ctx, low, high)
	if err != nil {
		return domain.Movie{}, err
	}

	if lastID == 0 {
		return domain.Movie{}, domain.ErrNoMovies
	}

	m, err := s.provider.GetMovie(ctx, lastID)
	if err != nil {
		return domain.Movie{}, err
	}

	return m, nil
}

func (s *Service) LoadMovies(ctx context.Context, from, to int) (domain.LoadResult, error) {
	if from < 1 || from > to {
		return domain.LoadResult{}, fmt.Errorf("Неверный диапазон: from=%d, to=%d: %w", from, to, domain.ErrInvalidRange)
	}

	var results domain.LoadResult

	for id := from; id <= to; id++ {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return results, ctxErr
		}

		movie, err := s.provider.GetMovie(ctx, id)

		if ctxErr := ctx.Err(); ctxErr != nil {
			return results, ctxErr
		}

		if err != nil {
			failure := domain.LoadFailure{ID: id, Message: err.Error()}
			results.Failures = append(results.Failures, failure)
			continue
		}
		results.Movies = append(results.Movies, movie)
	}

	return results, nil
}
