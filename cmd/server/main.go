package main

import (
	pb "M4ksR/VK-Golang-Homework.git/gen/moviesv1"
	"M4ksR/VK-Golang-Homework.git/internal/adapter/stand"
	"M4ksR/VK-Golang-Homework.git/internal/config"
	"M4ksR/VK-Golang-Homework.git/internal/handler/grpcapi"
	"M4ksR/VK-Golang-Homework.git/internal/handler/httpapi"
	"M4ksR/VK-Golang-Homework.git/internal/usecase"
	"log"
	"net"
	"net/http"

	_ "M4ksR/VK-Golang-Homework.git/docs"

	"google.golang.org/grpc"
)

// @title Movies API
// @version 1.0
// @description HTTP API для получения информации о фильмах из внешнего сервиса.
// @description Позволяет загрузить фильмы по диапазону ID и найти последний существующий фильм.
// @BasePath /
func main() {
	cfg := config.Load()
	client := stand.NewClient(cfg.MovieBaseURL)
	movieService := usecase.NewService(client)
	handler := httpapi.NewHandler(movieService)

	grpcServer := grpc.NewServer()
	pb.RegisterMovieServiceServer(grpcServer, grpcapi.NewServer(movieService))

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		err := grpcServer.Serve(listener)
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = http.ListenAndServe(cfg.HTTPAddr, handler.Router())
	if err != nil {
		log.Fatal(err)
	}
}
