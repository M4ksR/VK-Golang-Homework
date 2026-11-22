package domain

import "errors"

var ErrMovieNotFound = errors.New("Фильм не найден")
var ErrNoMovies = errors.New("На стенде отсутствуют фильмы")
var ErrInvalidRange = errors.New("Неверный диапазон")
