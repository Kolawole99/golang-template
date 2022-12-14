package main

import (
	"fmt"
	"golang-api/config"
	"golang-api/helper"
	"golang-api/routes"
	"golang-api/service"

	"log"
	"net/http"
	"os"

	"gorm.io/gorm"
)

// @title          Sample MicroService API built with Go + Gin
// @version        1.0
// @description    This is a server implementation of the microservice.
// @termsOfService http://swagger.io/terms/

// @contact.name  API Support
// @contact.url   http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host     localhost:8080
// @BasePath /api/v1
// @schemes  http
func main() {
	helper.LoadEnvironmentVariables()

	var db *gorm.DB = config.SetupDatabaseConnection()
	defer config.CloseDatabaseConnection(db)

	var jwt = service.NewJWTService()

	applicationPort := os.Getenv("PORT")

	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", applicationPort),
		Handler: routes.HandleRouting(db, jwt),
	}

	fmt.Println("Server is starting on port", applicationPort)

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
