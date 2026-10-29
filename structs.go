package main

// структура фильма
type Movie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}

// структура результата
type Result struct {
	ID    int
	Movie Movie
	Err   error
}
