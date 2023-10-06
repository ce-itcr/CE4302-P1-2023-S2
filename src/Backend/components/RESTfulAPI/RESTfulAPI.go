package restfulapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"Backend/components/MultiprocessingSystem"
)

var (
	mutex               sync.Mutex
	BroadcastData       string
	terminateRESTfulAPI chan struct{}
	mps                 *MultiprocessingSystem.MultiprocessingSystem
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func Restfulapi() {
	// Configura el enrutador.
	router := mux.NewRouter().StrictSlash(true)

	// Enable CORS middleware
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Ruta para obtener información sobre PE.
	router.HandleFunc("/about", GetAbouts).Methods("GET")
	router.HandleFunc("/aboutmetrics", GetMetrics).Methods("GET")

	// Ruta para establecer datos.
	router.HandleFunc("/setinitialize", SetInitialize).Methods("POST")
	router.HandleFunc("/setaction", SetAction).Methods("POST")
	router.HandleFunc("/setlj", SetLj).Methods("POST")

	// Servidor HTTP con CORS middleware.
	handlerWithCORS := handlers.CORS(originsOk, headersOk, methodsOk)(router)

	terminateRESTfulAPI = make(chan struct{})

	// Inicia el servidor HTTP en una goroutine.
	go func() {
		err := http.ListenAndServe(":8080", handlerWithCORS)
		log.Fatal(err)
		if err != nil {
			fmt.Println("Error al iniciar el servidor:", err)
		}
	}()

	// Espera hasta que se reciba un mensaje en el canal de terminación.
	<-terminateRESTfulAPI
}

// ... (rest of your code remains unchanged)


// ... (rest of your code remains unchanged)



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

	if mps.AreWeFinished() {
		fmt.Println("The Multiprocessing System has already finished.")
		aboutMetrics, _ := mps.AboutResults()
		fmt.Println(aboutMetrics)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, aboutMetrics)
	} else {
		fmt.Println("The Multiprocessing System has not finished yet.")
		aboutMetrics := "false"
		fmt.Println(aboutMetrics)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, aboutMetrics)
	}
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
		if newData1.Type == "MESI"  {
			// Procesar solicitud MESI aquí
			mps = MultiprocessingSystem.Start("MESI", newData1.LastCode, 4)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Solicitud MESI procesada exitosamente")
			return
		}
		if newData1.Type == "MOESI"{
			// Procesar solicitud MOESI aquí
			mps = MultiprocessingSystem.Start("MOESI", newData1.LastCode, 4)
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
		if newData3.Action == "all" {
			// Procesar solicitud de inicio aquí
			mps.StartProcessingElements()
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