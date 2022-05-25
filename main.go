package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const API_ROOT = "/api"

func main() {
	fmt.Println("Starting server ...")

	//Start the server
	server := CreateServer()

	done := make(chan bool)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Listen and serve: %v", err)
		}
		done <- true
	}()

	log.Printf("DONE!")
	router := mux.NewRouter()

	router.HandleFunc(fmt.Sprintf("%s/hash", API_ROOT), CreateHashPasswordRequest).Methods("POST")
	router.HandleFunc(fmt.Sprintf("%s/hash/{requestNum}", API_ROOT), GetHashedPassword).Methods("GET")
	router.HandleFunc(fmt.Sprintf("%s/stats", API_ROOT), GetHashStats).Methods("GET")
	router.HandleFunc(fmt.Sprintf("%s/shutdown", API_ROOT), func(writer http.ResponseWriter, request *http.Request) {
		OnShuttingDown(true)
		writer.WriteHeader(http.StatusAccepted)
		json.NewEncoder(writer).Encode(fmt.Sprintf(`{shuttingDown: %t}`, (true)))

		go func() {
			var timer = time.NewTimer(time.Second * 5)
			<-timer.C
			server.WaitShutdown()
		}()
	})

	server.Handler = router

	log.Fatal(http.ListenAndServe(":8000", router))
}
