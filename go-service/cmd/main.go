package main

import (
	"log"

	router "github.com/SLANGERES/go-service/internal/Routers"
)

func StartServer(addr string) {
	r := router.Router()
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func main() {
	StartServer(":9090")
}
