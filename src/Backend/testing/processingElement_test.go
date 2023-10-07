package testing

import (
	"testing"
	"fmt"
	"sync"
	"time"
	"os"
	"Backend/utils"
	"Backend/components/ProcessingElement"
)

// Test that the Processing Element the correct operations depending on the instruction type
func TestProcessingElementRun(t *testing.T) {
	fmt.Println("Starting Unit Test for the Processing Element Component")
	// Create a new random program with 50 random instructions to execute
	instructions := utils.GenerateRandomInstructions(1, 50)
	filename := fmt.Sprintf("../generated-programs/program0.txt")
	err := utils.WriteInstructionsToFile(filename, instructions[0])
	if err != nil {
		fmt.Printf("Could not create new instructions")
	}
	// Create a wait group for the threads
	var wg sync.WaitGroup
	
	// Initialize a new Processing Element
	quit := make(chan struct{})
	requestChannel := make(chan utils.RequestProcessingElement)
	responseChannel := make(chan utils.ResponseProcessingElement)
	pe, err := processingElement.New(0, requestChannel, responseChannel, "../generated-programs/program0.txt", "../logs/PE/PE", quit)
	if err != nil {
		t.Fatalf("Error creating ProcessingElement: %v", err)
	}
	// Start the Processing Element as a Thread
	wg.Add(1)
	go func() {
		defer wg.Done()
		pe.Run(&wg)
	}()

	// Start a thread with the simulation of the private Cache Controller responses
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <- quit:
				return
			case <-requestChannel:
				// Prepare a structure to respond to the Processing Element's request
				response := utils.ResponseProcessingElement {
					Data: 10,
					Status: true,
				}
				// Send the response
				responseChannel <- response
			}
		}
	}()

	done := false
	timeout := time.After(1 * time.Minute) // Set 1 minute of timeout
	ticker := time.NewTicker(300 * time.Millisecond) 
	
	for !done {
		select {
		case <-timeout:
			t.Fatal("Test timed out")
		case <-ticker.C:
			if !pe.IsDone && !pe.IsExecutingInstruction {
				// Send a control signal to the Processing Element
				pe.Control <- true
			}
			if pe.IsDone {
				done = true
				// Close everything
				close(quit)
				wg.Wait()
				close(requestChannel)
				close(responseChannel)
				pe.Logger.Writer().(*os.File).Close()
				break
			}
		}
	}
	// Stop the ticker when done
	ticker.Stop()
}
