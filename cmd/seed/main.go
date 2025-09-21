package main

import (
	"fmt"
	"github.com/thang1834/go-goss/config"
	"github.com/thang1834/go-goss/database"
	db "github.com/thang1834/go-goss/third_party/database"
)

func main() {
	cfg := config.New()
	store := db.NewSqlx(cfg.Database)

	seeder := database.Seeder(store.DB)
	seeder.SeedUsers()
	fmt.Println("seeding completed.")
}
