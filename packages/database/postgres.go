package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func NewPostgres(dsn string) *bun.DB {
	sqldb := sql.OpenDB(
		pgdriver.NewConnector(
			pgdriver.WithDSN(dsn),
			pgdriver.WithTimeout(5*time.Second),
			pgdriver.WithReadTimeout(30*time.Second),
			pgdriver.WithWriteTimeout(30*time.Second),
		),
	)

	db := bun.NewDB(sqldb, pgdialect.New())

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	log.Println("connected to database")
	return db
}
