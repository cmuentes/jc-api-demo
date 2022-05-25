package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

type WebServer struct {
	http.Server
	shutdownReq chan bool
	reqCount    uint32
}

func CreateServer() *WebServer {
	//create server
	s := &WebServer{
		Server: http.Server{
			Addr:         ":8080",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		shutdownReq: make(chan bool),
	}

	return s
}

func (s *WebServer) WaitShutdown() {
	log.Printf("Stopping http server ...")

	//Create shutdown context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	//shutdown the server
	err := s.Shutdown(ctx)
	if err != nil {
		log.Printf("Shutdown request error: %v", err)
	} else {
		log.Printf("Server shut down successfully ...")
	}
}
