package mainMemory

import (
	"MultiprocessingSystem/utils"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

// MainMemory representa la memoria principal del sistema.
type MainMemory struct {
	datos [16]uint32 // Array de 16 entradas de 32 bits.
	// mutex           sync.Mutex                // Mutex para garantizar RequestChannel seguro a la memoria.
	RequestChannel  chan utils.RequestMemory  // Canal para solicitudes de interconect a la memoria.
	ResponseChannel chan utils.ResponseMemory // Canal para respond de memoria a interconect.
	Quit            chan struct{}
	Logger          *log.Logger
}

func New(requestChannel chan utils.RequestMemory, responseChannel chan utils.ResponseMemory, quit chan struct{}) (*MainMemory, error) {
	logFile, err := os.Create("logs/MM/MM.log")
	if err != nil {
		log.Fatalf("Error creating log file for Main Memory: %v", err)
		// RequestChannel: make(chan *RequestChannelMemoria),
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
	// go mm.gestionarRequestChannels()

	return &MainMemory{
		datos:           dataInitialized,
		RequestChannel:  requestChannel,
		ResponseChannel: responseChannel,
		Quit:            quit,
		Logger:          logger,
	}, nil
}

// Read accede a la memoria principal para Read un value en una dirección.
func (mm *MainMemory) Read(address int) uint32 {
	// mm.mutex.Lock()
	// defer mm.mutex.Unlock()
	return mm.datos[address]
}

// Write actualiza un value en una dirección en la memoria principal.
func (mm *MainMemory) Write(address int, value uint32) {
	// mm.mutex.Lock()
	// defer mm.mutex.Unlock()
	mm.datos[address] = value
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
			mm.Logger.Printf(" - MM is Received request from Interconect.\n")

			// Initialize a struct to pack the response
			mm.Logger.Printf(" - MM is initializing a response object.\n")
			response := utils.ResponseMemory{
				Status: false,
				Type:   request.Type,
				Value:  randomNum,
				Time:   1,
			}

			switch request.Type {
			case "READ":
				mm.Logger.Printf(" - MM is processing a READ request.\n")
				mm.Logger.Printf(" - Address: %d.\n", request.Address)
				time.Sleep(time.Duration(READTIMECOST) * time.Second)
				response.Value = mm.Read(request.Address)
				response.Time = READTIMECOST
				response.Status = true

			case "WRITE":
				mm.Logger.Printf(" - MM is processing a WRITE request.\n")
				mm.Logger.Printf(" - Address: %d, Data: %d.\n", request.Address, request.Value)
				time.Sleep(time.Duration(WRITETIMECOST) * time.Second)
				response.Value = request.Value
				response.Time = WRITETIMECOST
				response.Status = true

			}

			// Enviar respuesta al interconnect
			mm.Logger.Printf(" - MM is about to send the response to IC.\n")
			mm.ResponseChannel <- response
			mm.Logger.Printf(" - MM has sent the response back to IC.\n")

		case <-mm.Quit:
			mm.Logger.Printf(" - MM has received an external signal to terminate.\n")
			return
		}

	}
}
