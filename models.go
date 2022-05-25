package main

import (
	"time"
)

type HashResponse struct {
	requestNumber  int
	initiatedOn    time.Time
	hashedOn       time.Time
	clearPassword  string
	hashedPassword string
}

type HashRequest struct {
	Password string `json:"password"`
}
