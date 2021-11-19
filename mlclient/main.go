package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	portEnvVar                       = "PORT"
	machineLearningServiceHostEnvVar = "ML_SERVICE_HOST"
)

var (
	port                       = "8080"
	machineLearningServiceHost = "http://localhost:8081"
)

type machineLearningServiceResponse struct {
	Question string `json:"question"`
	Response string `json:"response"`
}

type machineLearningServiceRequest struct {
	Question string `json:"question"`
}

func init() {
	mlHost := os.Getenv(machineLearningServiceHostEnvVar)
	if mlHost != "" {
		machineLearningServiceHost = mlHost
	}
	definedPort := os.Getenv(portEnvVar)
	if definedPort != "" {
		port = definedPort
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	question := r.URL.Query().Get("question")
	if question == "" {
		w.Write([]byte("ask something"))
		return
	}
	log.Printf("question: %s\n", question)
	request := machineLearningServiceRequest{
		Question: question,
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		errMessage := fmt.Errorf("failed to marshal request, error %v", err)
		log.Println(errMessage)
		w.Write([]byte(errMessage.Error()))
		return
	}
	log.Printf("request body: %s\n", string(requestBytes))
	log.Println("caling machine learning service at " + machineLearningServiceHost)
	res, err := http.Post(machineLearningServiceHost, "application/json", bytes.NewBuffer(requestBytes))
	if err != nil {
		errMessage := fmt.Errorf("failed to request machine learning service at host [%s], error %v", machineLearningServiceHost, err)
		log.Println(errMessage)
		w.Write([]byte(errMessage.Error()))
		return
	}
	defer res.Body.Close()
	var mlResponse machineLearningServiceResponse
	bodyAsBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		errMessage := fmt.Errorf("failed to ready content of body, error %v", err)
		log.Println(errMessage)
		w.Write([]byte(errMessage.Error()))
		return
	}
	if err := json.Unmarshal(bodyAsBytes, &mlResponse); err != nil {
		errMessage := fmt.Errorf("failed to unmarshal ml response, error %v", err)
		log.Println(errMessage)
		w.Write([]byte(errMessage.Error()))
		return
	}
	log.Printf("%+v\n", mlResponse)
	w.Write(bodyAsBytes)
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("listening at port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
