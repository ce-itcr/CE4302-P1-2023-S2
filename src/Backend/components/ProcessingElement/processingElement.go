package processingElement

import (
	"strconv"
    "sync"
    "bufio"
    "os"
    "strings"
    "log"
    "encoding/json"

    "Backend/utils"
)


// Represents a Processing Element
type ProcessingElement struct {
    ID          int                                         // Identifier of the PE
    Instructions utils.QueueS                                // Queue of instructions loaded
    Logger          *log.Logger                             // Local logger
    Control     chan bool                                   // Channel for external control
    RequestChannel chan utils.RequestProcessingElement      // Channel to send a request to a CacheController
    ResponseChannel chan utils.ResponseProcessingElement    // Channel to wait for a response from a CacheController
    register int                                            // The one and only register
    IsDone bool                                             // Flag to know when a PE hasn't finished executing instructions
    IsExecutingInstruction bool                             // Flag to know when a PE is currently executing an instruction
    Quit chan struct{}                                      // A signal to terminate the goroutine
    Status string                                           // Status for every momment of the execution
    Filename string                                         // The name of the text file where the instructions are
}

// New creates a new ProcessingElement instance with the required information to operate.
func New(
        id int, 
        RequestChannelCC chan utils.RequestProcessingElement , 
        ResponseChannelCC chan utils.ResponseProcessingElement,
        filename string, 
        quit chan struct{}) (*ProcessingElement, error) {

    // Load the program from the text file
    instructions, err := readInstructionsFromFile(filename)
    if err != nil {
        return nil, err
    }

    // Create the log file for this object
    logFile, err := os.Create("logs/PE/PE" + strconv.Itoa(id) + ".log")
    if err != nil {
        log.Fatalf("Error creating log file for PE%d: %v", id, err)
    }

    // Initialize logger for the PE using its respective log file
    logger1 := log.New(logFile, "PE" + strconv.Itoa(id) + "_", log.Ldate|log.Ltime)

    // Return a pointer with the Processing Element instance
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
        Status: "Active",
        Filename: filename,
    }, nil
}

// Function to get a JSON string with the current state of the Processing Element
func (pe *ProcessingElement) About()(string, error){
    // Create an empty InstructionObjectList
    instructions := utils.InstructionObjectList{}

    // Get the values from the instructions queue of the Processing Element
    for i, item := range pe.Instructions.Items{
		// Create an InstructionObject
        instructionObj := utils.InstructionObject{
            Position: i,
            Instruction: item,
        }
        // Append the instruction object to the list
        instructions = append(instructions, instructionObj)
	}

    // Create a struct
    aboutPE := utils.AboutProcessingElement {
        ID: pe.ID,
        Register: pe.register,
        Status: pe.Status,
        Instructions: instructions,
    }

	// Marshal the PE struct into a JSON string
	jsonData, err := json.MarshalIndent((aboutPE), "", "    ")
	if err != nil {
		return "", err
	}

	// Convert the byte slice to a string
	jsonString := string(jsonData)

	return jsonString, nil
}

// Function to send a request to the Cache Controller
func (pe *ProcessingElement) RequestCacheController(Type string, Address int, Data int) (int, bool) {
    // Prepare the request message for the Cache Controller
    request := utils.RequestProcessingElement {
        Type: Type,
        Address: Address,                    
        Data: Data,                      
    }
    // Send the request to the CacheController
    pe.Status = "Sending request to CC"
    pe.RequestChannel <- request
    pe.Status = "Sent request to CC"
    pe.Logger.Printf(" - PE%d sent (Type: %s, Address: %d) to the Cache Controller.\n", pe.ID, Type, Address)

    // Wait for the response from the CacheController
    pe.Logger.Printf(" - PE%d is waiting for a response from the Cache Controller....\n", pe.ID)
    pe.Status = "Waiting for a response from CC"
    response := <- pe.ResponseChannel
    pe.Status = "Received response from CC"
    DataResponse := response.Data
    StatusResponse := response.Status

    // Process the response
    pe.Logger.Printf(" - PE%d received a response from the Cache Controller.\n", pe.ID)

    // Return the response values
    return DataResponse, StatusResponse
}

// Run simulates the execution of instructions for a ProcessingElement.
func (pe *ProcessingElement) Run(wg *sync.WaitGroup) {
    pe.Logger.Printf(" - PE%d is ready to execute instructions.\n", pe.ID)
    pe.Status = "Ready"
    for {
        select {
            // The PE receives a signal to execute an instruction
            case <- pe.Control:
                // Let others know the PE is currently busy executing an instruction
                pe.IsExecutingInstruction = true
                pe.Status = "Signal Received"

                // Check if there are still instructions to execute
                if pe.Instructions.IsEmpty() {
                    pe.Logger.Printf(" - PE%d has executed all instructions.\n", pe.ID)
                    // Notify the main that this PE has executed all instructions
                    pe.IsDone = true
                    pe.Status = "Done"
                    return
                }
                
                // Get the next instruction
                instruction := pe.Instructions.Dequeue()

                pe.Logger.Printf(" - PE%d received external signal to execute instruction: %s.\n", pe.ID, instruction)
                words := strings.Fields(instruction)
                operation := words[0]

                switch operation {
                // Increment
                case "INC":
                    pe.Status = "Executing INC"
                    pe.Logger.Printf(" - PE%d is executing a %s operation.\n", pe.ID, operation)
                    pe.register++
                    pe.Logger.Printf(" - Now the value of the register is %d.\n", pe.register)
                    pe.Logger.Printf(" - PE%d has finished with the instruction.\n", pe.ID)
                    // Let others know the PE is now available
                    pe.IsExecutingInstruction = false

                // Read a data from an specific memory address
                case "READ":
                    pe.Status = "Executing READ"
                    pe.Logger.Printf(" - PE%d is executing a %s operation.\n", pe.ID, operation)
                    // Create a request structure
                    address, err := strconv.Atoi(words[1])
                    if err != nil {
                        // Error parsing the integer
                        return
                    }
                    
                    // Send a READ request to the Cache Controller
                    Data, _ := pe.RequestCacheController(operation, address, 0)

                    // Process the response values
                    pe.Logger.Printf(" - PE%d received Data: %d.\n", pe.ID, Data)
                    pe.Status = "Updating Register"
                    pe.register = Data
                    pe.Logger.Printf(" - Updated local register: Rs = %d.\n", pe.register)
                    pe.Logger.Printf(" - PE%d has finished with the instruction.\n", pe.ID)

                    // Let others know the PE is now available
                    pe.IsExecutingInstruction = false
                    pe.Status = "Free"

                // Write dato into an specific memory address
                case "WRITE":
                    pe.Status = "Executing WRITE"
                    pe.Logger.Printf(" - PE%d is executing a %s operation.\n", pe.ID, operation)
                    // Create a request structure
                    address, err := strconv.Atoi(words[1])
                    if err != nil {
                        // Error parsing the integer
                        return
                    }

                    // Send a READ request to the Cache Controller
                    _, Status := pe.RequestCacheController(operation, address, pe.register)

                    // Process the response values
                    pe.Logger.Printf(" - PE%d received --> Status: %v.\n", pe.ID, Status)
                    pe.Logger.Printf(" - PE%d has finished with the instruction.\n", pe.ID)

                    // Let others know the PE is now available
                    pe.IsExecutingInstruction = false
                    pe.Status = "Free"

                }

            // When the PE receives a signal terminate
            case <- pe.Quit:
                pe.Logger.Printf(" - PE%d received termination signal and is exiting gracefully.\n", pe.ID)
                pe.Status = "Forced to quit"
                return
        }

    }
}

// // Run simulates the execution of instructions for a ProcessingElement.
// func (pe *ProcessingElement) ManualRun(wg *sync.WaitGroup) {
//     pe.Logger.Printf(" - PE%d is ready to execute instructions.\n", pe.ID)
//     pe.Status = "Ready"
//     for {
//         select {
//             // When the PE receives a signal to execute an instruction
//             case <- pe.Control:
//                 // Let others know the PE is currently busy executing an instruction
//                 pe.IsExecutingInstruction = true
//                 pe.Status = "Signal Received"

//                 // Extract the instruction
//                 instruction := pe.readInstructionFromFile(pe.Filename)

//                 pe.Logger.Printf(" - PE%d received external signal to execute instruction: %s.\n", pe.ID, instruction)
//                 words := strings.Fields(instruction)
//                 operation := words[0]

//                 switch operation {
//                 // Increment
//                 case "INC":
//                     pe.Status = "Executing INC"
//                     pe.Logger.Printf(" - PE%d is executing a %s operation.\n", pe.ID, operation)
//                     pe.register++
//                     pe.Logger.Printf(" - Now the value of the register is %d.\n", pe.register)
//                     pe.Logger.Printf(" - PE%d has finished with the instruction.\n", pe.ID)
//                     // Let others know the PE is now available
//                     pe.IsExecutingInstruction = false

//                 // Read a data from an specific memory address
//                 case "READ":
//                     pe.Status = "Executing READ"
//                     pe.Logger.Printf(" - PE%d is executing a %s operation.\n", pe.ID, operation)
//                     // Create a request structure
//                     address, err := strconv.Atoi(words[1])
//                     if err != nil {
//                         // Error parsing the integer
//                         return
//                     }
                    
//                     // Send a READ request to the Cache Controller
//                     Data, _ := pe.RequestCacheController(operation, address, 0)

//                     // Process the response values
//                     pe.Logger.Printf(" - PE%d received Data: %d.\n", pe.ID, Data)
//                     pe.Status = "Updating Register"
//                     pe.register = Data
//                     pe.Logger.Printf(" - Updated local register: Rs = %d.\n", pe.register)
//                     pe.Logger.Printf(" - PE%d has finished with the instruction.\n", pe.ID)

//                     // Let others know the PE is now available
//                     pe.IsExecutingInstruction = false
//                     pe.Status = "Free"

//                 // Write dato into an specific memory address
//                 case "WRITE":
//                     pe.Status = "Executing WRITE"
//                     pe.Logger.Printf(" - PE%d is executing a %s operation.\n", pe.ID, operation)
//                     // Create a request structure
//                     address, err := strconv.Atoi(words[1])
//                     if err != nil {
//                         // Error parsing the integer
//                         return
//                     }

//                     // Send a READ request to the Cache Controller
//                     _, Status := pe.RequestCacheController(operation, address, pe.register)

//                     // Process the response values
//                     pe.Logger.Printf(" - PE%d received --> Status: %v.\n", pe.ID, Status)
//                     pe.Logger.Printf(" - PE%d has finished with the instruction.\n", pe.ID)

//                     // Let others know the PE is now available
//                     pe.IsExecutingInstruction = false
//                     pe.Status = "Free"

//                 }

//             // When the PE receives a signal terminate
//             case <- pe.Quit:
//                 pe.Logger.Printf(" - PE%d received termination signal and is exiting gracefully.\n", pe.ID)
//                 pe.Status = "Forced to quit"
//                 return
//         }
//     }
// }

// Reads lines from a text file and returns them as a slice of strings.
func readInstructionsFromFile(filename string) (utils.QueueS, error) {
    // Create a new queue to store the instructions from the text file
	Instructions := utils.QueueS{}

    // Try opening the text file
    file, err := os.Open(filename)
    if err != nil {
        return Instructions, err
    }
    // Close the file after reading all lines
    defer file.Close()

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        // Check if the instruction is valid
        if (isValidInstruction(line)){
            Instructions.Enqueue(line)
        }
    }

    // If there was an error, return an empty queue
    if err := scanner.Err(); err != nil {
        return Instructions, err
    }

    // If not, return the queue with the instructions
    return Instructions, nil
}

// // Function that reads the first line of the program and returns it.
// func (pe *ProcessingElement) readInstructionFromFile(filename string) string {

//     file, err := os.Open(filename)
//     if err != nil {
//         return ""
//     }

//     defer file.Close()

//     scanner := bufio.NewScanner(file)
//     scanner.Scan()
//     line := strings.TrimSpace(scanner.Text())
//     if (isValidInstruction(line)){
//         return line
//     } else  {
//         return ""
//     }

// }


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
