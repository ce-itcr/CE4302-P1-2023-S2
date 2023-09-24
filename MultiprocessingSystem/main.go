package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"MultiprocessingSystem/components/CacheController"
	"MultiprocessingSystem/components/ProcessingElement"
	"MultiprocessingSystem/utils"
)

func main() {
	// Create channels for communication
	requestChannelCC := make(chan utils.Request)
	responseChannelCC := make(chan utils.Response)

	// Instantiate CacheController
	cacheController, err := CacheController.New(requestChannelCC, responseChannelCC)
	if err != nil {
		fmt.Printf("Error initializing CacheController: %v\n", err)
		return
	}

	// Instantiate ProcessingElement
	processingElement, err := processingElement.New(1, "PE1", requestChannelCC, responseChannelCC, "programs/program1.txt")
	if err != nil {
		fmt.Printf("Error initializing ProcessingElement: %v\n", err)
		return
	}

	// Use WaitGroup to synchronize goroutines
	var wg sync.WaitGroup

	// Start CacheController in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		cacheController.Run(&wg)
	}()

	// Start ProcessingElement in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		processingElement.Run(&wg)
	}()

	// Create a simple command-line interface for controlling your program
	fmt.Println("Welcome to your program CLI!")
	fmt.Println("Available commands:")
	fmt.Println("1. step - Send the Control signal to PE")
	fmt.Println("2. lj   - Terminate the program")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter a command: ")
		command, _ := reader.ReadString('\n')
		command = trimNewline(command)

		switch command {
		case "step":
			// Sending a control signal to ProcessingElement
			processingElement.Control <- true
			fmt.Println("Sent 'step' command to ProcessingElement")
		case "lj":
			// Perform cleanup here if needed

			// Signal the ProcessingElement to finish
			processingElement.Done <- true
			fmt.Println("Sent 'lj' command to terminate the program")
			// Wait for both goroutines to finish
			wg.Wait()
			// Close channels when done to avoid goroutine leaks
			close(requestChannelCC)
			close(responseChannelCC)
			return
		default:
			fmt.Println("Invalid command. Please enter 'step' or 'lj'.")
		}
	}
}

func trimNewline(s string) string {
	return s[:len(s)-1]
}
