package config

import (
	"golang-api/model"

	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupDatabaseConnection is creating a new connection to our database
func SetupDatabaseConnection() *gorm.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	databaseName := os.Getenv("DB_NAME")
	ssl := os.Getenv("SSL_MODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, username, password, databaseName, port, ssl)
	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN:                  dsn,
				PreferSimpleProtocol: true, // disables implicit prepared statement usage
			},
		),
		&gorm.Config{},
	)

	if err != nil {
		log.Fatal("Error connecting to the Database")
	}

	fmt.Println("Connected to Database")

	fmt.Println("Migrating Schema")

	err = db.AutoMigrate(
		model.User{},
	)

	if err != nil {
		fmt.Printf("Could not migrate to db. Errors are: %v", err)
	}

	fmt.Println("Migration successful")

	return db
}

// CloseDatabaseConnection is calling the close connection method to close the existing connection to our database
func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()

	if err != nil {
		panic("Error in closing connection")
	}

	dbSQL.Close()
}
