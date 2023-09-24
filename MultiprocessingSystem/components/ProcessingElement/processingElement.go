package processingElement

import (
	"strconv"
    "fmt"
    "sync"
    "bufio"
    "os"
    "strings"
    "time"
)

// Represents a processing element.
type ProcessingElement struct {
    ID          int
    Name        string
    Instructions []string
    Control     chan bool // Channel for external control
    Done        chan bool // Channel to signal completion
}

// New creates a new ProcessingElement instance.
func New(id int, name string, filename string) (*ProcessingElement, error) {

    // Load the program from the text file
    instructions, err := readInstructionsFromFile(filename)
    if err != nil {
        return nil, err
    }

    return &ProcessingElement{
        ID:           id,
        Name:         name,
        Instructions: instructions,
        Control:      make(chan bool),
        Done:         make(chan bool),
    }, nil
}

// Run simulates the execution of instructions for a ProcessingElement.
func (pe *ProcessingElement) Run(wg *sync.WaitGroup) {
    defer wg.Done()

    fmt.Printf("PE %d (%s) is ready to execute instructions:\n", pe.ID, pe.Name)
    for i, instruction := range pe.Instructions {
        select {
        case <-pe.Control:
            fmt.Printf(" - PE %d (%s) received external signal to execute instruction %d: %s\n", pe.ID, pe.Name, i+1, instruction)
            // Simulate execution time (you can replace this with the actual work)
            time.Sleep(1 * time.Second)
        case <-pe.Done:
            fmt.Printf(" - PE %d (%s) has completed execution of instructions.\n", pe.ID, pe.Name)
            return
        }
    }
    fmt.Printf(" - PE %d (%s) has executed all instructions.\n", pe.ID, pe.Name)
}

// Reads lines from a text file and returns them as a slice of strings.
func readInstructionsFromFile(filename string) ([]string, error) {
    var instructions []string

    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        instructions = append(instructions, line)
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return instructions, nil
}

// isValidInstruction checks if an instruction is valid.
func isValidInstruction(instruction string) bool {
    // Split the instruction into words
    words := strings.Fields(instruction)

    if len(words) == 1 {
        // Single-word instructions (e.g., "INC")
        return words[0] == "INC"
    } else if len(words) == 2 {
        // Two-word instructions (e.g., "READ 5" or "WRITE 10")
        if words[0] == "READ" || words[0] == "WRITE" {
            num, err := strconv.Atoi(words[1])
            if err != nil {
                // Error parsing the integer
                return false
            }
            return num >= 0 && num <= 15
        }
    }

    // Invalid instruction format
    return false
}



