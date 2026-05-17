package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func NewPostgres(dsn string) *bun.DB {

	sqldb := sql.OpenDB(
		pgdriver.NewConnector(
			pgdriver.WithDSN(dsn),
		),
	)

	db := bun.NewDB(sqldb, pgdialect.New())

	err := db.PingContext(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	log.Println("Successfully connected to the database")

	return db
}
