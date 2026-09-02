package cmd

import (
	"fmt"
	"log"
	"os"

	migrations "github.com/Pemilu-CSSMoRA/pemilu26-backend/migration"
	"gorm.io/gorm"
)

func Commands(db *gorm.DB) {
	migrate := false
	seed := false

	for _, arg := range os.Args[1:] {
		if arg == "--migrate" {
			migrate = true
		}
		if arg == "--seed" {
			seed = true
		}

		if arg == "--help" || arg == "-h" {
			printHelp()
			return
		}
	}

	if migrate {
		if err := migrations.Migrate(db); err != nil {
			log.Fatalf("error migration: %v", err)
		}
		log.Println("migration completed successfully")
	}

	if seed {
		if err := migrations.RunSeeder(db); err != nil {
			log.Fatalf("error migration seeder: %v", err)
		}
		log.Println("seeder completed successfully")
	}

}

func printHelp() {
	fmt.Println("==================================================")
	fmt.Println("      Schematics 26 Backend CLI Help Guide")
	fmt.Println("==================================================")
	fmt.Println("Berikut adalah daftar command yang tersedia:")
	fmt.Println("")
	fmt.Println("1. Menjalankan Server")
	fmt.Println("   - go run main.go           : Menjalankan server (default)")
	fmt.Println("")
	fmt.Println("2. Database Migrations & Seeding")
	fmt.Println("   - go run main.go --migrate : Menjalankan migrasi database (GORM AutoMigrate)")
	fmt.Println("     -> Make shortcut         : make migrate")
	fmt.Println("   - go run main.go --seed    : Menjalankan semua seeder utama")
	fmt.Println("     -> Make shortcut         : make seed")
	fmt.Println("")
	fmt.Println("3. Lain-lain")
	fmt.Println("   - go run main.go --help              : Menampilkan panduan ini")
	fmt.Println("     -> Shortcut                        : go run main.go -h")
	fmt.Println("==================================================")
}
