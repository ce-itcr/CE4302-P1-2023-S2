package mainMemory

import (
	"Backend/utils"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
	"encoding/json"
)

// MainMemory representa la memoria principal del sistema.
type MainMemory struct {
	Data [16]uint32 // Array de 16 entradas de 32 bits.
	RequestChannel  chan utils.RequestMainMemory  // Canal para solicitudes de interconect a la memoria.
	ResponseChannel chan utils.ResponseMainMemory // Canal para respond de memoria a interconect.
	Quit            chan struct{}
	Status 			string
	Logger          *log.Logger
}

func New(
		requestChannel chan utils.RequestMainMemory,
		responseChannel chan utils.ResponseMainMemory,
		logfilepath string,
		quit chan struct{}) (*MainMemory, error) {

	logFile, err := os.Create(logfilepath + "MM.log")
	if err != nil {
		log.Fatalf("Error creating log file for Main Memory: %v", err)
	}

	// Create seed for random numbers
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)

	// Inicializa la memoria con valuees predeterminados si es necesario.
	dataInitialized := [16]uint32{}
	for i := 0; i < len(dataInitialized); i++ {
		dataInitialized[i] = uint32(generator.Intn(51)) // Genera números entre 0 y 50.
	}

	logger := log.New(logFile, "MM"+"_", log.Ldate|log.Ltime)

	// Inicia el goroutine para gestionar las solicitudes de RequestChannel a la memoria.
	return &MainMemory{
		Data:           dataInitialized,
		RequestChannel:  requestChannel,
		ResponseChannel: responseChannel,
		Quit:            quit,
		Logger:          logger,
		Status: "Active",
	}, nil
}

func (mm *MainMemory) About()(string, error){
	// Create an empty BlockObjectList
	memoryBlocks := utils.BlockObjectList{}
    for i := 0; i <= 15; i++ {
        // Create a new BlockObject instance
		blockObj := utils.BlockObject{
			Address: i,
			Data: int(mm.Data[i]),
		}
		// Append the new BlockObject to the BlockObjectList
		memoryBlocks = append(memoryBlocks, blockObj)
    }

    // Create a the final JSON struct
    aboutMM := utils.AboutMainMemory{
		Status: mm.Status,
		Blocks: memoryBlocks,
	}

	// Marshal the PE struct into a JSON string
	jsonData, err := json.MarshalIndent(aboutMM, "", "    ")
	if err != nil {
		return "", err
	}

	// Convert the byte slice to a string
	jsonString := string(jsonData)

	return jsonString, nil
}

// Read accede a la memoria principal para Read un value en una dirección.
func (mm *MainMemory) Read(address int) uint32 {
	return mm.Data[address]
}

// Write actualiza un value en una dirección en la memoria principal.
func (mm *MainMemory) Write(address int, value uint32) {
	mm.Data[address] = value
}

func (mm *MainMemory) Run(wg *sync.WaitGroup) {
	// Define time cost per write and read operations
	WRITETIMECOST := 5
	READTIMECOST := 3
	randomNum := uint32(12)

	mm.Logger.Printf(" - MM is running.\n")
	for {
		// Listen to the interconect for a request
		select {
		case request := <-mm.RequestChannel:
			mm.Logger.Printf(" - MM Received request from Interconect.\n")

			// Initialize a struct to pack the response
			response := utils.ResponseMainMemory{
				Status: false,
				Type:   request.Type,
				Value:  randomNum,
				Time:   1,
			}

			switch request.Type {
			case "READ":
				mm.Logger.Printf(" - MM is processing a READ request.\n")
				mm.Logger.Printf(" - Address: %d.\n", request.Address)
				time.Sleep(3 * time.Second)
				response.Value = mm.Read(request.Address)
				response.Time = READTIMECOST
				response.Status = true

			case "WRITE":
				mm.Logger.Printf(" - MM is processing a WRITE request.\n")
				mm.Logger.Printf(" - Address: %d, Data: %d.\n", request.Address, request.Value)
				time.Sleep(5 * time.Second)
				mm.Write(request.Address, request.Value)
				response.Value = request.Value
				response.Time = WRITETIMECOST
				response.Status = true

			}

			// Enviar respuesta al interconnect
			mm.ResponseChannel <- response
			mm.Logger.Printf(" - MM has sent the response back to IC.\n")

		case <-mm.Quit:
			mm.Logger.Printf(" - MM has received an external signal to terminate.\n")
			return
		}

	}
}