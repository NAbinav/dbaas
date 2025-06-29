package main

import (
	"dbaas/db"
	"dbaas/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	db.Init_DB()

	r := gin.Default()

	r.GET("/hi", handler.Hi)
	r.POST("/create/:table_name", handler.Create_Table)
	r.GET("/:table/:cndn", handler.GetHandler)
	r.POST("/", handler.PostHandler)

	r.Run(":8080")
}
