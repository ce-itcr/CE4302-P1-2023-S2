package restfulapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"Backend/components/CacheController"
	interconnect "Backend/components/Interconnect"
	mainMemory "Backend/components/MainMemory"
	processingElement "Backend/components/ProcessingElement"
	"Backend/utils"
)

// Estructura para los datos que quieres compartir con el frontend.
type Datos struct {
	Valor string `json:"valor"`
}

var (
	mutex    sync.Mutex
	data     Datos
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	instructions     [][]utils.Instruction
	pes              []*processingElement.ProcessingElement
	interconnectInit *interconnect.Interconnect
	mainMemoryInit   *mainMemory.MainMemory
	ccs              []*CacheController.CacheController

	terminate           chan struct{}
	terminateRESTfulAPI chan struct{}

	wg sync.WaitGroup

	RequestChannelsM1         []chan utils.RequestProcessingElement
	ResponseChannelsM1        []chan utils.ResponseProcessingElement
	RequestChannelsM2         []chan utils.RequestInterconnect
	ResponseChannelsM2        []chan utils.ResponseInterconnect
	RequestChannelsBroadcast  []chan utils.RequestBroadcast
	ResponseChannelsBroadcast []chan utils.ResponseBroadcast

	RequestChannelM3  chan utils.RequestMainMemory
	ResponseChannelM3 chan utils.ResponseMainMemory

	semaphore chan struct{}
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func Restfulapi() {
	// Configura el enrutador.
	router := mux.NewRouter().StrictSlash(true)

	// Ruta para obtener información sobre PE.
	router.HandleFunc("/about", GetAbouts).Methods("GET")

	// Ruta para establecer datos.
	router.HandleFunc("/set", SetData).Methods("POST")

	// Inicia el servidor web.
	go StartWebSocketServer()

	// Servidor HTTP.
	router.HandleFunc("/", homeLink)

	terminateRESTfulAPI = make(chan struct{})

	// Inicia el servidor HTTP en una goroutine.
	go func() {
		err := http.ListenAndServe(":8080", router)
		log.Fatal(err)
		if err != nil {
			fmt.Println("Error al iniciar el servidor:", err)
		}
	}()

	// Espera hasta que se reciba un mensaje en el canal de terminación.
	<-terminateRESTfulAPI
}

func Test() {
	mutex.Lock()
	defer mutex.Unlock()

	fmt.Printf("Teeeeest del servidor sin correr \n")
}

// Handler para obtener los datos actuales. GetAbouts
func GetAbouts(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	aboutPE, _ := mainMemoryInit.About() // Donde se debe reemplzar para obtener todo el json con la info completa
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, aboutPE)
}

// Handler para establecer datos.
func SetData(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	// Intenta decodificar el JSON en la primera estructura.
	var newData1 struct {
		Type     string `json:"type"`
		LastCode bool   `json:"lastCode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&newData1); err == nil {
		// Verifica si el JSON corresponde al primer tipo.
		if newData1.Type == "MESI" && newData1.LastCode {
			SetProtocol("MESI", newData1.LastCode)
			// Envía una actualización a través del websocket.
			BroadcastUpdate()
			return
		}
		if newData1.Type == "MOESI" && newData1.LastCode {
			SetProtocol("MOESI", newData1.LastCode)
			// Envía una actualización a través del websocket.
			BroadcastUpdate()
			return
		}
	}

	var newData2 struct {
		Action string `json:"action"`
	}

	if err := json.NewDecoder(r.Body).Decode(&newData2); err == nil {
		// Verifica si el JSON corresponde al primer tipo.
		if newData2.Action == "close" {
			Close()
			// Envía una actualización a través del websocket.
			BroadcastUpdate()
			return
		}
	}

	var newData3 struct {
		Action string `json:"action"`
		Number string `json:"number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&newData3); err != nil {
		http.Error(w, "JSON no válido", http.StatusBadRequest)
		return
	} else {
		if newData3.Action == "Step" && (newData3.Number == "1" || newData3.Number == "2" || newData3.Number == "3") {
			SetProcessAction("Step", newData3.Number)
			// Envía una actualización a través del websocket.
			BroadcastUpdate()
			return
		}
		if newData3.Action == "Start" {
			SetProcessAction("Start", newData3.Number)
			// Envía una actualización a través del websocket.
			BroadcastUpdate()
			return
		}
	}

	// Envía una actualización a través del websocket.
	BroadcastUpdate()
}

// Función para enviar actualizaciones a través del websocket.
func BroadcastUpdate() {
	mutex.Lock()
	defer mutex.Unlock()

	// Convierte los datos en formato JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error al convertir a JSON:", err)
		return
	}

	// Envía el JSON a todos los clientes conectados a través del websocket.
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			fmt.Println("Error al escribir mensaje:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

var clients = make(map[*websocket.Conn]bool)

// Función para iniciar el servidor WebSocket.
func StartWebSocketServer() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error al actualizar la conexión:", err)
			return
		}

		defer conn.Close()

		// Agrega el cliente a la lista de clientes.
		clients[conn] = true

		// Lee mensajes del cliente (puedes implementar más acciones aquí).
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Error al leer mensaje:", err)
				delete(clients, conn)
				break
			}
		}
	})

	http.ListenAndServe(":8081", nil)
}

// Handler para cerrar el servidor.
func CloseServer(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	// Cierra el canal de terminación para señalar la terminación del servidor.
	close(terminateRESTfulAPI)

	// Envía una respuesta exitosa al cliente.
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Server is closing...")
}

func SetProtocol(Protocol string, lastCode bool) {
	mutex.Lock()
	defer mutex.Unlock()
	InitializeCacheProtocol(Protocol, lastCode)
	fmt.Printf("Protocolo seleccionado %s \n", Protocol)
}

func SetProcessAction(Action string, Process string) {
	mutex.Lock()
	defer mutex.Unlock()

	ProcessAction(Action, Process)

	fmt.Printf("Action: %s Proceso: %s \n", Action, Process)
}

func InitializeCacheProtocol(Protocol string, lastCode bool) {

	protocol := Protocol
	// generate the program files ****************************************************************************************************
	instructionsPerCore := 10
	instructions = utils.GenerateRandomInstructions(3, instructionsPerCore)
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
	terminate = make(chan struct{})

	// Declare the Communication Channels array for PE-CC
	RequestChannelsM1 = make([]chan utils.RequestProcessingElement, 3)
	ResponseChannelsM1 = make([]chan utils.ResponseProcessingElement, 3)

	// Declare the Communication Channels array for CC-IC
	RequestChannelsM2 = make([]chan utils.RequestInterconnect, 3)
	ResponseChannelsM2 = make([]chan utils.ResponseInterconnect, 3)

	// Declare the Broadcast Communication Channels array for CC-IC
	RequestChannelsBroadcast = make([]chan utils.RequestBroadcast, 3)
	ResponseChannelsBroadcast = make([]chan utils.ResponseBroadcast, 3)

	// Declare the Communication Channels for the Interconnect and Main Memory
	RequestChannelM3 = make(chan utils.RequestMainMemory)
	ResponseChannelM3 = make(chan utils.ResponseMainMemory)

	// Create and start 3 Cache Controllers with the communication channels
	ccs = make([]*CacheController.CacheController, 3) // Create an array of Cache Controllers

	semaphore = make(chan struct{}, 1) // Initialize with a count of 1

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
	pes = make([]*processingElement.ProcessingElement, 3) // Create an array of PEs

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
	interconnectInit, err := interconnect.New(
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
		interconnectInit.Run(&wg)
	}()

	// Create Main Memory with two channels, ready to connect the interconect
	mainMemoryInit, err := mainMemory.New(RequestChannelM3, ResponseChannelM3, terminate)
	if err != nil {
		fmt.Printf("Error initializing Main Memory: %v\n", err)
		return
	}
	// Start Main Memory
	wg.Add(1)
	go func() {
		defer wg.Done()
		mainMemoryInit.Run(&wg)
	}()
}

func ProcessAction(Action string, Process string) {
	if Action == "all" {
		for i, pe := range pes {
			if !pe.IsDone && !pe.IsExecutingInstruction {
				pe.Control <- true
				fmt.Printf("Sent 'step' command to PE%d...\n", i)
			} else {
				fmt.Printf("PE%d is not available...\n", pe.ID)
			}
		}
	} else {
		peIndex, err := strconv.Atoi(Process)
		if err != nil || peIndex < -1 || peIndex > len(pes)-1 {
			fmt.Println("Invalid PE number. Please enter a valid PE number or 'all'.")
		}

		pe := pes[peIndex]
		if !pe.IsDone && !pe.IsExecutingInstruction {
			pe.Control <- true
			fmt.Printf("Sent 'step' command to PE%d...\n", pe.ID)
		} else {
			fmt.Printf("PE%d is not available...\n", pe.ID)
		}
	}
}

func Close() {
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
}
