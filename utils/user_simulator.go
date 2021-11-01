package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"

	"time"

	"github.com/google/uuid"
	"github.com/unullmass/msg-store/models"
)

var (
	s        = rand.NewSource(time.Now().UnixNano())
	r        = rand.New(s)
	createWg sync.WaitGroup
)

const (
	createUrl       = "http://localhost:8080/mydata"
	retrieveUrlBase = "http://localhost:8080/mydata/document"
	maxBurst        = 10
)

// documentCreateRequest is used to hold the unmarshalled DocumentCreate request payload
// The fields extracted from this are passed on to the Document model
type documentCreateRequest struct {
	Id        uuid.UUID
	Attrs     []models.Attribute
	Timestamp string
}

type documentSearchRequest struct {
	Timestamp string
	Key       string
	Value     string
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
			//randStringBytes(r.Intn(10)),
			//randStringBytes(r.Intn(20)))
			fmt.Sprintf("%s%d", "key", i+1), fmt.Sprintf("%s%d", "value", i+1))
		if err != nil {
			log.Default().Printf("Skipping failed records %+v\n", err)
			continue
		}
		attrs = append(attrs, *na)
	}
	return attrs
}

func main() {

	n := r.Intn(10)
	log.Default().Printf("Generating %d new records\n", n)
	burstCount := 0
	for i := 0; i < n; i++ {
		if burstCount < maxBurst {
			createWg.Add(1)
			burstCount++
			go func() {
				defer createWg.Done()
				burstCount++

				dcrsp, err := DoCreate()
				if err != nil {
					log.Default().Printf("%+v\n", err)
					return
				}

				if dcrsp.Id == nil {
					log.Default().Print("Null ID\n", err)
					return
				}

				// retrieve document
				docRetRsp, err := DoRetrieve(*dcrsp.Id)
				if err != nil {
					log.Default().Printf("%+v\n", err)
					return
				}

				// perform a search to validate create
				var dsreq documentSearchRequest
				dsreq.Timestamp = fmt.Sprint(*docRetRsp.Timestamp)
				dsreq.Key = docRetRsp.Attrs[0].Key
				dsreq.Value = docRetRsp.Attrs[0].Value
				dsresp, err := DoSearch(&dsreq)
				if err != nil {
					log.Default().Printf("%+v\n", err)
					return
				}
				log.Default().Printf("Search returned %d records\n", len(dsresp.Docs))
			}()
		} else {
			createWg.Wait()
			burstCount = 0
		}
	}
	createWg.Wait()
}

func DoRetrieve(id uuid.UUID) (*retrieveDocumentResponse, error) {
	retrieveUrl := fmt.Sprintf("%s/%s", retrieveUrlBase, id)
	c := http.DefaultClient
	resp, err := c.Get(retrieveUrl)
	log.Default().Printf("Retrieve Request: %+v\n", retrieveUrl)

	if err != nil {
		log.Default().Printf("Retrieve Error: %+v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Default().Printf("Retrieve Error: %+v\n", err)
		return nil, err
	}
	var drresp retrieveDocumentResponse
	err = json.Unmarshal(body, &drresp)
	if err != nil {
		log.Default().Printf("Retrieve Error: %+v\n", err)
		drresp.Status = err.Error()
	}

	log.Default().Printf("Retrieve Response: %+v\n", drresp)
	return &drresp, nil
}

func DoCreate() (*documentCreateResponse, error) {
	dcr := documentCreateRequest{
		Id:        uuid.New(),
		Attrs:     generateRandomAttrs(),
		Timestamp: fmt.Sprint(time.Now().Unix()),
	}
	log.Default().Printf("Create Request: %+v\n")
	json.NewEncoder(os.Stdout).Encode(dcr)

	reqBytes, _ := json.Marshal(dcr)
	c := http.DefaultClient
	resp, err := c.Post(createUrl, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Default().Printf("Create Error: %+v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Default().Printf("Create Error: %+v\n", err)
		return nil, err
	}
	dcresp := documentCreateResponse{}
	err = json.Unmarshal(body, &dcresp)
	if err != nil {
		log.Default().Printf("Create Error: %+v\n", err)
		return nil, err
	}

	log.Default().Printf("Response: %+v\n", dcresp)
	return &dcresp, nil
}

func DoSearch(docSearchReq *documentSearchRequest) (*documentSearchResponse, error) {
	var docSearchResp documentSearchResponse
	searchUrl := fmt.Sprintf("%s/%s/%s/%s", createUrl, docSearchReq.Timestamp, docSearchReq.Key, docSearchReq.Value)
	c := http.DefaultClient
	resp, err := c.Get(searchUrl)
	log.Default().Printf("Search Request: %+v\n", searchUrl)

	if err != nil {
		log.Default().Printf("Search Error: %+v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Default().Printf("Search Error: %+v\n", err)
		return nil, err
	}
	err = json.Unmarshal(body, &docSearchResp)
	if err != nil {
		log.Default().Printf("Search Error: %+v\n", err)
		return nil, err
	}

	//log.Default().Printf("Response: %+v\n", dcrsp)
	return &docSearchResp, nil
}
