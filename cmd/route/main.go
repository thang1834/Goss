package main

import (
	"fmt"

	"github.com/thang1834/go-goss/internal/server"
)

func main() {
	s := server.New()
	s.InitDomains()
	fmt.Print("Registered Routes:\n\n")
	s.PrintAllRegisteredRoutes()
}
