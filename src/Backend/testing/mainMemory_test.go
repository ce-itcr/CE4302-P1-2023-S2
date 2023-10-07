package testing

import (
	"testing"
	"sync"
	"time"
	"math/rand"
	"fmt"
	"Backend/utils"
	"Backend/components/MainMemory"
)

// Test that the Main Memory requests and responses simulating requests from an Interconnect
func TestMainMemoryRun(t *testing.T) {
	fmt.Println("Starting Unit Test for the Main Memory Component")
	// Create a wait group for the threads
	var wg sync.WaitGroup
	
	// Initialize a new Processing Element
	quit := make(chan struct{})

	// Declare the Communication Channels for the Interconnect and Main Memory
	requestChannel := make(chan utils.RequestMainMemory)
	responseChannel:= make(chan utils.ResponseMainMemory)
	mm, err := mainMemory.New(requestChannel, responseChannel, "../logs/MM/", quit)
	if err != nil {
		t.Fatalf("Error creating Main Memory: %v", err)
	}
	// Start the Main Memory as a Thread
	wg.Add(1)
	go func() {
		defer wg.Done()
		mm.Run(&wg)
	}()

	// Start a thread with the simulation of the Interconnect requests
	counter := 10
	wg.Add(1)
	go func() {
		defer wg.Done()
		operation := "READ"
		for counter > 0 {
			// Create a struct for the main memory requests
			request := utils.RequestMainMemory{
				Type: operation,
				Address: rand.Intn(16),
				Value: uint32(1),
			}
			// Send the request to the main memory
			requestChannel <- request
			// Wait for a response
			<- responseChannel
			// Switch operation
			switch operation {
			case "READ":
				operation = "WRITE"
			case "WRITE":
				operation = "READ"
			}
			// Decrease the counter
			counter--
		}
	}()

	// Validate if all the requests are managed until the end
	timeout := time.After(1 * time.Minute) // Set 1 minute of timeout
	for counter > 0 {
		select {
		case <-timeout:
			t.Fatal("Test timed out")

		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
	// Close everything
	close(quit)
	wg.Wait()
	close(requestChannel)
	close(responseChannel)
}
