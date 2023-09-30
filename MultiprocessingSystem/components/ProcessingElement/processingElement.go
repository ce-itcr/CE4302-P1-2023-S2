package processingElement

import (
	"strconv"
    "sync"
    "bufio"
    "os"
    "strings"
    "log"

    "MultiprocessingSystem/utils"
)

// Represents a processing element.
type ProcessingElement struct {
    ID          int                         // Identifier of the PE
    Instructions []string                   // Array of instructions loaded
    Logger          *log.Logger              //
    Control     chan bool                   // Channel for external control
    RequestChannel chan utils.RequestM1        // Channel to send a request to a CacheController
    ResponseChannel chan utils.ResponseM1      // Channel to wait for a response from a CacheController
    register int                            // The one and only register
    IsDone bool                             // Flag to know when a PE hasn't finished executing instructions
    IsExecutingInstruction bool             // Flag to know when a PE is currently executing an instruction
    Quit chan struct{}                      // A signal to terminate the goroutine
}

// New creates a new ProcessingElement instance.
func New(id int, RequestChannelCC chan utils.RequestM1 , ResponseChannelCC chan utils.ResponseM1  ,filename string, quit chan struct{}) (*ProcessingElement, error) {
    // Load the program from the text file
    instructions, err := readInstructionsFromFile(filename)
    if err != nil {
        return nil, err
    }

    // Create the log file
    logFile, err := os.Create("logs/PE/PE" + strconv.Itoa(id) + ".log")
    if err != nil {
        log.Fatalf("Error creating log file for PE%d: %v", id, err)
    }

    // Initialize logger for the PE using its respective log file
    logger1 := log.New(logFile, "PE" + strconv.Itoa(id) + "_", log.Ldate|log.Ltime)
	

    return &ProcessingElement{
        ID:           id,
        Instructions: instructions,
        Logger: logger1,
        RequestChannel: RequestChannelCC,
        ResponseChannel: ResponseChannelCC,
        Control:      make(chan bool),
        register: 0,
        IsDone : false,
        IsExecutingInstruction: false,
        Quit: quit,
    }, nil
}

// Run simulates the execution of instructions for a ProcessingElement.
func (pe *ProcessingElement) Run(wg *sync.WaitGroup) {
    pe.Logger.Printf(" - PE%d is ready to execute instructions.\n", pe.ID)
    for i, instruction := range pe.Instructions {
        select {
            // When the PE receives a signal to execute an instruction
            case <- pe.Control:
                // Let others know the PE is currently busy executing an instruction
                pe.IsExecutingInstruction = true

                pe.Logger.Printf(" - PE%d received external signal to execute instruction %d: %s.\n", pe.ID, i+1, instruction)
                words := strings.Fields(instruction)
                operation := words[0]

                switch operation {
                // Increment
                case "INC":
                    pe.Logger.Printf(" - PE%d is executing a %s operation.\n", pe.ID, operation)
                    pe.register++
                    pe.Logger.Printf(" - Now the value of the register is %d.\n", pe.register)
                    pe.Logger.Printf(" - PE%d has finished with the instruction.\n", pe.ID)
                    // Let others know the PE is now available
                    pe.IsExecutingInstruction = false

                // Read a data from an specific memory address
                case "READ":
                    pe.Logger.Printf(" - PE%d is executing a %s operation.\n", pe.ID, operation)
                    // Create a request structure
                    address, err := strconv.Atoi(words[1])
                    if err != nil {
                        // Error parsing the integer
                        return
                    }
                    // Prepare the request message
                    request := utils.RequestM1 {
                        Type: operation,
                        Address: address,                    
                        Data: 0,                      
                    }

                    // Send the request to the CacheController
                    pe.RequestChannel <- request
                    pe.Logger.Printf(" - PE%d sent (Type: %s, Address: %d) to the Cache Controller.\n", pe.ID, operation, address)

                    // Wait for the response from the CacheController
                    pe.Logger.Printf(" - PE%d is waiting for a response from the Cache Controller....\n", pe.ID)
                    response := <- pe.ResponseChannel

                    // Process the response
                    pe.Logger.Printf(" - PE%d received --> Status: %v, Type: %s, Data: %d.\n", pe.ID,response.Status, response.Type, response.Data)
                    pe.register = response.Data
                    pe.Logger.Printf(" - Now the value of the register is %d.\n", pe.register)
                    pe.Logger.Printf(" - PE%d has finished with the instruction.\n", pe.ID)

                    // Let others know the PE is now available
                    pe.IsExecutingInstruction = false

                // Write dato into an specific memory address
                case "WRITE":
                    pe.Logger.Printf(" - PE%d is executing a %s operation.\n", pe.ID, operation)
                    // Create a request structure
                    address, err := strconv.Atoi(words[1])
                    if err != nil {
                        // Error parsing the integer
                        return
                    }

                    // Prepare the request message
                    request := utils.RequestM1 {
                        Type: operation,
                        Address: address,                   
                        Data: pe.register,                     
                    }

                    // Send the request to the CacheController
                    pe.RequestChannel <- request
                    pe.Logger.Printf(" - PE%d sent (Type: %s, Address: %d, Data: %d) to the Cache Controller.\n", pe.ID, operation, address, pe.register)


                    // Wait for the response from the CacheController
                    pe.Logger.Printf(" - PE%d is waiting for a response from the Cache Controller....\n", pe.ID)
                    response := <- pe.ResponseChannel

                    // Process the response
                    pe.Logger.Printf(" - PE%d received --> Status: %v, Type: %s.\n", pe.ID, response.Status, response.Type)
                    pe.Logger.Printf(" - PE%d has finished with the instruction.\n", pe.ID)

                    // Let others know the PE is now available
                    pe.IsExecutingInstruction = false

                }

            // When the PE receives a signal terminate
            case <- pe.Quit:
                pe.Logger.Printf(" - PE%d received termination signal and is exiting gracefully.\n", pe.ID)
                return
        }

    }
    pe.Logger.Printf(" - PE%d has executed all instructions.\n", pe.ID)
    // Notify the main that this PE has executed all instructions
    pe.IsDone = true
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



