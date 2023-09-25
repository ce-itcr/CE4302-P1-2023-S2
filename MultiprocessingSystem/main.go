package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"strings"
    "strconv"

	"MultiprocessingSystem/components/CacheController"
	"MultiprocessingSystem/components/ProcessingElement"
	"MultiprocessingSystem/utils"
)

func main() {
	// Create termination channel to signal termination
	terminate := make(chan struct{})

	// Create WaitGroup for PEs and CCs
	var wg sync.WaitGroup

	// Declare the Communication Channels for PE-CC
	RequestChannels := make([]chan utils.Request, 3)
	ResponseChannels := make([]chan utils.Response, 3)

	// Create and start 3 Cache Controllers
	cacheControllers := make([]*CacheController.CacheController, 3) // Create an array of Cache Controllers

	for i := 0; i < 3; i++ {
		requestChannel := make(chan utils.Request)
		responseChannel := make(chan utils.Response)

		cacheController, err := CacheController.New(requestChannel, responseChannel, terminate)
		if err != nil {
			fmt.Printf("Error initializing CacheController %d: %v\n", i+1, err)
			return
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			cacheController.Run(&wg)
		}()

		cacheControllers[i] = cacheController
		RequestChannels[i] = requestChannel
		ResponseChannels[i] = responseChannel
	}

	// Create and start 3 Processing Elements
	pes := make([]*processingElement.ProcessingElement, 3) // Create an array of PEs

	for i := 0; i < 3; i++ {
		peName := fmt.Sprintf("PE%d", i+1)

		pe, err := processingElement.New(i+1, peName, RequestChannels[i], ResponseChannels[i], fmt.Sprintf("programs/program%d.txt", i+1), terminate)
		if err != nil {
			fmt.Printf("Error initializing ProcessingElement %d: %v\n", i+1, err)
			return
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			pe.Run(&wg)
		}()

		pes[i] = pe
	}

	// Create a simple command-line interface for controlling your program
	fmt.Println("WELCOME TO MCKEVINHO CLI")
	fmt.Println("The available commands are:")
	fmt.Println("1. step <PE> - Send the Control signal to a specific PE (e.g., 'step 1' or 'step all')")
	fmt.Println("2. lj        - Terminate the program")

	reader := bufio.NewReader(os.Stdin)
PELoop:
	for {
		fmt.Print("\nEnter a command: ")
		command, _ := reader.ReadString('\n')
		command = trimNewline(command)

		args := strings.Split(command, " ")
		if len(args) < 1 {
			fmt.Println("Invalid command. Please enter 'step <PE>' or 'lj'.")
			continue
		}

		switch args[0] {
		case "step":
			if len(args) != 2 {
				fmt.Println("Invalid 'step' command. Please use 'step <PE>' or 'step all'.")
				continue
			}

			if args[1] == "all" {
				for i, pe := range pes {
					if !pe.IsDone && !pe.IsExecutingInstruction {
						pe.Control <- true
						fmt.Printf("Sent 'step' command to %s %d\n", pe.Name, i)
					} else {
						fmt.Printf("%s is not available\n", pe.Name)
					}
				}
			} else {
				peIndex, err := strconv.Atoi(args[1])
				if err != nil || peIndex < 1 || peIndex > len(pes) {
					fmt.Println("Invalid PE number. Please enter a valid PE number or 'all'.")
					continue
				}

				pe := pes[peIndex-1]
				if !pe.IsDone && !pe.IsExecutingInstruction {
					pe.Control <- true
					fmt.Printf("Sent 'step' command to %s\n", pe.Name)
				} else {
					fmt.Printf("%s is not available\n", pe.Name)
				}
			}
		case "lj":
			// Signal termination to both components
			fmt.Println("Sent 'lj' command to terminate the program")
			close(terminate)

			wg.Wait() // Wait for all goroutines to finish gracefully

                // Close the log files for all PEs
            for _, pe := range pes {
                pe.Logger.Writer().(*os.File).Close()
            }

			for i := 0; i < 3; i++ {
				close(RequestChannels[i])
				close(ResponseChannels[i])
			}
			break PELoop
		default:
			fmt.Println("Invalid command. Please enter 'step <PE>' or 'lj'.")
		}
	}
}

func trimNewline(s string) string {
	return s[:len(s)-1]
}