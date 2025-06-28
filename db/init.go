package db

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

type Dummy []any

func Init_DB() {
	var err error
	DB, err = pgx.Connect(context.Background(), os.Getenv("PSQL_URL"))
	fmt.Println(os.Getenv("PSQL_URL"))
	if err != nil {
		fmt.Print(err)
		return
	}
}

func Read(table string, condition string) error {
	conditions_list := strings.Split(condition, "/")
	query := fmt.Sprintf("SELECT * FROM %s", table)
	fmt.Println(conditions_list)
	// err := DB.QueryRow(context.Background(), query).Scan(&data.userid, &data.name, &data.email)
	rows, err := DB.Query(context.Background(), query)
	if err != nil {
		return err
	}

	defer rows.Close() // Ensure rows are closed to prevent resource leaks

	description := rows.FieldDescriptions()
	data := make([]interface{}, len(description))
	dataptrs := make([]interface{}, len(description))

	// Initialize pointers for scanning
	for i := range data {
		dataptrs[i] = &data[i]
	}

	var results []map[string]interface{}
	for rows.Next() {
		// Scan the current row into dataptrs
		if err := rows.Scan(dataptrs...); err != nil {
			return err
		}

		// Create a map to store row data with column names
		rowData := make(map[string]interface{})
		for i, desc := range description {
			rowData[string(desc.Name)] = data[i]
		}
		results = append(results, rowData)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return err
	}

	// Optional: Print first column description and results
	fmt.Println(description[0])
	for _, row := range results {
		fmt.Println(row)
	}

	return nil
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
	fmt.Println(query)
	_, err := DB.Exec(context.Background(), query, values...)
	return err

}
