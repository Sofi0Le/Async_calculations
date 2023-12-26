package main

import (
	"remote_calculations/internal/api"
	"net/http"
)

func main() {
	print("started\n")
	http.HandleFunc("/api/calculate_operations/", api.Calculate)
	http.ListenAndServe("localhost:8080", nil)
}