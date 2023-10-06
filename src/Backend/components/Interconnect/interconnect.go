package interconnect

import (
    "log"
    "os"
    "sync"
	"time"
	"encoding/json"
	"fmt"

    "Backend/utils"
)

type Interconnect struct {
    ID              int
    RequestChannelsCacheController []chan utils.RequestInterconnect     	// Request channels from CacheControllers
    ResponseChannelsCacheController []chan utils.ResponseInterconnect  	// Response channels to CacheControllers
	RequestChannelMainMemory chan utils.RequestMainMemory 			// Request channel for Main Memory
	ResponseChannelMainMemory chan utils.ResponseMainMemory		// Response channel for Main Memory	

	RequestChannelsBroadcast []chan utils.RequestBroadcast        // Request channels for Interconnect
	ResponseChannelsBroadcast []chan utils.ResponseBroadcast       // Response channels for Interconnect
    Quit            chan struct{}
	Protocol string
    Logger          *log.Logger
	Transactions utils.QueueS
	Logs		utils.QueueS
	Status string
	PowerConsumption float64
	ReadRequests int
	ReadExclusiveRequests int
	DataResponses int
	Invalidates int
	MemoryReads int
	MemoryWrites int

}

func New(
		requestChannelsCC []chan utils.RequestInterconnect,
		responseChannelsCC []chan utils.ResponseInterconnect, 
		requestChannelMM chan utils.RequestMainMemory,
		responseChannelMM chan utils.ResponseMainMemory,
		requestChannelsCCp []chan utils.RequestBroadcast,
		responseChannelsCCp []chan utils.ResponseBroadcast,
		protocol string,
		quit chan struct{}) (*Interconnect, error) {

    // Create the log file for the Interconnect
    logFile, err := os.Create("logs/IC/IC.log")
    if err != nil {
        log.Fatalf("Error creating log file for Interconnect: %v", err)
    }

    // Initialize logger for the Interconnect using its respective log file
    logger := log.New(logFile, "IC" +"_", log.Ldate|log.Ltime)

	// Create a new queue to store the Bus Transactions
	transactionsQueue := utils.QueueS{}

	// Create a new queue to store the Bus logs
	busQueue := utils.QueueS{}

    return &Interconnect{
        RequestChannelsCacheController: requestChannelsCC,
        ResponseChannelsCacheController: responseChannelsCC,
		RequestChannelMainMemory: requestChannelMM,
		ResponseChannelMainMemory: responseChannelMM,
		RequestChannelsBroadcast: requestChannelsCCp,
		ResponseChannelsBroadcast: responseChannelsCCp,
        Quit:            quit,
		Protocol: protocol,
        Logger:          logger,
		Transactions: transactionsQueue,
		Logs: busQueue,
		Status: "Active",
		PowerConsumption: 0,
		ReadRequests: 0,
		ReadExclusiveRequests: 0,
		DataResponses: 0,
		Invalidates: 0,
    }, nil
}

// Function to get a JSON string with the current state of the Interconnect
func (ic *Interconnect) About()(string, error){
    // Create an empty LogObjectList
    logs := utils.LogObjectList{}

    // Get the values from the transactions queue of the Interconnect
    for i, item := range ic.Logs.Items{
		// Create an LogObject
        logObj := utils.LogObject{
            Order: i,
            Log: item,
        }
        // Append the transaction object to the list
        logs = append(logs, logObj)
	}

    // Create a struct
    aboutIC := utils.AboutInterconnect {
        Status: ic.Status,
		Logs: logs,
    }

	// Marshal the PE struct into a JSON string
	jsonData, err := json.MarshalIndent((aboutIC), "", "    ")
	if err != nil {
		return "", err
	}
	// Convert the byte slice to a string
	jsonString := string(jsonData)

	return jsonString, nil
}
// Function to send a write request to Main Memory
func (ic *Interconnect) WriteToMainMemory(address int, data int) bool{
	// Create a struct for the request
	requestMainMemory := utils.RequestMainMemory{
		Type: "WRITE",
		Address: address,
		Value: uint32(data),
	}
	ic.Logs.Enqueue(fmt.Sprintf("%s - Writing the value %d to memory address %d.", time.Now().Format("15:04:05"), data, address))
	ic.Transactions.Enqueue(time.Now().Format("15:04:05") + "-write-to-memory")
	ic.MemoryWrites++
	ic.PowerConsumption += 3.0

	// Send the request to the Main Memory
	ic.RequestChannelMainMemory <- requestMainMemory
	ic.Logger.Printf(" - IC has sent a WRITE request to the Main Memory.\n")

	// Wait for a response from the Main Memory
	ic.Logger.Printf(" - IC is waiting for a response from the Main Memory.\n")
	responseMainMemory := <- ic.ResponseChannelMainMemory
	requestStatus := responseMainMemory.Status
	ic.Logs.Enqueue(fmt.Sprintf("%s - Finished writing in memory address %d.", time.Now().Format("15:04:05"), address))
	
	return requestStatus
}

// Function to send a read request to Main Memory
func (ic *Interconnect) ReadFromMainMemory(address int) int{
	// Create a struct for the request
	requestMainMemory := utils.RequestMainMemory{
		Type: "READ",
		Address: address,
		Value: uint32(0),
	}
	ic.Logs.Enqueue(fmt.Sprintf("%s - Reading address %d from memory...", time.Now().Format("15:04:05"), address))
	ic.Transactions.Enqueue(time.Now().Format("15:04:05") + "-read-from-memory")
	ic.MemoryReads++
	ic.PowerConsumption += 2.0

	// Send the request to the Main Memory
	ic.RequestChannelMainMemory <- requestMainMemory
	ic.Logger.Printf(" - IC has sent a READ request to the Main Memory.\n")

	// Wait for a response from the Main Memory
	responseMainMemory := <- ic.ResponseChannelMainMemory
	dataResponse := int(responseMainMemory.Value)
	ic.Logs.Enqueue(fmt.Sprintf("%s - Received the value %d from memory.", time.Now().Format("15:04:05"), dataResponse))

	ic.Logger.Printf(" - IC received the value %d from Main Memory.\n", dataResponse)
	return dataResponse
}

// Function to send a DataResponse to an specific Cache Controller
func (ic *Interconnect) SendDataResponseToCacheController(ccID int, data int, status string) {
	// Prepare the data response struct
	dataResponse := utils.ResponseInterconnect{
		Data: data,
		NewStatus: status,
	}

	// Send it to the Cache Controller who requested the data
	ic.ResponseChannelsCacheController[ccID] <- dataResponse
	ic.Logs.Enqueue(fmt.Sprintf("%s - Sent data response to CC%d.", time.Now().Format("15:04:05"), ccID))
	ic.Transactions.Enqueue(time.Now().Format("15:04:05") + "-data-response")
	ic.DataResponses++
	ic.PowerConsumption += 0.8

	ic.Logger.Printf(" - IC sent a data response back to CC%d.\n", ccID)
	ic.Logger.Printf(" - Sent the value %d to CC%d.\n", dataResponse.Data, ccID)
}

// Function to send a Status to an specific Cache Controller
func (ic *Interconnect) SendStatusResponseToCacheController(ccID int, status string) {
	// Prepare the data response struct
	dataResponse := utils.ResponseInterconnect{
		NewStatus: status,
	}

	// Send it to the Cache Controller who requested the data
	ic.ResponseChannelsCacheController[ccID] <- dataResponse
	ic.Logger.Printf(" - IC sent a status response back to CC%d.\n", ccID)
	ic.Logger.Printf(" - Sent a confirmation status to CC%d.\n", ccID)
}

// Function to send a broadcast message to the IDLE Cache Controllers
func (ic *Interconnect) BroadcastMessage(ccID int, requestType string, AR string, address int) (bool, string, int) {
	// Prepare the output values
	Found := false 			// This flag indicates that the data was found
	Status := "I"			// This string represents the final status of the address
	Data := 0				// This int represents the data provided from a remote cache for a data response AR

	BPC1 := 0.0
	BPC2 := 0.0

	var (
		M, O, E, S int
		mu         sync.Mutex // Mutex to protect the counters
	)

	// Prepare a struct for the broadcast message
	broadcastRequest := utils.RequestBroadcast {
		Type: requestType,
		Address: address,
	}
	ic.Logger.Printf(" - IC will send a broadcast message to the CCs.\n")
	ic.Logs.Enqueue(fmt.Sprintf("%s - IC will send a broadcast message to the CCs.", time.Now().Format("15:04:05")))

	// Asign a weight based on the primary request type
	switch requestType {
	case "ReadRequest":
		BPC1 = 1.0
	case "ReadExclusiveRequest":
		BPC1 = 2.0
	}

	// Asign a weight based on the secondary request type
	switch AR {
	case "DataResponse":
		BPC2 = 0.8
	case "Invalidate":
		BPC2 = 1.5
	}

	var wg sync.WaitGroup

	for cc := range ic.RequestChannelsBroadcast {
		// Ask everyone exept the CC the IC is attending
		if (cc == ccID) {continue}

		// Increment the WaitGroup counter
		wg.Add(1)
		
		// Send a goroutine to broadcast the message
		go func(cc int){
			defer wg.Done()

			// Send the broadcast message to all the Cache Controllers
			ic.RequestChannelsBroadcast[cc] <- broadcastRequest
			ic.Logger.Printf(" - IC sent a broadcast %s to CC%d.\n", requestType, cc)
			ic.Logs.Enqueue(fmt.Sprintf("%s - IC sent a broadcast %s to CC%d.", time.Now().Format("15:04:05"), requestType, cc))
		
			// Wait for a response from the Cache Controller
			broadcastResponse := <- ic.ResponseChannelsBroadcast[cc]

			// Process the incoming response
			Matched := broadcastResponse.Match
			BlockStatus := broadcastResponse.Status

			// Update counters inside a critical section
			mu.Lock()
			defer mu.Unlock()

			// For each core broadcast message, sum the request power consumption
			ic.PowerConsumption += BPC1
			ic.PowerConsumption += BPC2

			if !Matched {
				ic.Logger.Printf(" - CC%d doesn't have the data.\n", cc)
				ic.Logs.Enqueue(fmt.Sprintf("%s - CC%d doesn't have the data.", time.Now().Format("15:04:05"), cc))
				return
			}

			switch BlockStatus {
			case "M":
				ic.Logger.Printf(" - CC%d has the data and its status is Modified.\n", cc)
				ic.Logs.Enqueue(fmt.Sprintf("%s - CC%d has the data and its status is 'Modified'.", time.Now().Format("15:04:05"), cc))
				Found = true
				M++
			case "O":
				ic.Logger.Printf(" - CC%d has the data and its status is Owned.\n", cc)
				ic.Logs.Enqueue(fmt.Sprintf("%s - CC%d has the data and its status is 'Owned'.", time.Now().Format("15:04:05"), cc))
				Found = true
				O++
			case "E":
				ic.Logger.Printf(" - CC%d has the data and its status is Exclusive.\n", cc)
				ic.Logs.Enqueue(fmt.Sprintf("%s - CC%d has the data and its status is 'Exclusive'.", time.Now().Format("15:04:05"), cc))
				Found = true
				E++
			case "S":
				ic.Logger.Printf(" - CC%d has the data and its status is Shared.\n", cc)
				ic.Logs.Enqueue(fmt.Sprintf("%s - CC%d has the data and its status is 'Shared'.", time.Now().Format("15:04:05"), cc))
				Found = true
				S++
			}
		}(cc)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Handle the statuses
	if Found {
		if M > 0{
			Status = "M"

		} else if O > 0{
			Status = "O"

		}else if E > 0{
			Status = "E"
			
		}else if S > 0{
			Status = "S"
		}
	}
	// Return the results after the loop
	return Found, Status, Data
}

// Function to handle the requests from a Cache Controller
func (ic *Interconnect) handleRequestFromCC(ccID int, request utils.RequestInterconnect) {
	requestType := request.Type
	requestAddress := request.Address
	requestAR := request.AR
	ic.Logger.Printf(" - IC received a %s from CC%d.\n",requestType, ccID)

	// Get the current time
	currentTime := time.Now()
	timeString := currentTime.Format("15:04:05")

	ic.Logs.Enqueue(fmt.Sprintf("%s - Received a %s from CC%d", timeString, requestType, ccID))

	// Send a broadcast message to the IDLE Cache Controllers
	RemoteFound, RemoteStatus, RemoteData := ic.BroadcastMessage(ccID, requestType, requestAR, requestAddress)
	ic.Logger.Printf(" - IC finished broadcast.\n")
	switch requestType {
	// Handle Read-Request
	case "ReadRequest":
		ic.Transactions.Enqueue(timeString + "-read-request")
		ic.ReadRequests++
		// Add the power consumption for the requesting cache
		ic.PowerConsumption += 1.0

		time.Sleep(3 * time.Second)
		// MESI protocol ***************************************************************************************************************
		if (ic.Protocol == "MESI"){
			// The data was found in a remote cache
			if (RemoteFound){
				// The Action Required is a Data Response
				if (requestAR == "DataResponse"){
					// The data was found with a 'E' or 'S' status
					if (RemoteStatus == "E" || RemoteStatus == "S"){
						// Send the Data provided by the remote cache back to the requesting Cache Controller
						ic.SendDataResponseToCacheController(ccID, RemoteData, "S")
					}
					// The data was found with a 'M' status
					if (RemoteStatus == "M"){
						// Flush the data back to Main Memory
						ic.WriteToMainMemory(requestAddress, RemoteData)
						// Send the Data provided by the remote cache back to the requesting Cache Controller
						ic.SendDataResponseToCacheController(ccID, RemoteData, "S")
					}
				}
			}
			// The data was not found in a remote cache
			if (!RemoteFound){
				// Bring the data from Main Memory
				dataFromMemory := ic.ReadFromMainMemory(requestAddress)
				// Send the data response back to the Cache Controller waiting for a response
				ic.SendDataResponseToCacheController(ccID, dataFromMemory, "E")	
			}
		}
		// MOESI protocol *********************************************************************************************************
		if (ic.Protocol == "MOESI"){
			// The data was found in a remote cache
			if (RemoteFound){
				// The Action Required is a Data Response
				if (requestAR == "DataResponse"){
					if (RemoteStatus == "M" || RemoteStatus == "E" || RemoteStatus == "S"){
						// Send the Data provided by the remote cache back to the requesting Cache Controller
						ic.SendDataResponseToCacheController(ccID, RemoteData, "S")
					}
					// The data was found with a 'O' status
					if (RemoteStatus == "O"){
						// Send the Data provided by the remote cache back to the requesting Cache Controller
						ic.SendDataResponseToCacheController(ccID, RemoteData, "S")
					}
				}
			}
			// The data was not found in a remote cache
			if (!RemoteFound){
				// The Action Required is a Data Response
				if (requestAR == "DataResponse"){
					// Bring the data from Main Memory
					dataFromMemory := ic.ReadFromMainMemory(requestAddress)
					// Send the data response back to the Cache Controller waiting for a response
					ic.SendDataResponseToCacheController(ccID, dataFromMemory, "E")
				}
			}
		}
	// Handle Read-Request
	case "ReadExclusiveRequest":
		ic.Transactions.Enqueue(timeString + "-read-exclusive-request")
		ic.ReadExclusiveRequests++
		// Add the power consumption for the requesting cache
		ic.PowerConsumption += 1.2

		time.Sleep(3 * time.Second)
		// MESI protocol *****************************************************************************************************************
		if (ic.Protocol == "MESI"){
			// The data was found in a remote cache
			if (RemoteFound){
				// The Action Required is a Data Response
				if (requestAR == "DataResponse"){
					// The data was found with a 'E' status
					if (RemoteStatus == "E"){
						// Send the Data provided by the remote cache back to the requesting Cache Controller
						ic.SendDataResponseToCacheController(ccID, RemoteData, "M")
					}
					// The data was found with a 'M' status
					if (RemoteStatus == "M"){
						// Flush the data back to Main Memory
						ic.WriteToMainMemory(requestAddress, RemoteData)
						// Send the Data provided by the remote cache back to the requesting Cache Controller
						ic.SendDataResponseToCacheController(ccID, RemoteData, "M")
					}
				}
				// The Action Required is Invalidate
				if (requestAR == "Invalidate"){
					ic.Transactions.Enqueue(timeString + "-invalidate")
					ic.Invalidates++
					ic.Logs.Enqueue(fmt.Sprintf("%s - Sent Invalidate.", time.Now().Format("15:04:05")))
					// The data was found with a 'S' status
					if (RemoteStatus == "S" || RemoteStatus == "M" || RemoteStatus == "E"){
						// Send the status response back to the requesting Cache Controller
						ic.SendStatusResponseToCacheController(ccID, "M")
					}
				}
			}
			// The data was not found in a remote cache
			if (!RemoteFound){
				// Bring the data from Main Memory
				dataFromMemory := ic.ReadFromMainMemory(requestAddress)
				// Send the data response back to the Cache Controller waiting for a response
				ic.SendDataResponseToCacheController(ccID, dataFromMemory, "M")
			}
		}
		// MOESI protocol *************************************************************************************************
		if (ic.Protocol == "MOESI"){
			// The data was found in a remote cache
			if (RemoteFound){
				// The Action Required is a Data Response
				if (requestAR == "DataResponse"){
					// The data was found with a 'E' status
					if (RemoteStatus == "E"){
						// Send the Data provided by the remote cache back to the requesting Cache Controller
						ic.SendDataResponseToCacheController(ccID, RemoteData, "M")
					}
					// The data was found with a 'M' status
					if (RemoteStatus == "M"){
						// Flush the data back to Main Memory
						ic.WriteToMainMemory(requestAddress, RemoteData)
						// Send the Data provided by the remote cache back to the requesting Cache Controller
						ic.SendDataResponseToCacheController(ccID, RemoteData, "M")
					}
					// The data was found with a 'O' status
					if (RemoteStatus == "O"){
						// Send the Data provided by the remote cache back to the requesting Cache Controller
						ic.SendDataResponseToCacheController(ccID, RemoteData, "M")
					}
				}
				// The Action Required is Invalidate
				if (requestAR == "Invalidate"){
					ic.Transactions.Enqueue(timeString + "-invalidate")
					ic.Invalidates++
					ic.Logs.Enqueue(fmt.Sprintf("%s - Sent Invalidate.", time.Now().Format("15:04:05")))
					// The data was found with a 'S' status
					if (RemoteStatus == "S" || RemoteStatus == "O"){
						// Send the status response back to the requesting Cache Controller
						ic.SendStatusResponseToCacheController(ccID, "M")
					}
				}
			}
			// The data was not found in a remote cache
			if (!RemoteFound){
				// The Action Required is a Data Response
				if (requestAR == "DataResponse"){
					// Bring the data from Main Memory
					dataFromMemory := ic.ReadFromMainMemory(requestAddress)
					// Send the data response back to the Cache Controller waiting for a response
					ic.SendDataResponseToCacheController(ccID, dataFromMemory, "M")
				}
			}
		}
	}	
}


func (ic *Interconnect) Run(wg *sync.WaitGroup) {
	ic.Logger.Printf(" - IC is running.\n")
	ic.Logs.Enqueue(fmt.Sprintf("%s - Ready to handle bus requests.", time.Now().Format("15:04:05")))

	for {
		select {
		case request, ok := <-ic.RequestChannelsCacheController[0]:
			if ok {
				ic.handleRequestFromCC(0, request)
			}

		case request, ok := <-ic.RequestChannelsCacheController[1]:
			if ok {
				ic.handleRequestFromCC(1, request)
			}

		case request, ok := <-ic.RequestChannelsCacheController[2]:
			if ok {
				ic.handleRequestFromCC(2, request)
			}

		// Wait for termination
		case <-ic.Quit:
			ic.Logger.Printf(" - IC has received an external signal to terminate.\n")
			return
		}
	}
}
