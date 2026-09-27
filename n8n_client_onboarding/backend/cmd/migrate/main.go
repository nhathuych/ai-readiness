package main

import (
	"backend/internal/config"
	"backend/internal/database"
	"context"
	"log"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx := context.Background()

	db, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	// defer close the database pool when main() finishes
	// defer means this runs at the end of the function
	defer db.Close()

	if err := database.RunMigrations(ctx, db); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	log.Println("migrations successfully completed")

}
