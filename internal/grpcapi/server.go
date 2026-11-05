package grpcapi

import (
	pb "M4ksR/VK-Golang-Homework.git/gen/moviesv1"
	"M4ksR/VK-Golang-Homework.git/internal/service"
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedMovieServiceServer
	service *service.Service
}

func NewServer(s *service.Service) *Server {
	return &Server{service: s}
}

func (s *Server) GetLastMovie(ctx context.Context, req *pb.GetLastMovieRequest) (*pb.Movie, error) {
	m, err := s.service.GetLastMovie(ctx)
	if err != nil {
		if errors.Is(err, service.ErrNoMovies) {
			return nil, status.Error(codes.NotFound, "Фильмов нет")
		}

		if errors.Is(err, context.Canceled) {
			return nil, status.Error(codes.Canceled, "Контекст вызова отменен")
		}

		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, "Истек deadline вызова")
		}

		return nil, status.Error(codes.Unavailable, "Ошибка внешнего стенда")
	}

	return &pb.Movie{
		Id:       int64(m.ID),
		Title:    m.Title,
		Year:     int32(m.Year),
		Director: m.Director,
	}, nil
}

func (s *Server) LoadMovies(ctx context.Context, req *pb.LoadMoviesRequest) (*pb.LoadMoviesResponse, error) {
	from := int(req.GetFrom())
	to := int(req.GetTo())

	if from < 1 || from > to {
		return nil, status.Error(codes.InvalidArgument, "Неверный диапазон")
	}

	result, err := s.service.LoadMovies(ctx, from, to)
	if errors.Is(err, context.Canceled) {
		return nil, status.Error(codes.Canceled, "Контекст вызова отменен")
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return nil, status.Error(codes.DeadlineExceeded, "Истек deadline вызова")
	}

	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "Ошибка загрузки фильмов: %v", err)
	}

	response := &pb.LoadMoviesResponse{}
	for _, movie := range result.Movies {
		response.Movies = append(response.Movies, &pb.Movie{
			Id:       int64(movie.ID),
			Title:    movie.Title,
			Year:     int32(movie.Year),
			Director: movie.Director,
		})
	}

	for _, failure := range result.Failures {
		response.Failures = append(response.Failures, &pb.LoadFailure{
			Id:      int64(failure.ID),
			Message: failure.Message,
		})
	}

	return response, nil
}
