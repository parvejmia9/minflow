package db

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB initializes the database connection
func ConnectDB() {
	var dsn string
	
	// Check if DATABASE_URL exists (Render/Production)
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		// Production: Use DATABASE_URL directly
		// Render provides: postgresql://user:pass@host:port/dbname
		// Ensure SSL mode is set for production
		if !strings.Contains(databaseURL, "sslmode=") {
			if strings.Contains(databaseURL, "?") {
				databaseURL += "&sslmode=require"
			} else {
				databaseURL += "?sslmode=require"
			}
		}
		dsn = databaseURL
		log.Println("Using production database (DATABASE_URL)")
	} else {
		// Development: Use individual env vars
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "minflow")
		password := getEnv("DB_PASSWORD", "minflow")
		dbname := getEnv("DB_NAME", "minflow")
		sslmode := getEnv("DB_SSLMODE", "disable")

		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode)
		log.Println("Using development database (individual env vars)")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	env := getEnv("ENV", "development")
	log.Printf("Database connected successfully! Environment: %s", env)
}


func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
