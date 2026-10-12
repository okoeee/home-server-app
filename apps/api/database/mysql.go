package database

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

func Open() (*sql.DB, error) {

	config := mysql.NewConfig()

	config.User = os.Getenv("DB_USER")
	config.Passwd = os.Getenv("DB_PASSWORD")
	config.DBName = os.Getenv("DB_NAME")
	config.ParseTime = true
	config.Net = "tcp"

	config.Addr = net.JoinHostPort(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	database, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	pingErr := database.PingContext(ctx)
	if pingErr != nil {
		database.Close()

		return nil, fmt.Errorf("ping database: %w", pingErr)
	}

	return database, nil

}
