package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/A-Victory/blog/auth"
	"github.com/A-Victory/blog/database"
	"github.com/A-Victory/blog/database/conn"
	"github.com/A-Victory/blog/routes"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load godotenv: %v", err)
	}

	dbConfig := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")

	dbConnection := database.NewDBConn(dbConfig, dbName)
	if err := dbConnection.Initialize(); err != nil {
		log.Fatalf("failed to initialize tables: %v", err)
	}
	conn := conn.NewConn(dbConnection)
	validator := auth.NewValidator()

	serverConfig := routes.ServerConfig{
		DB: conn,
		VA: validator,
	}

	server := routes.NewServer(serverConfig)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	address := fmt.Sprintf(":%s", port)
	log.Printf("starting server on port %s", address)
	if err := http.ListenAndServe(address, server); err != nil {
		log.Fatal("failed to start server on port " + address)
	}

}
