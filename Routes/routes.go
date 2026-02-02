package routes

import (
	"To_DO_Assistant/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.POST("/task", handler.Post)
	router.GET("/task", handler.Get)
	router.DELETE("/task/delete/:id", handler.Delete)
	router.GET("/task/search", handler.SearchQuery)
	router.GET("/ask", handler.Ask)
	router.GET("/askrag", handler.AskRAG)

}
