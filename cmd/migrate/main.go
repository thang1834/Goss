package main

import (
	"log"
	"os"

	"github.com/thang1834/go-goss/config"
	"github.com/thang1834/go-goss/database"
	db "github.com/thang1834/go-goss/third_party/database"
)

// Version is injected using ldflags during build time
var Version string

func main() {
	log.Printf("Version: %s\n", Version)

	cfg := config.New()
	store := db.NewSqlx(cfg.Database)
	migrator := database.Migrator(store.DB)

	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "up":
		log.Println("Running UP migrations")
		migrator.Up()

	case "down":
		log.Println("Running DOWN migration")
		migrator.Down()

	default:
		log.Fatalf("unknown migration command: %q", cmd)
	}
}
