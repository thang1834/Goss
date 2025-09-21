package main

import (
	"github.com/thang1834/go-goss/internal/server"
)

// Version is injected using ldflags during build time
var Version = "v0.1.0"

// @title GoGoSS
// @version 0.1.0
// @description Go + Postgres + Chi router + sqlx + ent + Testing starter kit for API development.
// @contact.name Thangnd1597
// @contact.url https://github.com/thang1834/go-goss
// @contact.email thangnd1597@gmail.com
// @host localhost:3080
// @BasePath /api/v1
func main() {
	s := server.New(server.WithVersion(Version))
	s.Init()
	s.Run()
}
