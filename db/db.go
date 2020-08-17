package db

import (
	"context"
	"fmt"
	"os"

	"github.com/go-pg/pg/v10"
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	query, _ := q.FormattedQuery()
	fmt.Println(string(query))
	return nil
}

// Connect opens a connection to the db
func Connect() *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
	})
	if os.Getenv("DB_LOGGING") == "enabled" {
		db.AddQueryHook(dbLogger{})
	}
	return db
}
