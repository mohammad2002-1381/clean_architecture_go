package goca

import (
	"go-ca/cmd/api"
)

// @title           Clean Architecture API
// @version         1.0
// @description     This is a sample server built with Clean Architecture in Go.

// @contact.name   API Support
// @contact.email  support@example.com

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	api.main()
}