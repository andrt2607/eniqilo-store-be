package main

import (
	"context"
	"log"
	"net/http"

	"eniqilo-store-be/internal/cfg"
	"eniqilo-store-be/internal/handler"
	"eniqilo-store-be/internal/repo"
	"eniqilo-store-be/internal/service"
	"eniqilo-store-be/pkg/env"
	"eniqilo-store-be/pkg/postgre"
	"eniqilo-store-be/pkg/router"
	"eniqilo-store-be/pkg/validator"
)

func main() {
	env.LoadEnv()

	ctx := context.Background()
	router := router.NewRouter()
	conn := postgre.GetConn(ctx)
	defer conn.Close()
	validator := validator.New()

	cfg := cfg.Load()
	repo := repo.NewRepo(conn)
	service := service.NewService(repo, validator, cfg)
	handler.NewHandler(router, service, cfg)

	log.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalln("fail start server:", err)
	}
}
