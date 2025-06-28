package db

import (
	"context"
	"fmt"
	"strings"
)

func Read(table string, condition string) (any, error) {
	fmt.Println(strings.Split(table, "/"))
	conditions_list := strings.Split(condition, "/")
	query := fmt.Sprintf("SELECT * FROM %s", table)
	fmt.Println(conditions_list)
	// err := DB.QueryRow(context.Background(), query).Scan(&data.userid, &data.name, &data.email)
	rows, err := DB.Query(context.Background(), query)
	if err != nil {
		return "", err
	}

	defer rows.Close()

	description := rows.FieldDescriptions()
	data := make([]any, len(description))
	dataptrs := make([]any, len(description))

	for i := range data {
		dataptrs[i] = &data[i]
	}

	var results []map[string]any
	for rows.Next() {
		if err := rows.Scan(dataptrs...); err != nil {
			return "", err
		}
		rowData := make(map[string]any)
		for i, desc := range description {
			rowData[string(desc.Name)] = data[i]
		}
		results = append(results, rowData)
	}
	if err := rows.Err(); err != nil {
		return "", err
	}

	for _, row := range results {
		fmt.Println(row)
	}

	return results, nil
}

func Insert(table string, data map[string]any) error {
	// Query example:
	// 	(INSERT INTO TABLE (userid,name,age) VALUES ($1 $2 $3)),values
	var values []any
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
