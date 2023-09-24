package processingElement

import (
	"strconv"
    "fmt"
    "sync"
    "bufio"
    "os"
    "strings"
    //"time"
    "MultiprocessingSystem/utils"
)

// Represents a processing element.
type ProcessingElement struct {
    ID          int                 // Identifier of the PE
    Name        string              // Name for the PE
    Instructions []string           // Array of instructions loaded
    Control     chan bool           // Channel for external control
    Done        chan bool           // Channel to signal completion
    RequestChannel chan utils.Request     // Channel to send a request to a CacheController
    ResponseChannel chan utils.Response   // Channel to wait for a response from a CacheController
    register int                    // The one and only register
}

// New creates a new ProcessingElement instance.
func New(id int, name string, RequestChannelCC chan utils.Request, ResponseChannelCC chan utils.Response ,filename string) (*ProcessingElement, error) {

    // Load the program from the text file
    instructions, err := readInstructionsFromFile(filename)
    if err != nil {
        return nil, err
    }

    return &ProcessingElement{
        ID:           id,
        Name:         name,
        Instructions: instructions,
        RequestChannel: RequestChannelCC,
        ResponseChannel: ResponseChannelCC,
        Control:      make(chan bool),
        Done:         make(chan bool),
        register: 0,
    }, nil
}

// Run simulates the execution of instructions for a ProcessingElement.
func (pe *ProcessingElement) Run(wg *sync.WaitGroup) {
    defer wg.Done()

    fmt.Printf("PE %d (%s) is ready to execute instructions:\n", pe.ID, pe.Name)
    for i, instruction := range pe.Instructions {
        select {
        case <- pe.Control:
            fmt.Printf(" - PE %d (%s) received external signal to execute instruction %d: %s\n", pe.ID, pe.Name, i+1, instruction)
            // Execute the instruction **********************************************************************************************
            words := strings.Fields(instruction)
            operation := words[0]

            switch operation {
            case "INC":
                fmt.Printf(" - PE %d (%s) is executing a %s operation\n", pe.ID, pe.Name, operation)
                pe.register++

            case "READ":
                fmt.Printf(" - PE %d (%s) is executing a %s operation\n", pe.ID, pe.Name, operation)
                // Create a request structure
                address, err := strconv.Atoi(words[1])
                if err != nil {
                    // Error parsing the integer
                    return
                }
                request := utils.Request{
                    Type: operation,
                    Address: address,                    
                    Data: 0,                      
                }

                // Send the request to the CacheController
                pe.RequestChannel <- request

                // Wait for the response from the CacheController
                response := <- pe.ResponseChannel

                // Process the response
                fmt.Printf(" - PE %d (%s) received the response --> Status: %v, Type: %s, Data: %d\n", pe.ID, pe.Name, response.Status, response.Type, response.Data)


            case "WRITE":
                fmt.Printf(" - PE %d (%s) is executing a %s operation\n", pe.ID, pe.Name, operation)
                // Create a request structure
                address, err := strconv.Atoi(words[1])
                if err != nil {
                    // Error parsing the integer
                    return
                }
                request := utils.Request{
                    Type: operation,
                    Address: address,                   
                    Data: pe.register,                     
                }

                // Send the request to the CacheController
                pe.RequestChannel <- request

                // Wait for the response from the CacheController
                response := <- pe.ResponseChannel

                // Process the response
                fmt.Printf(" - PE %d (%s) received the response --> Status: %v, Type: %s, Data: %d\n", pe.ID, pe.Name, response.Status, response.Type, response.Data)


            }
            fmt.Printf(" - PE %d (%s) has finished with the instruction.\n", pe.ID, pe.Name)
            
        case <-pe.Done:
            fmt.Printf(" - PE %d (%s) has completed execution of instructions.\n", pe.ID, pe.Name)
            return
        }
    }
    fmt.Printf(" - PE %d (%s) has executed all instructions.\n", pe.ID, pe.Name)
    return
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


