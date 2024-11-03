package cmd

import (
	"log"
	"os"

	"github.com/Amierza/e-wallet/migrations"
	"gorm.io/gorm"
)

func Command(db *gorm.DB) {
	migrate := false
	seed := false

	for _, arg := range os.Args[1:] {
		if arg == "--migrate" {
			migrate = true
		}
		if arg == "--seed" {
			seed = true
		}
	}

	if migrate {
		if err := migrations.Migrate(db); err != nil {
			log.Fatal("error migration: %v", err)
		}
		log.Println("migration completed successfully!")
	}

	if seed {
		if err := migrations.Seeder(db); err != nil {
			log.Fatal("error migration seeder: %v", err)
		}
		log.Println("seeder completed successfully!")
	}
}
