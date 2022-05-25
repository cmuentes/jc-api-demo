package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var HashRequests = make(map[int]HashResponse)
var Mutex sync.RWMutex

func CreateHashPasswordRequest(writer http.ResponseWriter, request *http.Request) {
	numRequests := len(HashRequests)
	numRequests++

	var req HashRequest
	decoder := json.NewDecoder(request.Body)
	decoder.Decode(&req)

	Mutex.Lock()

	HashRequests[numRequests] = HashResponse{
		numRequests,
		time.Now(),
		time.Time{},
		req.Password,
		"",
	}

	Mutex.Unlock()

	writer.Header().Add("Content-type", "application/json")
	writer.WriteHeader(http.StatusAccepted)
	json.NewEncoder(writer).Encode(fmt.Sprintf(`{requestNumber: %d}`, (numRequests)))
}

func GetHashStats(writer http.ResponseWriter, request *http.Request) {
	count := len(HashRequests)
	var avg int64

	for _, val := range HashRequests {
		avg += val.initiatedOn.Sub(val.hashedOn).Microseconds()
	}

	avg = (avg / int64(count))

	writer.Header().Add("Content-type", "application/json")
	json.NewEncoder(writer).Encode(fmt.Sprintf(`{total: %d, average: %d}`, count, avg))
}

func GetHashedPassword(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	reqNum := params["requestNum"]

	num, parsed := strconv.Atoi(reqNum)

	if parsed != nil {
		panic(parsed)
	}

	item, ok := HashRequests[num]

	if ok {
		t := item.initiatedOn.Add(time.Second * 5)
		diff := time.Now().Sub(t)

		if diff.Seconds() >= 5 {
			Mutex.Lock()
			item.hashedOn = time.Now()
			var salt = generateRandomSalt(10)
			var hash = createPasswordHash(item.clearPassword, salt)
			item.hashedPassword = hash

			writer.Header().Add("Content-type", "application/json")
			json.NewEncoder(writer).Encode(fmt.Sprintf(`{hashedPassword: %s}`, item.hashedPassword))

			Mutex.Unlock()
		} else {
			writer.WriteHeader(http.StatusAccepted)
		}
	}
}

func createPasswordHash(password string, salt []byte) string {
	var passwordBytes = []byte(password)
	var sha512Hasher = sha512.New()

	passwordBytes = append(passwordBytes, salt...)
	sha512Hasher.Write(passwordBytes)
	var hashedPasswordBytes = sha512Hasher.Sum(nil)

	var base64EncodedPasswordHash = base64.URLEncoding.EncodeToString(hashedPasswordBytes)

	return base64EncodedPasswordHash
}

func generateRandomSalt(saltSize int) []byte {
	var salt = make([]byte, saltSize)

	_, err := rand.Read(salt[:])

	if err != nil {
		panic(err)
	}

	return salt
}
