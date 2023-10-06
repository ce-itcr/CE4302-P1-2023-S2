package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"Backend/components/MultiprocessingSystem"
)

func main() {
	// Create a new Multiprocessing system
	mps := MultiprocessingSystem.Start("MESI", true, 4)

	// THIS IS WHERE THE CLI STARTS *****************************************************************************
	fmt.Println("WELCOME TO MCKEVINHO CLI")
	fmt.Println("The available commands are:")
	fmt.Println("1. step <PE> 	- Send the Control signal to a specific PE (e.g., 'step 1' or 'step all')")
	fmt.Println("2. about     	- Print the time stamp of the Procesing Elements")
	fmt.Println("3. results   	- Print the statistical information")
	fmt.Println("4. restart   	- Start a new Multiprocessing System")
	fmt.Println("4. finished? 	- Start a new Multiprocessing System")
	fmt.Println("5. lj        	- Check iuf the Multiprocessing System has finished")

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
				// Start executing all instructions for all Processing Elements
				mps.StartProcessingElements()

			} else {
				peIndex, err := strconv.Atoi(args[1])
				if err != nil || peIndex < -1 || peIndex > len(mps.ProcessingElements)-1 {
					fmt.Println("Invalid PE number. Please enter a valid PE number or 'all'.")
					continue
				}
				// Stepping for a Processing Element
				mps.SteppingProcessingElement(peIndex)
			}

		case "about":
			if len(args) == 1 {
				aboutMps, _ := mps.GetState()
				fmt.Println(aboutMps)
			}

		case "results":
			if len(args) == 1 {
				aboutMps, _ := mps.AboutResults()
				fmt.Println(aboutMps)
			}

		case "restart":
			if len(args) == 1 {
				// Start a new Multiprocessing System
				mps.Stop()
				mps = MultiprocessingSystem.Start("MESI", true, 4)
			}
		case "finished?":
			if len(args) == 1 {
				// Check if the Multiprocessing System is done executing
				if mps.AreWeFinished() {
					fmt.Println("The Multiprocessing System has already finished.")
				} else {
					fmt.Println("The Multiprocessing System has not finished yet.")
				}
			}

		case "lj":
			// Signal termination to both components
			fmt.Println("Sent 'lj' command to terminate the program")

			// Stop/Delete the Multiprocessing system ****************************
			mps.Stop()

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
