// cmd/api/main.go
package main

import (
	"clean_architecture_go/internal/api"
	"clean_architecture_go/internal/app"
	"clean_architecture_go/internal/infra"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting application...")
	infra.AddInfra()
	app.AddApp()
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server is running on Gin!",
		})
	})

	routers := router.Group("/api/v1")
	api.InjectRouters(routers)

	port := ":8090"
	log.Printf("Listening on %s\n", port)
	err := router.Run(port)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
