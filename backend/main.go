package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"to-mrz/api"
	"to-mrz/db"
)

func main() {
	// Initialize database
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	// Seed the database with initial data
	if err := db.SeedDB(); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// Setup router
	router := api.SetupRouter()

	// Start the server in a goroutine
	go func() {
		if err := router.Run(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("Server started on port 8080")

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// Kill (no param) default sends syscall.SIGTERM
	// Kill -2 (SIGINT) is the same as CTRL+C
	// Kill -9 (SIGKILL) can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Close database connection
	if err := db.CloseDB(); err != nil {
		log.Fatalf("Error closing database: %v", err)
	}
	log.Println("Server exited properly")
}
