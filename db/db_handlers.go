package db

import (
	"context"
	"dbaas/model"
	"fmt"
	"strings"
)

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
		sql_data_type, exists := model.SimpleNameToSQL[data_type]
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

func Delete_table(table_name string) error {
	query := "DROP TABLE gopgx_schema." + table_name
	_, err := DB.Exec(context.Background(), query)
	return err
}
