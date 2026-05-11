package main

import (
	"log"
	"os"

	"expense-tracker/routes"
)

func main() {
	frontendDir := os.Getenv("FRONTEND_DIR")

	if frontendDir == "" {
		frontendDir = "../frontend/"
	}

	r := routes.SetupRoutes(frontendDir)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on http://localhost:%s", port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}