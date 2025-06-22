package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/ws", websocketHandler)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
