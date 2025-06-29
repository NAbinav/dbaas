package db

import (
	"context"
	"dbaas/helpers"
	"fmt"
	"strings"
)

func Read(table string, condition map[string][]string, path string) (any, error) {
	// fmt.Println(strings.Split(table, "/"))
	conditions_list := strings.Split(path, "/")
	fmt.Println(condition)
	condition_query, err := helpers.QueryRefiner(condition)
	if err != nil {
		return "", err
	}
	fmt.Println(conditions_list, condition)
	query := fmt.Sprintf("SELECT %s FROM %s %s", conditions_list[2], table, condition_query)
	fmt.Println(query)
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

	// for _, row := range results {
	// 	fmt.Println(row)
	// }

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

func Create_Table(table_name string, table_details map[string]string) error {
	fmt.Println(table_name)
	query := "CREATE TABLE " + table_name + "("

	for column_name, data_type := range table_details {
		sql_data_type, exists := helpers.SimpleNameToSQL[data_type]
		if exists == false {
			return fmt.Errorf("DataTypes not valid")
		}
		query += fmt.Sprintf("%s %s,", column_name, sql_data_type)
		fmt.Println(column_name, sql_data_type, exists)

	}
	query = query[:len(query)-1]
	query += ");"
	fmt.Println(query)
	_, err := DB.Exec(context.Background(), query)
	return err

}
