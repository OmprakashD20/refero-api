package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/OmprakashD20/refero-api/config"
)

func InitDB(config *config.DBConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName, config.SSLMode,
	)

	conn, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify the connection
	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, err
}
