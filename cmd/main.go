package main

import (
	routes "To_DO_Assistant/Routes"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("[STARTUP] To-Do-Assistant")
	gin.SetMode(gin.DebugMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	routes.RegisterRoutes(router)
	router.Run(":8080")

}
