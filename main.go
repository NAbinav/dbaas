package main

import (
	"dbaas/auth"
	"dbaas/db"
	"dbaas/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Init_DB()
	auth.Init_auth()

	r := gin.Default()

	r.GET("/", handler.Hi)
	r.GET("/callback", auth.App.CallbackHandler)
	r.GET("/login", auth.App.LoginHandler)
	r.POST("/create/:table_name", handler.Create_Table)
	r.GET("/:table/:column", handler.GetHandler)
	r.POST("/", handler.PostHandler)
	r.PUT("/:table_name", handler.UpdateTable)
	r.DELETE("/delete/:table_name", handler.Delete_table)
	r.DELETE("/:table_name", handler.DeleteRowHandler)
	r.GET("/newApiKey", handler.NewAPIKey)

	r.Run(":8081")
}
