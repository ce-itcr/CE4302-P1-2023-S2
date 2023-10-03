package RESTfulAPI

import (
	"encoding/json"
	"fmt"
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

func RESTfulAPI() {
	// Configura el enrutador.
	router := mux.NewRouter()

	// Ruta para obtener los datos actuales.
	router.HandleFunc("/data", GetData).Methods("GET")

	// Ruta para establecer datos.
	router.HandleFunc("/set", SetData).Methods("POST")

	// Inicia el servidor web.
	go StartWebSocketServer()

	// Servidor HTTP.
	http.Handle("/", router)

	// Inicia el servidor HTTP.
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
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

	// Decodifica el JSON recibido en el cuerpo de la solicitud y actualiza los datos.
	var newData Datos
	if err := json.NewDecoder(r.Body).Decode(&newData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data = newData

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
