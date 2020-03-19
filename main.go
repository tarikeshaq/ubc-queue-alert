package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const DEFAULTCOURSE string = "CPSC 213"

var threshold int64
var course string
var webhook string
var token string
var queueUrl string

type SlackMessage struct {
	Text string `json:"text"`
}

type Course struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Queue struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Course        Course `json:"course"`
	QuestionCount int    `json:"questionCount"`
}

func main() {
	delay := setupEnv()
	interval := time.Duration(delay) * time.Minute
	ticker := time.NewTicker(interval)
	done := make(chan bool)
	go func() {
		for {
			<-ticker.C
			go checkQueue()
		}
	}()
	<-done
}

func setupEnv() int64 {
	thresholdEnv := os.Getenv("THRESHOLD")
	var err error
	threshold, err = strconv.ParseInt(thresholdEnv, 10, 32)
	if err != nil {
		threshold = 10
	}
	course = os.Getenv("COURSE")
	if course == "" {
		course = DEFAULTCOURSE
	}

	webhook = os.Getenv("WEBHOOK")
	if webhook == "" {
		log.Fatal("No webhook provided, please add a webhook to your environment variables")
	}

	token = os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("No Queue token provided, please add a token to your environment variables")

	}

	queueUrl = os.Getenv("URL")
	if queueUrl == "" {
		log.Fatal("No queue URL provided, please add a url to your environment variables")
	}

	delay := os.Getenv("DELAY")
	delayNum, err := strconv.ParseInt(delay, 10, 32)
	if err != nil {
		delayNum = 10
	}
	return delayNum
}

func checkQueue() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/queues", queueUrl), nil)
	if err != nil {
		log.Printf("error creating request: %v\n", err)
		return
	}
	req.Header.Set("Private-Token", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var queues []Queue
	err = decoder.Decode(&queues)
	if err != nil {
		log.Printf("error reading body: %v\n", err)
		return
	}
	queueLength := 0
	for _, queue := range queues {
		if queue.Course.Name == course {
			queueLength += queue.QuestionCount
		}
	}
	if int64(queueLength) > threshold {
		alertSlack(queueLength)
	}
}

func alertSlack(queueLength int) {
	sm := SlackMessage{
		Text: fmt.Sprintf("The queue is now %d long please send help", queueLength),
	}
	body, err := json.Marshal(sm)
	if err != nil {
		log.Printf("Error unmarshalling message %v\n", err)
	}
	req, err := http.NewRequest("POST", webhook, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Unable to parse request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	client.Do(req)
}
