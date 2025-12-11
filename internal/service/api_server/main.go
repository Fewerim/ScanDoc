package main

import (
	"fmt"
	"log"
	hand "proWeb/internal/service/api_server/handlers"
	serv "proWeb/internal/service/api_server/server"
	"time"
)

func main() {
	handlers := hand.NewHTTPHandler()
	server := serv.NewServer(handlers)

	go func() {
		if err := serv.StartServer(server); err != nil {
			log.Fatalf("Error starting server: %s", err)
		}
	}()
	log.Println("Server started")

	time.Sleep(15 * time.Second)
	fmt.Println("eeee")
	if err := server.Stop(); err != nil {
		log.Fatalf("Error stopping server: %s", err)
	}
	log.Println("Server stopped")
	time.Sleep(15 * time.Second)
	fmt.Println("eeee")
}
