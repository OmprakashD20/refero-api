package main

import (
	"context"
	"log"
	"time"

	"github.com/OmprakashD20/refero-api/cmd/api"
	"github.com/OmprakashD20/refero-api/config"
	"github.com/OmprakashD20/refero-api/database"
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

	// Run the server
	server := api.NewAPIServer(config.Envs.Port, conn)

	if err := server.Run(); err != nil {
		log.Fatalf("Failed to run the server: %v", err)
	}
}
