package testing

import (
	"testing"
	"fmt"
	"sync"
	"math/rand"
	"time"
	"Backend/utils"
	"Backend/components/Interconnect"
)

// Test the Cache Controller communication with a Processing Element and Interconnect
func TestInterconnectRun(t *testing.T) {
	fmt.Println("Starting Unit Test for the Interconnect Component")
	
	// Create the communication channels for the Cache Controller and the Interconnect
	requestChannelsInterconnect := make([]chan utils.RequestInterconnect, 3)
	responseChannelsInterconnect := make([]chan utils.ResponseInterconnect, 3)

	// Create the communication channel with broadcast messaging in the bus
	requestChannelsBroadcast := make([]chan utils.RequestBroadcast, 3)
	responseChannelsBroadcast := make([]chan utils.ResponseBroadcast, 3)

	// Declare the Communication Channels for the Interconnect and Main Memory
	requestChannelMainMemory := make(chan utils.RequestMainMemory)
	responseChannelMainMemory:= make(chan utils.ResponseMainMemory)


	// Define the cache coherence protocol
	protocol := "MOESI"
	quit := make(chan struct{})

	// Create a wait group for the threads
	var wg sync.WaitGroup

	// Create 3 threats simulating the Cache Controllers
	// Create the Bus semaphore
	semaphore := make(chan struct{}, 1)
	cacheControllersDone := [3]bool{false, false, false}
	for i := 0; i < 3; i++ {
		requestChannelInterconnect := make(chan utils.RequestInterconnect)
		responseChannelInterconnect := make(chan utils.ResponseInterconnect)
		requestChannelBroadcast := make(chan utils.RequestBroadcast)
		responseChannelBroadcast := make(chan utils.ResponseBroadcast)
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// request message structure for Interconnect
			requestIC := utils.RequestInterconnect{
				Type: "",
				AR: "",
				Address: 0,
				Data: 0,            
			}
			// response broadcast structure
			responseBC :=  utils.ResponseBroadcast {
				Match: false,
				Status: "I",
				Data: 0,
			}

			// Execute a thread to simulate the Cache Controller requests to the Interconnect
			go func(i int) {
				instructions := utils.GenerateRandomInstructions(1, 0)[0]
				for _, item := range instructions {
					select {
					case semaphore <- struct{}{}:
						operation := item.Type
						address := item.Address
						switch operation {
						case "READ":
							requestIC = utils.RequestInterconnect {
								Type: "ReadRequest",
								AR: "DataResponse",
								Address: address,
								Data: rand.Intn(100),
							}
							// send the request to the Interconnect
							requestChannelInterconnect <- requestIC
							// Wait for a response
							<- responseChannelInterconnect
							// Release the semaphore
							<- semaphore
							time.Sleep(time.Second*2)
						case "WRITE":
							requestIC = utils.RequestInterconnect {
								Type: "ReadExclusiveRequest",
								AR: "DataResponse",
								Address: address,
								Data: rand.Intn(100),
							}
							// send the request to the Interconnect
							requestChannelInterconnect <- requestIC
							// Wait for a response
							<- responseChannelInterconnect
							// Release the semaphore
							<- semaphore
							time.Sleep(time.Second*2)

						case "INC":
							// Release the semaphore
							<- semaphore
							time.Sleep(time.Second*2)
						}
					}
				}
				// Update the cache controller done status
				cacheControllersDone[i] = true
			}(i)
			
			// Create a thread to simulate the broadcast responses
			go func(i int) {
				for {
					select {
					case request := <- requestChannelBroadcast:
						requestType := request.Type
						switch requestType {
						case "ReadRequest":
							// response broadcast structure
							responseBC =  utils.ResponseBroadcast {
								Match: false,
								Status: "I",
								Data: 0,
							}
							// send the response back to the Interconnect
							responseChannelBroadcast <- responseBC
		
						case "ReadExclusiveRequest":
							// response broadcast structure
							responseBC =  utils.ResponseBroadcast {
								Match: true,
								Status: "M",
								Data: 0,
							}
							// send the response back to the Interconnect
							responseChannelBroadcast <- responseBC
						}
					case <- quit:
						return
					}
				}
			}(i)
		}(i)

		// Add the communication channels to the arrays
		requestChannelsInterconnect[i] = requestChannelInterconnect
		responseChannelsInterconnect[i] = responseChannelInterconnect
		requestChannelsBroadcast[i] = requestChannelBroadcast
		responseChannelsBroadcast[i] = responseChannelBroadcast
	}
	// Simulate the Main Memory responses
	wg.Add(1)
	go func(){
		defer wg.Done()
		// Initialize a struct to pack the response
		response := utils.ResponseMainMemory{
			Status: false,
			Value:  0,
		}
		for {
			select {
			case <- requestChannelMainMemory:
				// prepare the response to send back to the Interconnect
				response = utils.ResponseMainMemory {
					Status: true,
					Value: 100,
				}
				// Send the response
				responseChannelMainMemory <- response
			case <- quit:
				return
			}
		}
	}()

	// Create Interconnect
	ic, err := interconnect.New(
		requestChannelsInterconnect,
		responseChannelsInterconnect,
		requestChannelMainMemory,
		responseChannelMainMemory,
		requestChannelsBroadcast,
		responseChannelsBroadcast,
		protocol,
		"../logs/IC/",
		quit,
	)
	if err != nil {
		t.Fatalf("Error initializing Interconnect: %v\n", err)
	}
	// Start Interconnect
	wg.Add(1)
	go func() {
		defer wg.Done()
		ic.Run(&wg)
	}()

	// Validate if all the requests are managed until the end
	timeout := time.After(1 * time.Minute) // Set 1 minute of timeout
	for (!cacheControllersDone[0]) && (!cacheControllersDone[1]) && (!cacheControllersDone[2]) {
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
	for i := 0; i < 3; i++ {
		close(requestChannelsInterconnect[i])
		close(responseChannelsInterconnect[i])
		close(requestChannelsBroadcast[i])
		close(responseChannelsBroadcast[i])
	}
	close(requestChannelMainMemory)
	close(responseChannelMainMemory)
	close(semaphore)
}
