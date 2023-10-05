package CacheController

import (
	"sync"
	"log"
	"os"
	"strconv"
	"time"
	"encoding/json"
	
	"Backend/utils"
)

type CacheController struct {
	ID int
	Cache *Cache
    RequestChannelProcessingElement chan utils.RequestProcessingElement 
    ResponseChannelProcessingElement chan utils.ResponseProcessingElement

	RequestChannelInterconnect chan utils.RequestInterconnect 
    ResponseChannelInterconnect chan utils.ResponseInterconnect

	RequestChannelBroadcast chan utils.RequestBroadcast
    ResponseChannelBroadcast chan utils.ResponseBroadcast
	Semaphore chan struct {}

	Quit chan struct{}
	Protocol string
	Logger *log.Logger
	ReplacementQueue utils.Queue
	Status string
	CacheHits int
	CacheMisses int
}

func New(
			id int, 
			requestChannelPE chan utils.RequestProcessingElement, 
			responseChannelPE chan utils.ResponseProcessingElement,
			requestChannelIC chan utils.RequestInterconnect,
			responseChannelIC chan utils.ResponseInterconnect,
			requestChannelBroadcast chan utils.RequestBroadcast,
			responseChannelBroadcast chan utils.ResponseBroadcast,
			semaphore chan struct{},
			protocol string,
			quit chan struct{}) (*CacheController, error) {

    // Create the log file
    logFile, err := os.Create("logs/CC/CC" + strconv.Itoa(id) + ".log")
    if err != nil {
        log.Fatalf("Error creating log file for CC%d: %v", id, err)
    }

    // Initialize logger for the PE using its respective log file
    logger1 := log.New(logFile, "CC" + strconv.Itoa(id) + "_", log.Ldate|log.Ltime)

	// Create a new queue to handle the cache lines replacement
	myQueue := utils.Queue{}
	
	return &CacheController{
		ID: id,
		Cache: NewCache(),
		Logger: logger1,
        RequestChannelProcessingElement: requestChannelPE,
        ResponseChannelProcessingElement: responseChannelPE,
		RequestChannelInterconnect: requestChannelIC,
		ResponseChannelInterconnect: responseChannelIC,
		RequestChannelBroadcast: requestChannelBroadcast,
		ResponseChannelBroadcast: responseChannelBroadcast,
		Semaphore: semaphore,
		Protocol: protocol,
		Quit: quit,
		ReplacementQueue: myQueue,
		Status: "Active",
		CacheMisses: 0,
		CacheHits: 0,
    }, nil
}

// Function to obtain all the information about the Cache Controller
func (cc *CacheController) About()(string, error){
	// Create an empty CacheObjectList
	cacheBlocks := utils.CacheObjectList{}
    for i := 0; i <= 3; i++ {
        // Create a new CacheObject instance
		cacheObj := utils.CacheObject{
			Block:   i,
			Address: cc.Cache.GetAddress(i),
			Data:    cc.Cache.GetData(i),
			State:   cc.Cache.GetState(i),
		}
		// Append the new CacheObject to the CacheObjectList
		cacheBlocks = append(cacheBlocks, cacheObj)
    }

    // Create a the final JSON struct
    aboutCC := utils.AboutCacheController {
		ID: cc.ID,
		Status: cc.Status,
		Cache: cacheBlocks,
		CacheMisses: cc.CacheMisses,
		CacheHits: cc.CacheHits,
	}

	// Marshal the PE struct into a JSON string
	jsonData, err := json.MarshalIndent(aboutCC, "", "    ")
	if err != nil {
		return "", err
	}

	// Convert the byte slice to a string
	jsonString := string(jsonData)

	return jsonString, nil
}

// Function to return the status of an address in the local cache
func (cc *CacheController) GetAddressStatus(address int) string{
	pos := 0
	for _, addr := range cc.Cache.address{
		if(address == addr){
			break
		}
		pos++
	}
	if(pos == 4){
		return "I"
	}
	return cc.Cache.status[pos]
}

// Function to know if a data is in the local cache
func (cc *CacheController) DataInCache(address int) bool{
	found := false
	for _, addr := range cc.Cache.address{
		if(address == addr){
			found = true
			break
		}
	}
	return found 
}

// Function to write a new data in the local cache, using FIFO for replacement policy
func (cc *CacheController) WriteDataToCache(address int, data int, status string){
	cc.Logger.Printf(" - CC%d is storing the value %d into the local cache at the address %d.\n", cc.ID, data, address)
	// Get the new block to replace
	newLine := cc.ReplacementQueue.Size()

	// Check if the data is in the local cache
	if (cc.DataInCache(address)){
		cc.Logger.Printf(" - The address exists in the local cache.\n")
		pos := 0
		for _, addr := range cc.Cache.address{
			if(address == addr){
				break
			}
			pos++
		}
		newLine = pos
	}

	if (!cc.DataInCache(address)){
		cc.Logger.Printf(" - The address doesn't exist in the local cache.\n")
		if (newLine == 4){
			newLine = cc.ReplacementQueue.Dequeue()
			cc.Logger.Printf(" - CC%d is replacing the the block %d.\n", cc.ID, newLine)
		}
	}
	
	// Replace the contents of the cache line
	cc.Cache.SetData(newLine, data)
	cc.Cache.SetAddress(newLine, address)
	cc.Cache.SetState(newLine, status)

	// Add the cache line to the queue
	cc.ReplacementQueue.Enqueue(newLine)
	cc.Logger.Printf(" - CC%d stored the value %d at the memory address %d and the cache block %d.\n", cc.ID, data, address, newLine)
	cc.Logger.Printf(" - The new state of address %d is '%s'.\n", address, status)
}

// Function to change the status of a local cache line
func (cc *CacheController) ChangeCacheLineStatus(address int, newStatus string) bool {
	found := false
	cacheLine := 0
	for _, addr := range cc.Cache.address{
		if(address == addr){
			found = true
			break
		}
		cacheLine++
	}
	// If the address is not in the cache
	if(!found){
		return false
	}

	// Change the status of that cache line
	cc.Cache.SetState(cacheLine, newStatus)
	cc.Logger.Printf(" - CC%d changed the status of the cache block %d to %s.\n", cc.ID, cacheLine, newStatus)

	return true
}

// Function to get a data from a local cache address
func (cc *CacheController) GetDataFromCache(address int) int{
	cacheLine := 0
	for _, addr := range cc.Cache.address{
		if(address == addr){
			break
		}
		cacheLine++
	}

	// Return the data at the required address
	return cc.Cache.GetData(cacheLine)
}

// Function to send a read-request to the Interconnect
func (cc *CacheController) RequestToInterconnect(requestType string, AR string, address int) (int, string){
	// Prepsre a struct for the request
	readRequest := utils.RequestInterconnect {
		Type: requestType,
		AR: AR,
		Address: address,
	}
	// Wait for 2 seconds
	time.Sleep(2 * time.Second)

	// Send the request to the Interconnect
	cc.Logger.Printf(" - CC%d is about to send a %s to the Interconnect.\n", cc.ID, requestType)
	cc.RequestChannelInterconnect <- readRequest
	cc.Logger.Printf(" - CC%d sent a %s to the Interconnect\n", cc.ID, requestType)

	cc.Logger.Printf(" - CC%d is waiting for the data from the Interconnect.\n", cc.ID)
	dataResponse := <- cc.ResponseChannelInterconnect
	cc.Logger.Printf(" - CC%d received a response from the Interconnect with the value: %d.\n", cc.ID, dataResponse.Data)
	Data := dataResponse.Data
	NewStatus := dataResponse.NewStatus

	// Return the response data
	return Data, NewStatus
}

// Function to send a response to the Processing Element
func (cc *CacheController) RespondToProcessingElement(Data int, Status bool) {
	// Prepare a structure for the response
	peResponse := utils.ResponseProcessingElement {
		Data: Data,
		Status: Status,
	}
	time.Sleep(time.Second)
	cc.Logger.Printf(" - CC%d will send a response to the PE.\n", cc.ID)
	cc.ResponseChannelProcessingElement <- peResponse
	cc.Logger.Printf(" - CC%d sent the response to the PE.\n", cc.ID)
}

// Function to respond to a Broadcast Message from the Interconnect
func (cc *CacheController) RespondToBroadcast(Match bool, Status string, Data int) {
	// Prepare a struct to respond to the broadcast message
	statusResponse := utils.ResponseBroadcast{
		Match: Match,
		Status: Status,
		Data: Data,
	}
	// Send the response to the broadcast
	cc.ResponseChannelBroadcast <- statusResponse
	cc.Logger.Printf(" - CC%d responded to the Broadcast Message with Match: %v, Data: %d, Status: %s.\n", cc.ID, Match, Data, Status)
}

// This is the function that is executed in parallel to handle the requests from the Processing Element
func (cc *CacheController) Run(wg *sync.WaitGroup) {
	cc.Logger.Printf(" - CC%d is running.\n", cc.ID)

	//Create a separate goroutine to handle requests from PE (ChannelM1)
	go func() {
		for {
			select {
			case <- cc.Quit:
				return

			case cc.Semaphore <- struct{}{}:
				go func() {
					select {
						// Listen for the requests from the Processing Element
						case request := <-cc.RequestChannelProcessingElement:
							cc.Logger.Printf(" - CC%d received a request from PE%d.\n", cc.ID, cc.ID)
							requestAddress := request.Address
							requestData := request.Data
							cacheLineStatus := cc.GetAddressStatus(requestAddress)

							switch request.Type {
							// Read request from the Processing Element
							case "READ":
								cc.Logger.Printf(" - CC%d is processing a READ request.\n", cc.ID)
								cc.Logger.Printf(" - The Address to read from is: %d.\n", request.Address)

								// MESI protocol ******************************************************************************
								if (cc.Protocol == "MESI") {
									// The address is not in the local cache
									if (cacheLineStatus == "I"){
										cc.CacheMisses++
										cc.Logger.Printf(" - The address %d is not in the local cache.\n", requestAddress)
				
										// Send a Read-Request to the Interconnect
										Data, NewStatus := cc.RequestToInterconnect("ReadRequest", "DataResponse", requestAddress)
								
										// Update the cache line with the new data and line status
										cc.WriteDataToCache(requestAddress, Data, NewStatus)
				
										// Send a response status to the Processing Element
										cc.RespondToProcessingElement(Data, true)
										
										// Release the semaphore
										<-cc.Semaphore
									}
									// The address is in the local cache
									if (cacheLineStatus == "E" || cacheLineStatus == "M" || cacheLineStatus == "S") {
										cc.CacheHits++
										cc.Logger.Printf(" - The address %d is in the local cache.\n", requestAddress)
										cc.Logger.Printf(" - Communication with the Interconnect is no required.\n")

										// Get the data from the local cache
										DataFromCache := cc.GetDataFromCache(requestAddress)

										// Send the local copy to the Processing Element
										cc.RespondToProcessingElement(DataFromCache, true)
							
										// Release the semaphore
										<-cc.Semaphore
									}
								}

								// MOESI protocol *********************************************************************************
								if (cc.Protocol == "MOESI") {
									// The address is not in the local cache
									if (cacheLineStatus == "I"){
										cc.CacheMisses++
										cc.Logger.Printf(" - The address %d is not in the local cache.\n", requestAddress)
				
										// Send a Read-Request to the Interconnect
										Data, NewStatus := cc.RequestToInterconnect("ReadRequest", "DataResponse", requestAddress)
								
										// Update the cache line with the new data and line status
										cc.WriteDataToCache(requestAddress, Data, NewStatus)
				
										// Send a response status to the Processing Element
										cc.RespondToProcessingElement(Data, true)
										
										// Release the semaphore
										<-cc.Semaphore
									}
									// The address is in the local cache
									if (cacheLineStatus == "E" || cacheLineStatus == "M" || cacheLineStatus == "O") {
										cc.CacheHits++
										cc.Logger.Printf(" - The address %d is in the local cache.\n", requestAddress)
										cc.Logger.Printf(" - Communication with the Interconnect is no required.\n")

										// Get the data from the local cache
										DataFromCache := cc.GetDataFromCache(requestAddress)

										// Send the local copy to the Processing Element
										cc.RespondToProcessingElement(DataFromCache, true)
							
										// Release the semaphore
										<-cc.Semaphore
									}

									// If the address status is Shared
									if (cacheLineStatus == "S") {
										cc.CacheHits++
										cc.Logger.Printf(" - The address %d is in the local cache.\n", requestAddress)

										// Send a Read-Request to the Interconnect
										Data, NewStatus := cc.RequestToInterconnect("ReadRequest", "DataResponse", requestAddress)

										// Update the cache line with the new data and line status
										cc.WriteDataToCache(requestAddress, Data, NewStatus)

										// Send the local copy to the Processing Element
										cc.RespondToProcessingElement(Data, true)
									}
								}
			
							// Write request from the Processing Element
							case "WRITE":
								cc.Logger.Printf(" - CC%d is processing a WRITE request.\n", cc.ID)
								cc.Logger.Printf(" - The Address to write to is: %d.\n", request.Address)
							
								// MESI protocol ************************************************************************************
								if (cc.Protocol == "MESI"){
									// The address is not in the local cache
									if (cacheLineStatus == "I"){
										cc.CacheMisses++
										cc.Logger.Printf(" - The address %d is not in the local cache.\n", requestAddress)
				
										// Send a Read-Exclusive-Request to the Interconnect
										Data, NewStatus := cc.RequestToInterconnect("ReadExclusiveRequest", "Invalidate", requestAddress)
									
										// Update the cache line with the new data and line status
										cc.WriteDataToCache(requestAddress, requestData, NewStatus)
				
										// Send the the status to the Processing Element
										cc.RespondToProcessingElement(Data, true)
							
										// Release the semaphore
										<-cc.Semaphore
									}
				
				
									// The address is in the local cache and its status is 'Shared'
									if (cacheLineStatus == "S"){
										cc.CacheHits++
										cc.Logger.Printf(" - The address %d is in the local cache.\n", requestAddress)
				
										// Send a Read-Exclusive-Request to the Interconnect
										Data, NewStatus := cc.RequestToInterconnect("ReadExclusiveRequest", "Invalidate", requestAddress)
									
										// Update the cache line with the new data and line status
										cc.WriteDataToCache(requestAddress, requestData, NewStatus)
				
										// Send the the status to the Processing Element
										cc.RespondToProcessingElement(Data, true)
							
										// Release the semaphore
										<-cc.Semaphore
									}
				
									// The address is in the local cache and its status is 'Exclusive'
									if (cacheLineStatus == "E" || cacheLineStatus == "M"){
										cc.CacheHits++
										cc.Logger.Printf(" - The address %d is in the local cache.\n", requestAddress)
										cc.Logger.Printf(" - The new data can be writen without using the Interconnect.\n")

										// Write the data in the local cache
										cc.WriteDataToCache(requestAddress, requestData, "M")
				
										// Send the the status to the Processing Element
										cc.RespondToProcessingElement(requestData, true)
							
										// Release the semaphore
										<-cc.Semaphore
									}
								}

								// MOESI protocol **********************************************************************************************
								if (cc.Protocol == "MOESI") {
									// The address is not in the local cache
									if (cacheLineStatus == "I"){
										cc.CacheMisses++
										cc.Logger.Printf(" - The address %d is not in the local cache.\n", requestAddress)
				
										// Send a Read-Exclusive-Request to the Interconnect
										Data, NewStatus := cc.RequestToInterconnect("ReadExclusiveRequest", "DataResponse", requestAddress)
									
										// Update the cache line with the new data and line status
										cc.WriteDataToCache(requestAddress, requestData, NewStatus)
				
										// Send the the status to the Processing Element
										cc.RespondToProcessingElement(Data, true)
							
										// Release the semaphore
										<-cc.Semaphore
									}
				
				
									// The address is in the local cache and its status is 'Shared'
									if (cacheLineStatus == "S"){
										cc.CacheHits++
										cc.Logger.Printf(" - The address %d is in the local cache.\n", requestAddress)
				
										// Send a Read-Exclusive-Request to the Interconnect
										_, NewStatus := cc.RequestToInterconnect("ReadExclusiveRequest", "Invalidate", requestAddress)
									
										// Update the cache line with the new data and line status
										cc.WriteDataToCache(requestAddress, requestData, NewStatus)
				
										// Send the the status to the Processing Element
										cc.RespondToProcessingElement(requestData, true)
							
										// Release the semaphore
										<-cc.Semaphore
									}
				
									// The address is in the local cache and its status is 'Exclusive'
									if (cacheLineStatus == "E" || cacheLineStatus == "M"){
										cc.CacheHits++
										cc.Logger.Printf(" - The address %d is in the local cache.\n", requestAddress)
										cc.Logger.Printf(" - The new data can be writen without using the Interconnect.\n")

										// Write the data in the local cache
										cc.WriteDataToCache(requestAddress, requestData, "M")
				
										// Send the the status to the Processing Element
										cc.RespondToProcessingElement(requestData, true)
							
										// Release the semaphore
										<-cc.Semaphore
									}

									// The address is in the local cache and its status is 'Exclusive'
									if (cacheLineStatus == "O"){
										cc.CacheHits++
										cc.Logger.Printf(" - The address %d is in the local cache.\n", requestAddress)
										cc.Logger.Printf(" - The new data can be writen without using the Interconnect.\n")

										// Send a Upgrade-Request to the Interconnect
										_, NewStatus := cc.RequestToInterconnect("ReadExclusiveRequest", "Invalidate", requestAddress)

										// Update the cache line with the new data and line status
										cc.WriteDataToCache(requestAddress, requestData, NewStatus)
				
										// Send the the status to the Processing Element
										cc.RespondToProcessingElement(requestData, true)
							
										// Release the semaphore
										<-cc.Semaphore
									}
								}

							}

						// Release the semaphore if the Cache Controller acquired it but didn't use the Interconnect
						case <-time.After(time.Millisecond * 200):
							<-cc.Semaphore
			
						case <-cc.Quit:
							return
					}
				}()
			}
	
		}
	}()

	// Create a separate goroutine to handle the broadcast messages from the Interconnect
	go func() {
		for {
			select {
			case broadcastRequest := <-cc.RequestChannelBroadcast:
				address := broadcastRequest.Address
				dataInCache := cc.DataInCache(address)
				addressStatus := cc.GetAddressStatus(address)
				Type := broadcastRequest.Type
				cc.Logger.Printf(" - CC%d received a broadcast message from Interconnect for the address %d.\n", cc.ID, address)

				cc.Logger.Printf(" - The type of request is %s.\n", Type)

				// Check if the request from the broadcast is read-request
				if (Type == "ReadRequest"){
					// MESI protocol ********************************************************************************
					if (cc.Protocol == "MESI"){
						// The data does not exist in the local cache
						if (addressStatus == "I"){
							cc.Logger.Printf(" - The data is not in the local cache.\n")
							// Tell the Interconnect that this cache does not have the data
							cc.RespondToBroadcast(false, addressStatus, -1)
						}
						// The data is in the local cache
						if (dataInCache) {
							Data := cc.GetDataFromCache(address)
							cc.Logger.Printf(" - The data is in the local cache.\n")
							// Tell the Interconnect that this cache has the data
							cc.RespondToBroadcast(true, addressStatus, Data)
							
							// Change the cache line status to Invalid
							cc.ChangeCacheLineStatus(address, "S")
						
						}
					}

					// MOESI protocol *******************************************************************************
					if (cc.Protocol == "MOESI"){
						// The data does not exist in the local cache
						if (addressStatus == "I"){
							cc.Logger.Printf(" - The data is not in the local cache.\n")
							// Tell the Interconnect that this cache does not have the data
							cc.RespondToBroadcast(false, addressStatus, -1)

						}else if (dataInCache) {
							Data := cc.GetDataFromCache(address)
							cc.Logger.Printf(" - The data is in the local cache.\n")
							
							// Check if the addres status is Modified
							if (addressStatus == "M") {
								cc.RespondToBroadcast(true, addressStatus, Data)				
								// Change the cache line status to Owned
								cc.ChangeCacheLineStatus(address, "O")
								
							} else if (addressStatus == "O") {
								cc.RespondToBroadcast(true, addressStatus, Data)
					
							} else if (addressStatus == "E") {
								cc.RespondToBroadcast(true, addressStatus, Data)
								// Change the cache line status to Shared
								cc.ChangeCacheLineStatus(address, "S")
					
							} else if (addressStatus == "S") {
								cc.RespondToBroadcast(true, addressStatus, -1)

							}
						}
					}
				}

				// Check if the request from the broadcast is read-exclusive-request
				if (Type == "ReadExclusiveRequest"){
					// MESI protocol ********************************************************************************************
					if (cc.Protocol == "MESI"){
						// Verify if the data does not exist in the local cache
						if (addressStatus == "I"){
							// Tell the Interconnect that this cache does not have the data
							cc.RespondToBroadcast(false, addressStatus, -1)
						}
						// Verify if the data is in the local cache
						if (dataInCache) {
							Data := cc.GetDataFromCache(address)
							cc.Logger.Printf(" - The data is in the local cache.\n")
							// Tell the Interconnect that this cache has the data
							cc.RespondToBroadcast(true, addressStatus, Data)

							// Change the cache line status to Invalid
							cc.ChangeCacheLineStatus(address, "I")
						}
					}

					// MOESI protocol ******************************************************************************************
					if (cc.Protocol == "MOESI"){
						// Verify if the data does not exist in the local cache
						if (addressStatus == "I"){
							// Tell the Interconnect that this cache does not have the data
							cc.RespondToBroadcast(false, addressStatus, 0)

						}else if (dataInCache) {
							Data := cc.GetDataFromCache(address)
							cc.Logger.Printf(" - The data is in the local cache.\n")
							// Tell the Interconnect that this cache has the data
							cc.RespondToBroadcast(true, addressStatus, Data)

							// Change the cache line status to Invalid
							cc.ChangeCacheLineStatus(address, "I")
						}
					}
		
				}

			case <-cc.Quit:
				return
			}
		}
	}()

	// Wait for termination
	<-cc.Quit
	cc.Logger.Printf(" - CC%d has received an external signal to terminate.\n", cc.ID)
}


