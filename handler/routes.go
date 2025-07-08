package handler

import (
	"dbaas/auth"
	"dbaas/db"
	"fmt"
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var TablePrefix = "gopgx_schema."

func GetHandler(c *gin.Context) {
	tableName := c.Param("table_name")
	// cndn := c.Param("cndn") // Optional: If you plan to use it later
	path := c.Request.URL.Path
	// fmt.Println(queries)

	// apiKey := c.GetHeader("X-API-Key")
	// schema_name := db.FormatAPIKEY(apiKey)
	queries := c.Request.URL.Query()
	result, err := db.Read(TablePrefix+tableName, queries, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

var Welcome = `API Documentation

Welcome to the API documentation for the DBaaS (Database as a Service) application. This service allows you to interact with a database using HTTP requests. You can perform CRUD (Create, Read, Update, Delete) operations on any table in your database.

    GET /: Displays a simple welcome message.


    GET /newApiKey: Generates a new API key for authenticated users.

    POST /create/:table_name: Creates a new table in the database with the specified name.

    GET /:table_name/:column: Fetches data from a specific column in the specified table.

    POST /:table_name: Inserts a new row of data into the specified table.

    PUT /:table_name: Updates data in the specified table.

    DELETE /delete/:table_name: Deletes the entire table from the database.

    DELETE /:table_name: Deletes a row from the specified table.`

func Hi(c *gin.Context) {
	c.String(http.StatusOK, Welcome)
}

func PostHandler(c *gin.Context) {
	var body map[string]interface{}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	apiKey := c.GetHeader("X-API-Key")
	fmt.Println(apiKey)
	schema_name, err := db.FormatAPIKEY(apiKey)
	if err != nil {
		c.JSON(400, err)
	}

	table := c.Param("table_name")
	db.Insert(schema_name+"."+table, body)
	c.JSON(http.StatusCreated, gin.H{"status": "inserted"})
	return
}

func Create_Table(c *gin.Context) {
	var table_details map[string]string
	if err := c.BindJSON(&table_details); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	apiKey := c.GetHeader("X-API-Key")
	schema_name, err := db.FormatAPIKEY(apiKey)
	fmt.Println(schema_name)
	fmt.Println(table_details)
	err = db.Create_Table(schema_name+"."+c.Param("table_name"), table_details)
	if err != nil {
		c.JSON(400, err)
		return
	}
}
func Delete_table(c *gin.Context) {
	var table_name string
	table_name = c.Param("table_name")
	apiKey := c.GetHeader("X-API-Key")
	schema_name, err := db.FormatAPIKEY(apiKey)
	if err != nil {
		c.JSON(400, err)
	}
	err = db.Delete_table(schema_name + "." + table_name)
	if err != nil {
		c.JSON(400, err)
	}
}

func DeleteRowHandler(c *gin.Context) {
	table_name := c.Param("table_name")
	condition := c.Request.URL.Query()
	apiKey := c.GetHeader("X-API-Key")
	schema_name, err := db.FormatAPIKEY(apiKey)
	if err != nil {
		c.JSON(400, err)
	}
	err = db.DeleteRow(schema_name+table_name, condition)
	if err != nil {
		c.JSON(400, "Bad Request")
		return
	}

}

func UpdateTable(c *gin.Context) {
	table_name := c.Param("table_name")
	condition := c.Request.URL.Query()
	apiKey := c.GetHeader("X-API-Key")
	schema_name, err := db.FormatAPIKEY(apiKey)
	if err != nil {
		c.JSON(400, err)
	}
	var body map[string]interface{}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	err = db.UpdateRow(schema_name+table_name, condition, body)
	if err != nil {
		c.JSON(400, err)
	}

}

func NewAPIKey(c *gin.Context) {
	cookies, ok := auth.CheckAndVerifyCookies(c)
	if ok == false {
		c.String(400, "Baddd")

	}
	fmt.Println(cookies)
	apiKey, err := db.InsertAPI(cookies)
	// err = db.UserSchemaCreation(apiKey)
	fmt.Println(apiKey)
	fmt.Println(err)
	if err != nil {
		c.JSON(400, fmt.Errorf("Bad Request"))
		c.Abort()

	}
	c.String(200, apiKey)

	// c.JSON(200, apiKey)
	// if err != nil {
	// 	c.AbortWithStatusJSON(400, "Baddd")
	// }
}
