package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Ammar022/evv-logger-backend/internal/api"
	"github.com/Ammar022/evv-logger-backend/internal/middleware"
	"github.com/Ammar022/evv-logger-backend/internal/repo"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	_ "github.com/Ammar022/evv-logger-backend/docs" // Imported docs for Swagger
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Mini EVV Logger API
// @version 1.0
// @description This is the backend server for the Caregiver Shift Tracker application.
// @host localhost:8080
// @BasePath /
func main() {
	godotenv.Load()

	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbSource = fmt.Sprintf("postgresql://%s:%s@localhost:5432/%s?sslmode=disable", dbUser, dbPassword, dbName)
	}

	db, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("could not ping database: %v", err)
	}

	log.Println("Database connected successfully!")

	store := repo.NewPostgresStore(db)
	
	if err := store.Init(); err != nil {
		log.Fatalf("Could not initialize database tables: %v", err)
	}
	log.Println("Database tables initialized successfully!")
	
	if err := store.SeedData(); err != nil {
		log.Printf("Warning: Could not seed sample data: %v", err)
	} 
	router := mux.NewRouter()
	
	router.Use(middleware.CORS)
	
	api.RegisterRoutes(router, store)
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
