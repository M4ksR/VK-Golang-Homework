package domain

type Movie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}

type LoadFailure struct {
	ID      int
	Message string
}

type LoadResult struct {
	Movies   []Movie
	Failures []LoadFailure
}
