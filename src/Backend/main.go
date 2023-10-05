package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"strings"
    "strconv"

	"Backend/components/CacheController"
	"Backend/components/ProcessingElement"
	"Backend/components/Interconnect"
	"Backend/components/MainMemory"
	"Backend/utils"
)



func main() {
	// generate the program files ****************************************************************************************************
	instructionsPerCore := 10
	instructions := utils.GenerateRandomInstructions(3, instructionsPerCore)
	// Write instructions to files
	for coreID, coreInstructions := range instructions {
		filename := fmt.Sprintf("generated-programs/program%d.txt", coreID)
		err := utils.WriteInstructionsToFile(filename, coreInstructions)
		if err != nil {
			fmt.Printf("Error writing to file for Core %d: %v\n", coreID, err)
		} else {
			fmt.Printf("Instructions for Core %d written to %s\n", coreID, filename)
		}
	}


	// Create termination channel to signal the termination to all threads
	terminate := make(chan struct{})

	// Create WaitGroup for PEs and CCs
	var wg sync.WaitGroup

	// Declare the Communication Channels array for PE-CC
	RequestChannelsM1 := make([]chan utils.RequestProcessingElement, 3)
	ResponseChannelsM1 := make([]chan utils.ResponseProcessingElement, 3)

	// Declare the Communication Channels array for CC-IC
	RequestChannelsM2 := make([]chan utils.RequestInterconnect, 3)
	ResponseChannelsM2 := make([]chan utils.ResponseInterconnect, 3)

	// Declare the Broadcast Communication Channels array for CC-IC
	RequestChannelsBroadcast := make([]chan utils.RequestBroadcast, 3)
	ResponseChannelsBroadcast:= make([]chan utils.ResponseBroadcast, 3)

	// Declare the Communication Channels for the Interconnect and Main Memory
	RequestChannelM3 := make(chan utils.RequestMainMemory)
	ResponseChannelM3 := make(chan utils.ResponseMainMemory)

	// Create and start 3 Cache Controllers with the communication channels
	ccs := make([]*CacheController.CacheController, 3) // Create an array of Cache Controllers

	semaphore := make(chan struct{}, 1) // Initialize with a count of 1

	protocol := "MESI"

	for i := 0; i < 3; i++ {
		// Create the Request and Response channels for PE and IC communications
		requestChannelM1 := make(chan utils.RequestProcessingElement)
		responseChannelM1 := make(chan utils.ResponseProcessingElement)

		requestChannelM2 := make(chan utils.RequestInterconnect)
		responseChannelM2 := make(chan utils.ResponseInterconnect)

		requestChannelBroadcast := make(chan utils.RequestBroadcast)
		responseChannelBroadcast := make(chan utils.ResponseBroadcast)

		// Create the CacheController with its ID and communication channels
		cacheController, err := CacheController.New(
				i, 
				requestChannelM1, 
				responseChannelM1, 
				requestChannelM2, 
				responseChannelM2, 
				requestChannelBroadcast, 
				responseChannelBroadcast,
				semaphore,
				protocol,
				terminate)
		if err != nil {
			fmt.Printf("Error initializing CacheController %d: %v\n", i+1, err)
			return
		}

		// Add the CacheController to the Wait Group
		wg.Add(1)
		go func() {
			defer wg.Done()
			cacheController.Run(&wg)
		}()

		// Save the CacheController and the communicatio channels created
		ccs[i] = cacheController
		RequestChannelsM1[i] = requestChannelM1
		ResponseChannelsM1[i] = responseChannelM1

		RequestChannelsM2[i] = requestChannelM2
		ResponseChannelsM2[i] = responseChannelM2
		
		RequestChannelsBroadcast[i] = requestChannelBroadcast
		ResponseChannelsBroadcast[i] = responseChannelBroadcast
	}

	// Create and start 3 Processing Elements
	pes := make([]*processingElement.ProcessingElement, 3) // Create an array of PEs

	for i := 0; i < 3; i++ {
		pe, err := processingElement.New(i, RequestChannelsM1[i], ResponseChannelsM1[i], fmt.Sprintf("generated-programs/program%d.txt", i), terminate)
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

	// Create the Interconnect and attach the communication channels with the 3 CacheControllers
	// Create Interconnect
	interconnect, err := interconnect.New(
						RequestChannelsM2, 
						ResponseChannelsM2, 
						RequestChannelM3, 
						ResponseChannelM3, 
						RequestChannelsBroadcast, 
						ResponseChannelsBroadcast,
						protocol, 
						terminate)
	if err != nil {
		fmt.Printf("Error initializing Interconnect: %v\n", err)
		return
	}

	// Start Interconnect
	wg.Add(1)
	go func() {
		defer wg.Done()
		interconnect.Run(&wg)
	}()

	// Create Main Memory with two channels, ready to connect the interconect
	mainMemory, err := mainMemory.New(RequestChannelM3, ResponseChannelM3, terminate)
	if err != nil {
		fmt.Printf("Error initializing Main Memory: %v\n", err)
		return
	}
	// Start Main Memory
	wg.Add(1)
	go func() {
		defer wg.Done()
		mainMemory.Run(&wg)
	}()

	// THIS IS WHERE THE CLI STARTS *****************************************************************************
	fmt.Println("WELCOME TO MCKEVINHO CLI")
	fmt.Println("The available commands are:")
	fmt.Println("1. step <PE> - Send the Control signal to a specific PE (e.g., 'step 1' or 'step all')")
	fmt.Println("2. about pe  - Print the time stamp of the Procesing Elements")
	fmt.Println("3. about cc  - Print the time stamp of the Cache Controllers")
	fmt.Println("4. about mm  - Print the time stamp of the Main Memory")
	fmt.Println("5. lj        - Terminate the program")

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
						fmt.Printf("Sent 'step' command to PE%d...\n", i)
					} else {
						fmt.Printf("PE%d is not available...\n", pe.ID)
					}
				}
			} else {
				peIndex, err := strconv.Atoi(args[1])
				if err != nil || peIndex < -1 || peIndex > len(pes)-1 {
					fmt.Println("Invalid PE number. Please enter a valid PE number or 'all'.")
					continue
				}

				pe := pes[peIndex]
				if !pe.IsDone && !pe.IsExecutingInstruction {
					pe.Control <- true
					fmt.Printf("Sent 'step' command to PE%d...\n", pe.ID)
				} else {
					fmt.Printf("PE%d is not available...\n", pe.ID)
				}
			}

		case "about":
			if len(args) != 2 {
				fmt.Println("Invalid 'about' command. Please use 'about <PE>' or 'step all'.")
				continue
			}

			if args[1] == "pe" {
				for _, pe := range pes {
					aboutPE, _ := pe.About()
					fmt.Println(aboutPE + "\n\n")
				}
			}

			if args[1] == "cc" {
				for _, cc := range ccs{
					aboutCC, _ := cc.About()
					fmt.Println(aboutCC + "\n\n")
				}
			}

			if args[1] == "mm" {
				aboutMM, _ := mainMemory.About()
				fmt.Println(aboutMM + "\n\n")

			}

			if args[1] == "ic" {
				aboutIC, _ := interconnect.About()
				fmt.Println(aboutIC + "\n\n")

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
				close(RequestChannelsM1[i])
				close(ResponseChannelsM1[i])
				close(RequestChannelsM2[i])
				close(ResponseChannelsM2[i])
				close(RequestChannelsBroadcast[i])
				close(ResponseChannelsBroadcast[i])
			}
			close(RequestChannelM3)
			close(ResponseChannelM3)
			close(semaphore)
			break PELoop
		default:
			fmt.Println("Invalid command. Please enter 'step <PE>' or 'lj'.")
		}
	}
}

func trimNewline(s string) string {
	return s[:len(s)-1]
}

// Function