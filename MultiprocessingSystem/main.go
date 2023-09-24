package main

import (
    "fmt"
    "os"
    "os/signal"
    "strconv"
    "sync"
    "syscall"
    "MultiprocessingSystem/components/ProcessingElement"
)

func main() {
    // Create three ProcessingElement instances with different instructions
    pe1, err1 := processingElement.New(1, "PE1", "programs/program1.txt")
    pe2, err2 := processingElement.New(2, "PE2", "programs/program2.txt")
    pe3, err3 := processingElement.New(3, "PE3", "programs/program3.txt")

    if err1 != nil || err2 != nil || err3 != nil {
        fmt.Println("Error:", err1, err2, err3)
        return
    }

    var wg sync.WaitGroup

    // Start each ProcessingElement as a goroutine (Thread)
    wg.Add(3)
    go pe1.Run(&wg)
    go pe2.Run(&wg)
    go pe3.Run(&wg)

    // Set up a signal handler for external control
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

    // Channel to receive special termination word
    specialTerminationCh := make(chan struct{})

    // Map to keep track of active PEs
    activePEs := map[int]*processingElement.ProcessingElement{
        pe1.ID: pe1,
        pe2.ID: pe2,
        pe3.ID: pe3,
    }

    fmt.Println("Press Ctrl+C to control the Processing Elements.")
    fmt.Println("Type 'terminate' to stop all Processing Elements.")

    // Listen for external control signals
    go func() {
        for sig := range sigCh {
            fmt.Printf("Received signal: %v\n", sig)
            handleSignal(sig, activePEs)
        }
    }()

    // Listen for the special termination word
    go func() {
        var specialWord string
        fmt.Scanln(&specialWord)
        if specialWord == "terminate" {
            fmt.Println("Terminating all Processing Elements...")
            close(specialTerminationCh)
        }
    }()

    // Wait for all processing elements to finish or special termination word
    select {
    case <-specialTerminationCh:
        // Termination word received
        // Signal all active PEs to terminate
        for _, pe := range activePEs {
            pe.Done <- true
        }
    case <-pe1.Done:
        // PE1 completed, remove it from active PEs
        delete(activePEs, pe1.ID)
    case <-pe2.Done:
        // PE2 completed, remove it from active PEs
        delete(activePEs, pe2.ID)
    case <-pe3.Done:
        // PE3 completed, remove it from active PEs
        delete(activePEs, pe3.ID)
    }

    // Wait for all remaining PEs to finish
    wg.Wait()

    fmt.Println("All Processing Elements have finished.")
}

// handleSignal handles external control signals to start execution of instructions.
func handleSignal(sig os.Signal, activePEs map[int]*processingElement.ProcessingElement) {
    if sig == syscall.SIGINT {
        var peIDs []string
        for id := range activePEs {
            peIDs = append(peIDs, strconv.Itoa(id))
        }
        fmt.Println("Available Processing Elements:", peIDs)
        fmt.Print("Enter the ID of the PE to execute an instruction (or 'all'): ")
        var input string
        fmt.Scanln(&input)

        if input == "all" {
            for _, pe := range activePEs {
                pe.Control <- true
            }
        } else {
            peID, err := strconv.Atoi(input)
            if err != nil {
                fmt.Println("Invalid input. Please enter a valid PE ID or 'all'.")
                return
            }

            if pe, ok := activePEs[peID]; ok {
                pe.Control <- true
            } else {
                fmt.Println("PE with ID", peID, "not found.")
            }
        }
    } else if sig == syscall.SIGTERM {
        for _, pe := range activePEs {
            pe.Done <- true
        }
    }
}
