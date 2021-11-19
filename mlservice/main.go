package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
)

const (
	portEnvVar = "PORT"
)

var (
	port = "8081"
)

type machineLearningServiceResponse struct {
	Question string `json:"question"`
	Response string `json:"response"`
}

type machineLearningServiceRequest struct {
	Question string `json:"question"`
}

func init() {
	definedPort := os.Getenv(portEnvVar)
	if definedPort != "" {
		port = definedPort
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("handling machine learning request")
	var request machineLearningServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errMessage := fmt.Errorf("failed to read request body, error %v", err)
		log.Println(errMessage)
		w.Write([]byte(errMessage.Error()))
		return
	}
	log.Printf("%+v\n", request)
	mlResponse := machineLearningServiceResponse{
		Question: request.Question,
		Response: doMachineLearningProcess(),
	}
	b, err := json.Marshal(mlResponse)
	if err != nil {
		errMessage := fmt.Errorf("failed to marshal machine learning response, error %v", err)
		w.Write([]byte(errMessage.Error()))
		return
	}
	w.Write(b)
}

func doMachineLearningProcess() string {
	min := 0
	max := 2
	output := rand.Intn(max-min) + min
	if output == 0 {
		return "yes"
	}
	return "no"
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("listening at port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
