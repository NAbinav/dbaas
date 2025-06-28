package handler

import (
	"dbaas/db"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHandler(c *gin.Context) {
	tableName := c.Param("table")
	cndn := c.Param("cndn") // Optional: If you plan to use it later
	queries := c.Request.URL.Query()
	fmt.Println(queries)

	fmt.Println(cndn)
	result, err := db.Read("gopgx_schema."+tableName, "")
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
