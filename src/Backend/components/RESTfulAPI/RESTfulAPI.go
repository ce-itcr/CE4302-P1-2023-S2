package restfulapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func Restfulapi() {
	// Configura el enrutador.
	router := mux.NewRouter().StrictSlash(true)

	// Ruta para obtener los datos actuales.
	router.HandleFunc("/data", GetData).Methods("GET")

	// Ruta para establecer datos.
	router.HandleFunc("/set", SetData).Methods("POST")

	// Inicia el servidor web.
	go StartWebSocketServer()

	// Servidor HTTP.
	router.HandleFunc("/", homeLink)

	// Inicia el servidor HTTP.
	err := http.ListenAndServe(":8080", router)
	log.Fatal(err)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}

func Test() {
	mutex.Lock()
	defer mutex.Unlock()

	fmt.Printf("Teeeeest del servidor sin correr \n")
}

// Handler para obtener los datos actuales.
func GetData(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	// Responde con los datos actuales en formato JSON.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
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
			SetProtocol("MESI")
			// Envía una actualización a través del websocket.
			BroadcastUpdate()
			return
		}
		if newData1.Type == "MOESI" && newData1.LastCode {
			SetProtocol("MOESI")
			// Envía una actualización a través del websocket.
			BroadcastUpdate()
			return
		}
	}

	// Si no coincide con el primer tipo, intenta decodificarlo en la segunda estructura.
	var newData2 []struct {
		Number string   `json:"number"`
		Code   []string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&newData2); err == nil {
		// Verifica el contenido del JSON y ejecuta la función si es necesario.
		for _, item := range newData2 {
			if item.Number == "pe1" {
				SetProcessCode(item.Number, item.Code)
				BroadcastUpdate()
				return
			} else if item.Number == "pe2" {
				SetProcessCode(item.Number, item.Code)
				BroadcastUpdate()
				return
			} else if item.Number == "pe3" {
				SetProcessCode(item.Number, item.Code)
				BroadcastUpdate()
				return
			}
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
		if newData3.Action == "Step" && (newData3.Number == "pe1" || newData3.Number == "pe2" || newData3.Number == "pe3") {
			SetProcessAction("Step", newData3.Number)
			// Envía una actualización a través del websocket.
			BroadcastUpdate()
			return
		}
		if newData3.Action == "Start" && (newData3.Number == "pe1" || newData3.Number == "pe2" || newData3.Number == "pe3") {
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

func SetProtocol(Protocol string) {
	mutex.Lock()
	defer mutex.Unlock()

	fmt.Printf("Protocolo seleccionado %s \n", Protocol)
}

func SetProcessCode(Process string, code []string) {
	mutex.Lock()
	defer mutex.Unlock()

	fmt.Printf("Proceso: %s Codigo: %q \n", Process, code)
}

func SetProcessAction(Action string, Process string) {
	mutex.Lock()
	defer mutex.Unlock()

	fmt.Printf("Action: %s Proceso: %s \n", Action, Process)
}
