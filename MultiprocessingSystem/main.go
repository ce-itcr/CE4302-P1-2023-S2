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

    // Create a termination channel to signal termination
    terminate := make(chan struct{})
    

    // Instantiate CacheController
    cacheController, err := CacheController.New(requestChannelCC, responseChannelCC, terminate)
    if err != nil {
        fmt.Printf("Error initializing CacheController: %v\n", err)
        return
    }

    // Instantiate ProcessingElement
    processingElement, err := processingElement.New(1, "PE1", requestChannelCC, responseChannelCC, "programs/program1.txt", terminate)
    if err != nil {
        fmt.Printf("Error initializing ProcessingElement: %v\n", err)
        return
    }

    // Use separate WaitGroups for CacheController and ProcessingElement
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
PELoop:
    for {
        fmt.Print("Enter a command: ")
        command, _ := reader.ReadString('\n')
        command = trimNewline(command)

        switch command {
        case "step":
            if !processingElement.IsDone  && !processingElement.IsExecutingInstruction{
                // Sending a control signal to ProcessingElement
                processingElement.Control <- true
                fmt.Println("Sent 'step' command to ProcessingElement")
            } else {
                fmt.Println("PE is not available")
            }
        case "lj":
            // Signal termination to both components
            fmt.Println("Sent 'lj' command to terminate the program")
            close(terminate)


            wg.Wait() // Wait for both goroutines to finish gracefully
        
            close(requestChannelCC)
            close(responseChannelCC)
            break PELoop
        default:
            fmt.Println("Invalid command. Please enter 'step' or 'lj'.")
        }
    }
}


func trimNewline(s string) string {
    return s[:len(s)-1]
}
