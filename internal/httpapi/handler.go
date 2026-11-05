package httpapi

import (
	"M4ksR/VK-Golang-Homework.git/internal/movie"
	"M4ksR/VK-Golang-Homework.git/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

// @Summary Получить последний фильм
// @Description Находит последний существующий фильм
// @Tags movies
// @Produce json
// @Success 200 {object} movie.Movie "Последний фильм"
// @Failure 404 {string} string "Фильмы не найдены"
// @Failure 502 {string} string "Ошибка внешнего стенда"
// @Router /last-movie [get]
func (h *Handler) getLastMovie(w http.ResponseWriter, r *http.Request) {
	m, err := h.service.GetLastMovie(r.Context())

	if err != nil {
		if errors.Is(err, service.ErrNoMovies) {
			http.Error(w, "Фильм не найден", http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка обращения к внешнему стенду", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

// @Summary Загрузить фильмы
// @Description Загружает фильмы по диапазону ID
// @Tags movies
// @Produce json
// @Param from query int true "Начальный ID фильма"
// @Param to query int true "Конечный ID фильма"
// @Success 200 {array} movie.Movie "Список фильмов"
// @Failure 400 {string} string "Некорректные параметры"
// @Failure 502 {string} string "Ошибка загрузки фильмов"
// @Router /load-movies [get]
func (h *Handler) loadMovies(w http.ResponseWriter, r *http.Request) {
	fromString := r.URL.Query().Get("from")
	toString := r.URL.Query().Get("to")

	from, errFrom := strconv.Atoi(fromString)
	to, errTo := strconv.Atoi(toString)

	if errFrom != nil || errTo != nil || from < 1 || from > to {
		http.Error(w, "Некорректные параметры", http.StatusBadRequest)
		return
	}

	result, err := h.service.LoadMovies(r.Context(), from, to)
	if err != nil {
		http.Error(w, "Ошибка загрузки фильмов", http.StatusBadGateway)
		return
	}

	if result.Movies == nil {
		result.Movies = []movie.Movie{}
	}

	w.Header().Set("Content-Type", "application/json")

	if size := len(result.Failures); size > 0 {
		IDs := make([]string, 0, size)
		for _, failure := range result.Failures {
			IDs = append(IDs, strconv.Itoa(failure.ID))
		}

		failedIDs := strings.Join(IDs, ",")
		failedCount := strconv.Itoa(size)

		w.Header().Set("X-Failed-Movie-IDs", failedIDs)
		w.Header().Set("X-Failed-Movie-Count", failedCount)
	}

	json.NewEncoder(w).Encode(result.Movies)
}

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /last-movie", h.getLastMovie)
	mux.HandleFunc("GET /load-movies", h.loadMovies)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return mux
}
