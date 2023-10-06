package restfulapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"Backend/components/MultiprocessingSystem"
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

	terminateRESTfulAPI chan struct{}

	mps *MultiprocessingSystem.MultiprocessingSystem
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func Restfulapi() {
	// Configura el enrutador.
	router := mux.NewRouter().StrictSlash(true)

	// Ruta para obtener información sobre PE.
	router.HandleFunc("/about", GetAbouts).Methods("GET")
	router.HandleFunc("/aboutmetrics", GetMetrics).Methods("GET")

	// Ruta para establecer datos.
	router.HandleFunc("/setinitialize", SetInitialize).Methods("POST")
	router.HandleFunc("/setaction", SetAction).Methods("POST")
	router.HandleFunc("/setlj", SetLj).Methods("POST")

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

	aboutMps, _ := mps.GetState()
	fmt.Println(aboutMps)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, aboutMps)
}

func GetMetrics(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	aboutMetrics, _ := mps.GetState()
	fmt.Println(aboutMetrics)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, aboutMetrics)
}

// Handler para establecer datos.
func SetInitialize(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	//BroadcastUpdate()

	var (
		newData1 struct {
			Type     string `json:"type"`
			LastCode bool   `json:"lastCode"`
		}
	)

	// Intenta decodificar en newData1
	if err := json.NewDecoder(r.Body).Decode(&newData1); err == nil {
		// Verificar y manejar newData1
		if newData1.Type == "MESI" && newData1.LastCode {
			// Procesar solicitud MESI aquí
			mps = MultiprocessingSystem.Start("MESI", true, 10)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Solicitud MESI procesada exitosamente")
			return
		}
		if newData1.Type == "MOESI" && newData1.LastCode {
			// Procesar solicitud MOESI aquí
			mps = MultiprocessingSystem.Start("MOESI", true, 10)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Solicitud MOESI procesada exitosamente")
			return
		}
	}

	http.Error(w, "JSON no válido", http.StatusBadRequest)
}

// Handler para establecer datos.
func SetAction(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	//BroadcastUpdate()

	var (
		newData3 struct {
			Action string `json:"action"`
			Number string `json:"number"`
		}
	)

	// Intenta decodificar en newData3
	if err := json.NewDecoder(r.Body).Decode(&newData3); err == nil {
		// Verificar y manejar newData3
		peIndex, err := strconv.Atoi(newData3.Number)
		if err != nil || peIndex < -1 || peIndex > len(mps.ProcessingElements)-1 {
			http.Error(w, "Número de PE no válido", http.StatusBadRequest)
			return
		}

		if newData3.Action == "step" {
			// Procesar solicitud de paso aquí
			mps.SteppingProcessingElement(peIndex)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Solicitud Step"+newData3.Number+"procesada exitosamente")
			return
		}
		if newData3.Action == "start" {
			// Procesar solicitud de inicio aquí
			fmt.Println("Execute all done.")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Solicitud ALL procesada exitosamente")
			return
		}
	}

	http.Error(w, "JSON no válido", http.StatusBadRequest)
}

// Handler para establecer datos.
func SetLj(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	//BroadcastUpdate()

	var (
		newData2 struct {
			Action string `json:"action"`
		}
	)

	// Intenta decodificar en newData2
	if err := json.NewDecoder(r.Body).Decode(&newData2); err == nil {
		// Verificar y manejar newData2
		if newData2.Action == "close" {
			// Procesar solicitud de cierre aquí
			mps.Stop()
			fmt.Println("Multiprocessing system se ha terminado")
			fmt.Println("Si deseas volver a ejecutar el programa, inicia nuevamente con el frontEnd")
			fmt.Println("Caso contrario, introduce 'C' para terminar el servidor")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Solicitud cerrar Multiprocessing System procesada exitosamente")
			return
		}
	}

	http.Error(w, "JSON no válido", http.StatusBadRequest)
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

// Handler para cerrar el servidor por http
func CloseServer(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	// Cierra el canal de terminación para señalar la terminación del servidor.
	close(terminateRESTfulAPI)

	// Envía una respuesta exitosa al cliente.
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Server is closing...")
}

// Handler para cerrar el servidor por terminal
func Close() {
	fmt.Println("Sent 'lj' command to terminate the program")
	close(terminateRESTfulAPI)
	fmt.Println("Server is closing...")
}
