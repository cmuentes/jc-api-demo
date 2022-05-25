package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const API_ROOT = "/api"

func main() {
	fmt.Println("Starting server ...")

	router := mux.NewRouter()

	router.HandleFunc(fmt.Sprintf("%s/hash", API_ROOT), CreateHashPasswordRequest).Methods("POST")
	router.HandleFunc(fmt.Sprintf("%s/hash/{requestNum}", API_ROOT), GetHashedPassword).Methods("GET")
	router.HandleFunc(fmt.Sprintf("%s/stats", API_ROOT), GetHashStats).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
