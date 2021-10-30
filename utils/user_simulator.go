package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"time"

	"github.com/google/uuid"
	"github.com/unullmass/msg-store/models"
)

var (
	s = rand.NewSource(time.Now().UnixNano())
	r = rand.New(s)
)

// documentCreateRequest is used to hold the unmarshalled DocumentCreate request payload
// The fields extracted from this are passed on to the Document model
type documentCreateRequest struct {
	Id        uuid.UUID
	Attrs     []models.Attribute
	Timestamp string
}

// documentCreateResponse is used to hold the response payload for Create Document flow
type documentCreateResponse struct {
	Id     *uuid.UUID
	Status string
}

// retrieveDocumentResponse is used to hold the retrieve document response payload
type retrieveDocumentResponse struct {
	Id        *uuid.UUID
	Attrs     []models.Attribute
	Timestamp *int64
	Status    string
}

// documentSearchResponse is used to hold the search document response payload
type documentSearchResponse struct {
	Docs []uuid.UUID
}

func randStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func generateRandomAttrs() []models.Attribute {
	attrs := []models.Attribute{}
	for i := 0; i < r.Intn(100)%10; i++ {
		na, err := models.NewAttribute(
			randStringBytes(r.Intn(256)),
			randStringBytes(r.Intn(256)))
		if err != nil {
			log.Default().Printf("Skipping failed records %+v\n", err)
			continue
		}
		attrs = append(attrs, *na)
	}
	return attrs
}

const (
	createUrl = "http://localhost:8080/mydata"
	maxBurst  = 50
)

var wg sync.WaitGroup

func main() {

	n := r.Intn(100)
	log.Default().Printf("Generating %d new records\n", n)
	burstCount := 0
	for i := 0; i < n; i++ {

		if burstCount < maxBurst {
			wg.Add(1)
			go DoCreate()
			burstCount++
		} else {
			wg.Wait()
			burstCount = 0
		}

	}
	wg.Wait()
}

func DoCreate() {
	defer wg.Done()

	dcr := documentCreateRequest{
		Id:        uuid.New(),
		Attrs:     generateRandomAttrs(),
		Timestamp: fmt.Sprint(time.Now().Unix()),
	}
	log.Default().Printf("Request: %+v\n", dcr)

	reqBytes, _ := json.Marshal(dcr)
	c := http.DefaultClient
	resp, err := c.Post(createUrl, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Default().Printf("Error: %+v\n", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Default().Printf("Error: %+v\n", err)
	}
	var dcrsp documentCreateResponse
	json.Unmarshal(body, &dcrsp)
	//log.Default().Printf("Response: %+v\n", dcrsp)
}
