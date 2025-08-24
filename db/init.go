package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

var DB *pgxpool.Pool

type Dummy []any

func Init_DB() {
	var err error
	DB, err = pgxpool.New(context.Background(), os.Getenv("PSQL_URL"))
	fmt.Println(os.Getenv("PSQL_URL"))
	if err != nil {
		fmt.Print(err)
		return
	}
}
