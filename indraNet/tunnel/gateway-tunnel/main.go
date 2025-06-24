package main

import (
	"fmt"
	"gateway-tunnel/internal/client"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/ws", client.WebsocketHandler)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
