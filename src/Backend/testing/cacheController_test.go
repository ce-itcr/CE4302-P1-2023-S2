package testing

import (
	"testing"
	"fmt"
	"sync"
	"math/rand"
	"time"
	"Backend/utils"
	"Backend/components/CacheController"
)

// Test the Cache Controller communication with a Processing Element and Interconnect
func TestCacheControllerRun(t *testing.T) {
	fmt.Println("Starting Unit Test for the Cache Controller Component")
	
	// Create the communication channels for a Processing Element and a Cache Controller
	requestChannelProcessingElement := make(chan utils.RequestProcessingElement)
	responseChannelProcessingElement := make(chan utils.ResponseProcessingElement)

	// Create the communication channels for the Cache Controller and the Interconnect
	requestChannelInterconnect := make(chan utils.RequestInterconnect)
	responseChannelInterconnct := make(chan utils.ResponseInterconnect)

	// Create the communication channel with broadcast messaging in the bus
	requestChannelBroadcast := make(chan utils.RequestBroadcast)
	responseChannelBroadcast := make(chan utils.ResponseBroadcast)

	// Create the Bus semaphore
	semaphore := make(chan struct{}, 1)

	// Define the cache coherence protocol
	protocol := "MESI"
	quit := make(chan struct{})

	// Create a wait group for the threads
	var wg sync.WaitGroup

	// Create an isolated Cache Controller
	cc, err := CacheController.New(
		0,
		requestChannelProcessingElement,
		responseChannelProcessingElement,
		requestChannelInterconnect,
		responseChannelInterconnct,
		requestChannelBroadcast,
		responseChannelBroadcast,
		semaphore,
		protocol,
		"../logs/CC/CC",
		quit,
	)
	if err != nil {
		t.Fatalf("Error creating Cache Controller: %v", err)
	}
	// Add the CacheController to the Wait Group
	wg.Add(1)
	go func() {
		defer wg.Done()
		cc.Run(&wg)
	}()

	// Start a thread to simulate Processing Element requests
	peIsDone := false
	wg.Add(1)
	go func() {
		defer wg.Done()
		// request message structure for the Cache Controller
		request := utils.RequestProcessingElement {
				Type: "",
				Address: 0,                    
				Data: 0,                      
		}
		instructions := utils.GenerateRandomInstructions(1, 20)[0]
		for _, item := range instructions {
			operation := item.Type
			address := item.Address
			if operation != "INC" {
				// Send a read request to the Cache Controller
				request = utils.RequestProcessingElement {
					Type: operation,
					Address: address,
					Data: rand.Intn(100),
				}
				requestChannelProcessingElement <- request
				// Wait for a response
				<- responseChannelProcessingElement
			}
		}
		// Update the flag to terminate the main loop
		peIsDone = true
	}()

	// Start a thread to simulate Interconncet responses
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <- quit:
				return
			case <- requestChannelInterconnect:
				// Create a response structure
				response := utils.ResponseInterconnect {
					Data: rand.Intn(100),
					NewStatus: "E",
				}
				// Send the response to the Cache Controller
				responseChannelInterconnct <- response
			}
		}
	}()

	// Validate if all the requests are managed until the end
	timeout := time.After(1 * time.Minute) // Set 1 minute of timeout
	for !peIsDone {
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
	close(requestChannelProcessingElement)
	close(responseChannelProcessingElement)
	close(responseChannelInterconnct)
	close(requestChannelBroadcast)
	close(responseChannelBroadcast)
	close(semaphore)
}
