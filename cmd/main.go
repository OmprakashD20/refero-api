package main

import (
	"context"
	"log"
	"time"

	"github.com/OmprakashD20/refero-api/cmd/api"
	"github.com/OmprakashD20/refero-api/config"
	"github.com/OmprakashD20/refero-api/database"
	"github.com/OmprakashD20/refero-api/repository"
)

func main() {
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer dbCancel()

	// Init Database Connection
	conn, err := database.InitDB(dbCtx, &config.Envs.DB)
	if err != nil {
		log.Fatalf("Failed connecting the database: %v", err)
	}

	defer conn.Close()

	log.Println("Connected to the database successfully")

	// Init SQLC Queries
	db := repository.New(conn)

	// Run the server
	server := api.NewAPIServer(config.Envs.Port, db)

	if err := server.Run(); err != nil {
		log.Fatalf("Failed to run the server: %v", err)
	}
}
