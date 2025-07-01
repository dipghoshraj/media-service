package main

import (
	"context"
	"fmt"
	"gateway-tunnel/internal/client"
	"gateway-tunnel/internal/serve"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", serve.HandleRequest)
	http.HandleFunc("/ws", client.WebsocketHandler)

	server := &http.Server{
		Addr:    ":8082",
		Handler: mux,
	}

	go func() {
		log.Println("Server started at :8082")
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	go func() {
		fmt.Println("Server started at :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to gracefully shutdown: %v\n", err)
	}

	log.Println("Server exited properly")

}
