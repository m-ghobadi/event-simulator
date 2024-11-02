package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type EventType int

const (
	Transaction EventType = iota
	Log
	Notification
	Command
	Query
)

var eventUrgencyLevel = map[EventType]map[int]float64{
	Transaction: {
		1: 0.05, 2: 0.10, 3: 0.15, 4: 0.30, 5: 0.40,
	},
	Log: {
		1: 0.40, 2: 0.30, 3: 0.15, 4: 0.05, 5: 0.10,
	},
	Notification: {
		1: 0.10, 2: 0.20, 3: 0.30, 4: 0.30, 5: 0.10,
	},
	Command: {
		1: 0.05, 2: 0.15, 3: 0.35, 4: 0.35, 5: 0.10,
	},
	Query: {
		1: 0.30, 2: 0.30, 3: 0.25, 4: 0.10, 5: 0.05,
	},
}

var eventRequestCount = map[EventType]int{
	Transaction:  100000,
	Log:          100000,
	Notification: 100000,
	Command:      100000,
	Query:        100000,
}

var eventDelay = map[EventType]time.Duration{
	Transaction:  0 * time.Millisecond, // 500 -> 0.5 seconds
	Log:          0 * time.Millisecond, // 0.25 seconds
	Notification: 0 * time.Millisecond, // 0.4 seconds
	Command:      0 * time.Millisecond, // 0.3 seconds
	Query:        0 * time.Millisecond, // 0.35 seconds
}

func main() {

	var wg sync.WaitGroup
	wg.Add(5)

	// Generate requests concurrently for each event type
	go func() { generateRequests(Transaction); wg.Done() }()
	go func() { generateRequests(Log); wg.Done() }()
	go func() { generateRequests(Notification); wg.Done() }()
	go func() { generateRequests(Command); wg.Done() }()
	go func() { generateRequests(Query); wg.Done() }()

	wg.Wait()

}

func generateRequests(eventType EventType) {
	count := 0
	// Generate requests for the given event type
	for i := 0; i < eventRequestCount[eventType]; i++ {
		count++
		urgencyLevel := getRandomUrgencyLevel(eventType)
		makeRequest(eventType, urgencyLevel)
		time.Sleep(eventDelay[eventType])
		fmt.Println("Generated request for event type:", eventTypeToString(eventType), " Total: :", count)
	}

}

func makeRequest(eventType EventType, urgencyLevel int) {

	eventID := uuid.New().String()
	// url := "http://50.19.22.56"
	// url := "http://127.0.0.1"
	url := "http://127.0.0.1:8181/event"
	// url := "http://192.168.1.104:8181/event"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("X-Event-Type", eventTypeToString(eventType))
	req.Header.Set("X-Event-Urgency", fmt.Sprintf("%d", urgencyLevel))
	req.Header.Set("X-Event-ID", eventID)
	req.Header.Set("X-Event-Request-Time", time.Now().Format(time.RFC3339))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

}

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func getRandomUrgencyLevel(eventType EventType) int {

	r := rnd.Float64()
	urgencyLevelMap := eventUrgencyLevel[eventType]

	for level, percentage := range urgencyLevelMap {
		if r < percentage {
			return level
		}
		r -= percentage
	}

	return 1 // Default to level 1 if no match
}

func eventTypeToString(eventType EventType) string {
	switch eventType {
	case Transaction:
		return "transaction"
	case Log:
		return "log"
	case Notification:
		return "notification"
	case Command:
		return "command"
	case Query:
		return "query"
	default:
		return "unknown"
	}
}
