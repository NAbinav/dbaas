package db

import (
	"context"
	"dbaas/auth"
	"dbaas/helpers"
	"dbaas/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strings"
)

func extract_value(data map[string]any) ([]string, []string, []any) {
	var values []any
	var columns []string
	var placeholders []string
	for column, column_value := range data {
		columns = append(columns, column)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)+1))
		values = append(values, column_value)
	}
	return columns, placeholders, values
}

func Insert(table string, data map[string]any) error {
	// Query example:
	// 	(INSERT INTO TABLE (userid,name,age) VALUES ($1 $2 $3)),values
	columns, placeholders, values := extract_value(data)
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

func DeleteRow(table_name string, condition map[string][]string) error {
	condition_query, err := helpers.Condition_extract(condition)
	if err != nil {
		return err
	}
	fmt.Println(condition_query)
	query := "DELETE FROM " + table_name + " " + condition_query
	fmt.Println(query)
	commandTag, err := DB.Exec(context.Background(), query)
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("No Rows Affected")
	}
	if err != nil {
		return err
	}
	return nil
}

func UpdateRow(table_name string, condition map[string][]string, changes map[string]any) error {
	column, placeholder, values := extract_value(changes)
	query, err := helpers.UpdateQuery(table_name, column, placeholder)
	if err != nil {
		return err
	}
	cndn_query, err := helpers.Condition_extract(condition)
	if err != nil {
		return err
	}
	query += cndn_query
	fmt.Println(query, values)
	_, err = DB.Exec(context.Background(), query, values...)
	if err != nil {
		return err
	}
	return nil
}

func InsertAPI(email string) (string, error) {
	apiKey, err := auth.GenerateAPIKey()
	fmt.Println(apiKey)
	if err != nil {
		return "", err
	}
	existingKey, found, err := APIExists(email)
	if err != nil {
		return "", err
	}
	if found {
		return existingKey, nil
	}

	err = InsertEmailApi(apiKey, email)
	if err != nil {
		return "", err
	}

	fmt.Println(apiKey)
	return apiKey, nil
}

func APIExists(email string) (string, bool, error) {
	query := `SELECT apikey FROM api_keys.keytable WHERE email_id = $1;`
	data, err := DB.Query(context.Background(), query, email)
	if err != nil {
		fmt.Println("Error querying DB:", err)
		return "", false, err
	}

	fmt.Println(data.FieldDescriptions())

	results, err := ReadFromQuery(data)
	if err != nil {
		return "", false, err
	}

	if len(results) > 0 {
		return (results[0]["apikey"]).(string), true, nil
	}

	return "", false, nil
}

func InsertEmailApi(api string, email string) error {
	query := "INSERT INTO api_keys.keytable  (apikey, email_id, tablenames) values ($1,$2,ARRAY[]::TEXT[]);"
	err := UserSchemaCreation(api)
	if err != nil {
		return err
	}
	_, err = DB.Exec(context.Background(), query, api, email)
	return err
}

func FormatAPIKEY(api string) (string, error) {
	// Regex: allow letters, numbers, underscore, dash; length 20+
	re := regexp.MustCompile(`^[a-zA-Z0-9_\-]{20,}$`)
	if !re.MatchString(api) {
		return "", fmt.Errorf("invalid API key format")
	}

	// Replace dashes with underscores (optional normalization)
	schema := "user_" + strings.ReplaceAll(api, "-", "_")

	// Wrap in double quotes to safely use as schema name in SQL
	return fmt.Sprintf(`"%s"`, schema), nil
}

func UserSchemaCreation(api string) error {
	schemaName, err := FormatAPIKEY(api)
	if err != nil {
		return err
	}
	query := "CREATE SCHEMA IF NOT EXISTS " + (schemaName)
	fmt.Println("Running query:", query)

	_, err = DB.Exec(context.Background(), query)
	if err != nil {
		fmt.Println("Error creating schema:", err)
	}
	return err
}
func IsValidAPIKey(key string) bool {
	query := `SELECT 1 FROM api_keys.keytable WHERE apikey = $1 LIMIT 1;`
	row := DB.QueryRow(context.Background(), query, key)
	var exists int
	err := row.Scan(&exists)
	return err == nil
}

func ValidateAPIHeader(key string) bool {

	query := `SELECT 1 FROM api_keys.keytable WHERE apikey = $1 LIMIT 1;`

	row := DB.QueryRow(context.Background(), query, key)
	var exists int
	err := row.Scan(&exists)
	return err == nil

}

func CheckTableWithAPI(apikey, table_name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM api_keys.keytable WHERE $1 = ANY(tablenames) AND apikey = $2);`
	fmt.Println(query)
	err := DB.QueryRow(context.Background(), query, table_name, apikey).Scan(&exists)
	return exists, err
}

func TableNameToAPIKEY(table string, api string) error {
	query := `UPDATE api_keys.keytable SET tablenames = array_append(tablenames, $1) WHERE apikey = $2 AND NOT ($1 = ANY(tablenames));`
	_, err := DB.Exec(context.Background(), query, table, api)
	return err
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		// email, err := GetEmailWithAPI(apiKey)
		// if err != nil {
		// 	c.JSON(416, "Nahh")
		// }
		// fmt.Println(email)
		table := c.Param("table_name")
		if apiKey == "" || !IsValidAPIKey(apiKey) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing API key"})
			c.Abort()
			return
		}

		if c.Request.Method == "POST" && strings.HasPrefix(c.FullPath(), "/create/:table_name") {
			c.Next()
			fmt.Println("Hello")
			err := TableNameToAPIKEY(table, apiKey)
			if err != nil {
				fmt.Printf("Failed to update permissions: %v\n", err)
			}
			return
		}
		allowed, err := CheckTableWithAPI(apiKey, table)
		fmt.Println(allowed, err)

		if err != nil || !allowed {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing API key"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetEmailWithAPI(api string) (string, error) {
	query := `SELECT email_id FROM api_keys.keytable WHERE apikey = $1;`
	data, err := DB.Query(context.Background(), query, api)
	if err != nil {
		fmt.Println("Error querying DB:", err)
		return "", err
	}

	results, err := ReadFromQuery(data)
	if err != nil {
		return "", err
	}

	if len(results) > 0 {
		return (results[0]["email_id"]).(string), nil
	}

	return "", fmt.Errorf("Not avail")

}
