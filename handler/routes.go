package handler

import (
	"dbaas/db"
	"fmt"

	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHandler(c *gin.Context) {
	tableName := c.Param("table")
	// cndn := c.Param("cndn") // Optional: If you plan to use it later
	path := c.Request.URL.Path
	// fmt.Println(queries)

	queries := c.Request.URL.Query()
	result, err := db.Read("gopgx_schema."+tableName, queries, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
func Hi(c *gin.Context) {
	c.String(http.StatusOK, "hi")
}

func PostHandler(c *gin.Context) {
	var body map[string]interface{}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	db.Insert("gopgx_schema.Users", body)
	c.JSON(http.StatusCreated, gin.H{"status": "inserted"})
}

func Create_Table(c *gin.Context) {
	var table_details map[string]string
	if err := c.BindJSON(&table_details); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	fmt.Println(table_details)
	err := db.Create_Table("gopgx_schema."+c.Param("table_name"), table_details)
	if err != nil {
		c.JSON(400, err)
	}
	return
}
