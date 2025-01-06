package main

import (
	"github.com/malbertzard/beelzebub/pkg/cli"
	"github.com/malbertzard/beelzebub/pkg/db"
	"log"
)

func main() {
	// Initialize the database
	if err := db.InitDB("./storage/Beelzebub.db"); err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}

	// Execute the CLI
	if err := cli.Execute(); err != nil {
		log.Fatalf("Error starting the CLI: %v", err)
	}
}

