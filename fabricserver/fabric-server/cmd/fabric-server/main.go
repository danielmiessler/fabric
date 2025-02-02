package main

import (
	"log"

	"github.com/natlamir/fabric-server/internal/server"
)

func main() {
	s, err := server.NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	
	log.Printf("Starting fabric server...")
	if err := s.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
