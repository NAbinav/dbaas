package db

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func Init_DB() {
	var err error
	DB, err = pgx.Connect(context.Background(), os.Getenv("PSQL_URL"))
	fmt.Println(os.Getenv("PSQL_URL"))
	if err != nil {
		fmt.Print(err)
		return
	}
}

func Insert(table string, data map[string]interface{}) error {
	// Query example:
	// 	(INSERT INTO TABLE (userid,name,age) VALUES ($1 $2 $3)),values
	var values []interface{}
	var columns []string
	var placeholders []string
	for column, column_value := range data {
		columns = append(columns, column)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)+1))
		values = append(values, column_value)
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, " ,"), strings.Join(placeholders, " ,"))
	_, err := DB.Exec(context.Background(), query, values...)
	return err

}
