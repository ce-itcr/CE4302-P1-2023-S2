package interconnect

import (
    "log"
    "os"
    "sync"
	"time"

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
    }, nil
}


// Function to send a write request to Main Memory
func (ic *Interconnect) WriteToMainMemory(address int, data int) bool{
	// Create a struct for the request
	requestMainMemory := utils.RequestMainMemory{
		Type: "WRITE",
		Address: address,
		Value: uint32(data),
	}
	// Send the request to the Main Memory
	ic.RequestChannelMainMemory <- requestMainMemory
	ic.Logger.Printf(" - IC has sent a WRITE request to the Main Memory.\n")

	// Wait for a response from the Main Memory
	ic.Logger.Printf(" - IC is waiting for a response from the Main Memory.\n")
	responseMainMemory := <- ic.ResponseChannelMainMemory
	requestStatus := responseMainMemory.Status
	ic.Logger.Printf(" - IC received a response from the Main Memory: %v.\n", requestStatus)
	
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
	// Send the request to the Main Memory
	ic.RequestChannelMainMemory <- requestMainMemory
	ic.Logger.Printf(" - IC has sent a READ request to the Main Memory.\n")

	// Wait for a response from the Main Memory
	responseMainMemory := <- ic.ResponseChannelMainMemory
	dataResponse := int(responseMainMemory.Value)

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
	ic.Logger.Printf(" - IC sent a data response back to CC%d.\n", ccID)
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
}

// Function to send a broadcast message to the IDLE Cache Controllers
func (ic *Interconnect) BroadcastMessage(ccID int, requestType string, address int) (bool, string, int) {
	// Prepare the output values
	Found := false 			// This flag indicates that the data was found
	Status := "I"			// This string represents the final status of the address
	Data := 0				// This int represents the data provided from a remote cache for a data response AR

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
		
			// Wait for a response from the Cache Controller
			broadcastResponse := <- ic.ResponseChannelsBroadcast[cc]

			// Process the incoming response
			Matched := broadcastResponse.Match
			BlockStatus := broadcastResponse.Status

			// Update counters inside a critical section
			mu.Lock()
			defer mu.Unlock()

			if !Matched {
				ic.Logger.Printf(" - CC%d doesn't have the data.\n", cc)
				return
			}

			switch BlockStatus {
			case "M":
				ic.Logger.Printf(" - CC%d has the data and its status is Modified.\n", cc)
				Found = true
				M++
			case "O":
				ic.Logger.Printf(" - CC%d has the data and its status is Owned.\n", cc)
				Found = true
				O++
			case "E":
				ic.Logger.Printf(" - CC%d has the data and its status is Exclusive.\n", cc)
				Found = true
				E++
			case "S":
				ic.Logger.Printf(" - CC%d has the data and its status is Shared.\n", cc)
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

// // Function to send a broadcast message to the IDLE Cache Controllers
// func (ic *Interconnect) BroadcastMessage(ccID int, requestType string, address int) (bool, string, int) {
// 	// Prepare the output values
// 	Found := false 			// This flag indicates that the data was found
// 	Status := "I"			// This string represents the final status of the address
// 	Data := 0				// This int represents the data provided from a remote cache for a data response AR

// 	M := 0
// 	O := 0
// 	E := 0
// 	S := 0

// 	// Prepare a struct for the broadcast message
// 	broadcastRequest := utils.RequestBroadcast {
// 		Type: requestType,
// 		Address: address,
// 	}
// 	ic.Logger.Printf(" - IC is about to send a broadcast message to the other CCs.\n")

// 	for cc := range ic.RequestChannelsBroadcast {
// 		// Ask everyone exept the CC the IC is attending
// 		if (cc == ccID) {continue}

// 		// Send the broadcast message to all the Cache Controllers
// 		ic.RequestChannelsBroadcast[cc] <- broadcastRequest
// 		ic.Logger.Printf(" - IC sent a broadcast %s to CC%d.\n", requestType, cc)
	
// 		broadcastResponse := <- ic.ResponseChannelsBroadcast[cc]
// 		Matched := broadcastResponse.Match
// 		BlockStatus := broadcastResponse.Status

// 		// Take the data only if the response was an Intervention (The Data field is a -1)
// 		if (broadcastResponse.Data != -1) {
// 			Data = broadcastResponse.Data
// 		}

// 		// Check if the Cache Controller doesn't have the data, just continue with the other CCs
// 		if (!Matched) {
// 			ic.Logger.Printf(" - CC%d doesn't have the data.\n", cc)
// 			continue
// 		}

// 		// Check if the Cache Controller has the data, and if the Status of the block id 'Modified'.
// 		if (Matched && BlockStatus == "M"){
// 			ic.Logger.Printf(" - CC%d has the data and its status is Modified.\n", cc)
// 			// Set the final startus to Modified and the found flag to true
// 			Found = true
// 			M++
// 			continue

// 		}

// 		// Check if the Cache Controller has the data, and if the Status of the block id 'Modified'.
// 		if (Matched && BlockStatus == "O"){
// 			ic.Logger.Printf(" - CC%d has the data and its status is Owned.\n", cc)
// 			// Set the final startus to Modified and the found flag to true
// 			Found = true
// 			// Break the loop because there are no other CCs containing the data
// 			O++
// 			continue
// 		}

// 		// Check if the Cache Controller has the data, and if the Status of the block id 'Exclusive'.
// 		if (Matched && BlockStatus == "E"){
// 			ic.Logger.Printf(" - CC%d has the data and its status is Exclusive.\n", cc)
// 			// Set the final startus to Modified and the found flag to true
// 			Found = true
// 			// Break the loop because there are no other CCs containing the data
// 			E++
// 			continue
// 		}

// 		// Check if the Cache Controller has the data, and if the Status of the block id 'Shared'.
// 		if (Matched && BlockStatus == "S"){
// 			ic.Logger.Printf(" - CC%d has the data and its status is Shared.\n", cc)
// 			// Set the final startus to Modified and the found flag to true
// 			Found = true
// 			// Continue with the loop because there may be other copies of the data
// 			S++
// 			continue
// 		}

// 	}

// 	// Handle the statuses
// 	if (Found) {
// 		if (M > 0){
// 			Status = "M"

// 		} else if (O > 0){
// 			Status = "O"

// 		}else if (E > 0){
// 			Status = "E"
			
// 		}else if (S > 0 ){
// 			Status = "S"
// 		}
// 	}
// 	// Return the results after the loop
// 	return Found, Status, Data
// }

func (ic *Interconnect) handleRequestFromCC(ccID int, request utils.RequestInterconnect) {
	requestType := request.Type
	requestAddress := request.Address
	requestAR := request.AR
	ic.Logger.Printf(" - IC received a %s from CC%d.\n",requestType, ccID)

	// Send a broadcast message to the IDLE Cache Controllers
	RemoteFound, RemoteStatus, RemoteData := ic.BroadcastMessage(ccID, requestType, requestAddress)
	switch requestType {
	// Handle Read-Request
	case "ReadRequest":
		time.Sleep(3 * time.Second)
		// MESI protocol ***************************************************************************************************************
		if (ic.Protocol == "MESI"){
			// The data was found in a remote cache
			if (RemoteFound){
				// The Action Required is a Data Response
				if (requestAR == "DataResponse"){
					// The data was found with a 'E' status
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
		time.Sleep(3 * time.Second)
		// MESI protocol *****************************************************************************************************************
		if (ic.Protocol == "MESI"){
			// The data was found in a remote cache
			if (RemoteFound){
				// The Action Required is a Data Response
				if (requestAR == "DataResponse"){
					// The data was found with a 'E' status
					// Here add the Invalidate request
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
					// The data was found with a 'S' status
					if (RemoteStatus == "S"){
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
					// Here add the Invalidate request
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
