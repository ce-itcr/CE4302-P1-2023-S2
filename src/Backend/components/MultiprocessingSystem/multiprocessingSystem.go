package MultiprocessingSystem

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"Backend/components/CacheController"
	interconnect "Backend/components/Interconnect"
	mainMemory "Backend/components/MainMemory"
	processingElement "Backend/components/ProcessingElement"
	"Backend/utils"
)

// Create a structure that holds all the references needed to control the multiprocessing system
type MultiprocessingSystem struct {
	CacheControllers          []*CacheController.CacheController
	ProcessingElements        []*processingElement.ProcessingElement
	Interconnect              *interconnect.Interconnect
	MainMemory                *mainMemory.MainMemory
	Terminate                 chan struct{}
	WG                        *sync.WaitGroup
	RequestChannelsM1         []chan utils.RequestProcessingElement
	ResponseChannelsM1        []chan utils.ResponseProcessingElement
	RequestChannelsM2         []chan utils.RequestInterconnect
	ResponseChannelsM2        []chan utils.ResponseInterconnect
	RequestChannelsBroadcast  []chan utils.RequestBroadcast
	ResponseChannelsBroadcast []chan utils.ResponseBroadcast
	RequestChannelM3          chan utils.RequestMainMemory
	ResponseChannelM3         chan utils.ResponseMainMemory
	Semaphore                 chan struct{}
}

// Function that initializes a new Multiprocessing System
func Start(Protocol string, CodeGenerator bool, InstructionsPerCore int) *MultiprocessingSystem {
	fmt.Println("Starting a new Multiprocessing System...")
	fmt.Printf("Initializing %s protocol...\n\n\n", Protocol)
	// Is it necessary to generate a random program for the Processing Elements??
	if CodeGenerator {
		instructions := utils.GenerateRandomInstructions(3, InstructionsPerCore)
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
	} else {
		fmt.Printf("Reusing the previous generated code \n")
		for i := 0; i < 3; i++ {
			filename := fmt.Sprintf("generated-programs/program%d.txt", i)
			isEmpty, err := FileIsEmpty(filename)
			if err != nil {
				fmt.Printf("Error reading file.")
			}
			// Check if the file has no instructions
			if isEmpty {
				fmt.Printf("generated-programs/program%d.txt is not valid\n", i)
			}

		}
	}

	// Create termination channel to signal the termination to all threads
	terminate := make(chan struct{})

	// Create WaitGroup for PEs and CCs
	var wg sync.WaitGroup

	// Declare the Communication Channels array for PE-CC
	RequestChannelsM1 := make([]chan utils.RequestProcessingElement, 3)
	ResponseChannelsM1 := make([]chan utils.ResponseProcessingElement, 3)

	// Declare the Communication Channels array for CC-IC
	RequestChannelsM2 := make([]chan utils.RequestInterconnect, 3)
	ResponseChannelsM2 := make([]chan utils.ResponseInterconnect, 3)

	// Declare the Broadcast Communication Channels array for CC-IC
	RequestChannelsBroadcast := make([]chan utils.RequestBroadcast, 3)
	ResponseChannelsBroadcast := make([]chan utils.ResponseBroadcast, 3)

	// Declare the Communication Channels for the Interconnect and Main Memory
	RequestChannelM3 := make(chan utils.RequestMainMemory)
	ResponseChannelM3 := make(chan utils.ResponseMainMemory)

	// Create and start 3 Cache Controllers with the communication channels
	ccs := make([]*CacheController.CacheController, 3) // Create an array of Cache Controllers

	semaphore := make(chan struct{}, 1) // Initialize with a count of 1

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
			Protocol,
			terminate)
		if err != nil {
			fmt.Printf("Error initializing CacheController %d: %v\n", i+1, err)
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
	pes := make([]*processingElement.ProcessingElement, 3) // Create an array of PEs

	for i := 0; i < 3; i++ {
		pe, err := processingElement.New(i, RequestChannelsM1[i], ResponseChannelsM1[i], fmt.Sprintf("generated-programs/program%d.txt", i), terminate)
		if err != nil {
			fmt.Printf("Error initializing ProcessingElement %d: %v\n", i+1, err)
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
	interconnect, err := interconnect.New(
		RequestChannelsM2,
		ResponseChannelsM2,
		RequestChannelM3,
		ResponseChannelM3,
		RequestChannelsBroadcast,
		ResponseChannelsBroadcast,
		Protocol,
		terminate)
	if err != nil {
		fmt.Printf("Error initializing Interconnect: %v\n", err)
	}

	// Start Interconnect
	wg.Add(1)
	go func() {
		defer wg.Done()
		interconnect.Run(&wg)
	}()

	// Create Main Memory with two channels, ready to connect the interconect
	mainMemory, err := mainMemory.New(RequestChannelM3, ResponseChannelM3, terminate)
	if err != nil {
		fmt.Printf("Error initializing Main Memory: %v\n", err)
	}
	// Start Main Memory
	wg.Add(1)
	go func() {
		defer wg.Done()
		mainMemory.Run(&wg)
	}()

	return &MultiprocessingSystem{
		CacheControllers:          ccs,
		ProcessingElements:        pes,
		Interconnect:              interconnect,
		MainMemory:                mainMemory,
		Terminate:                 terminate,
		WG:                        &wg,
		RequestChannelsM1:         RequestChannelsM1,
		ResponseChannelsM1:        ResponseChannelsM1,
		RequestChannelsM2:         RequestChannelsM2,
		ResponseChannelsM2:        ResponseChannelsM2,
		RequestChannelsBroadcast:  RequestChannelsBroadcast,
		ResponseChannelsBroadcast: ResponseChannelsBroadcast,
		RequestChannelM3:          RequestChannelM3,
		ResponseChannelM3:         ResponseChannelM3,
		Semaphore:                 semaphore,
	}
}

// Function to create a JSON object with all the information of the Multiprocessing System
func (mps *MultiprocessingSystem) GetState() (string, error) {

	// Create the AboutProcessingElementList
	pes := utils.AboutProcessingElementList{}
	for _, pe := range mps.ProcessingElements {
		// Create an empty InstructionObjectList
		instructions := utils.InstructionObjectList{}

		// Get the values from the instructions queue of the Processing Element
		for i, item := range pe.Instructions.Items {
			// Create an InstructionObject
			instructionObj := utils.InstructionObject{
				Position:    i,
				Instruction: item,
			}
			// Append the instruction object to the list
			instructions = append(instructions, instructionObj)
		}
		// Create a struct for the PE
		aboutPE := utils.AboutProcessingElement{
			ID:           pe.ID,
			Register:     pe.Register,
			Status:       pe.Status,
			Instructions: instructions,
		}

		// Add it to the list
		pes = append(pes, aboutPE)
	}

	// Create the AboutCacheControllerList
	ccs := utils.AboutCacheControllerList{}
	for _, cc := range mps.CacheControllers {
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
		aboutCC := utils.AboutCacheController{
			ID:     cc.ID,
			Status: cc.Status,
			Cache:  cacheBlocks,
		}
		// Add it to the list
		ccs = append(ccs, aboutCC)
	}

	// Create an empty LogObjectList
	logs := utils.LogObjectList{}

	// Get the values from the transactions queue of the Interconnect
	for i, item := range mps.Interconnect.Logs.Items {
		// Create an LogObject
		logObj := utils.LogObject{
			Order: i,
			Log:   item,
		}
		// Append the transaction object to the list
		logs = append(logs, logObj)
	}

	// Create a struct
	ic := utils.AboutInterconnect{
		Status: mps.Interconnect.Status,
		Logs:   logs,
	}

	// Create an empty BlockObjectList
	memoryBlocks := utils.BlockObjectList{}
	for i := 0; i <= 15; i++ {
		// Create a new BlockObject instance
		blockObj := utils.BlockObject{
			Address: i,
			Data:    int(mps.MainMemory.Data[i]),
		}
		// Append the new BlockObject to the BlockObjectList
		memoryBlocks = append(memoryBlocks, blockObj)
	}

	// Create a the final JSON struct
	mm := utils.AboutMainMemory{
		Status: mps.MainMemory.Status,
		Blocks: memoryBlocks,
	}

	// Create the final object
	JSON := utils.MultiprocessingSystemState{
		PEs: pes,
		CCs: ccs,
		IC:  ic,
		MM:  mm,
	}

	// Return a string with the JSON as a string
	// Marshal the PE struct into a JSON string
	jsonData, err := json.MarshalIndent((JSON), "", "    ")
	if err != nil {
		return "", err
	}

	// Convert the byte slice to a string
	jsonString := string(jsonData)

	return jsonString, nil
}

// Function to apply a steping to an individual Processing Element
func (mps *MultiprocessingSystem) SteppingProcessingElement(ID int) string {
	if ID < -1 || ID > 2 {
		return "Invalid PE number"
	}
	pe := mps.ProcessingElements[ID]
	if !pe.IsDone && !pe.IsExecutingInstruction {
		pe.Control <- true
		return "Sent 'step' command to PE"
	} else {
		return "PE is not available..."
	}
}

// Function to apply a stepping to an individual Processing Element
func (mps *MultiprocessingSystem) StartProcessingElements() {
	allDone := false

	go func() {
		for !allDone {
			select {
			case <-mps.Terminate:
				fmt.Println("Received termination signal. Gracefully terminating.")
				return
			default:
				// Execute the code block
				for _, pe := range mps.ProcessingElements {
					if !pe.IsDone && !pe.IsExecutingInstruction {
						pe.Control <- true
						// Introduce a delay here to avoid tight loops
						time.Sleep(2 * time.Second)
					}
				}

				// Check the condition to exit the loop
				allDone = true
				for _, pe := range mps.ProcessingElements {
					if !pe.IsDone {
						allDone = false
						break
					}
				}
			}
		}
	}()
}

// Function to obtain the results after the execution of the Multiprocessing System
func (mps *MultiprocessingSystem) AboutResults() (string, error) {
	// Create an empty TransactionObjectList
	transactions := utils.TransactionObjectList{}
	// Get the values from the transactions queue of the Interconnect
	for i, item := range mps.Interconnect.Transactions.Items {
		// Create an TransactionObject
		transactionObj := utils.TransactionObject{
			Order:       i,
			Transaction: item,
		}
		// Append the transaction object to the list
		transactions = append(transactions, transactionObj)
	}
	// Sum the Cache Misses and Cache Hits for the 3 Cache Controllers
	CacheMisses := 0
	CacheHits := 0
	totalMemoryAccesses := 0
	for _, cc := range mps.CacheControllers {
		CacheMisses += cc.CacheMisses
		CacheHits += cc.CacheHits
	}
	totalMemoryAccesses = CacheHits + CacheMisses
	// Calculate the Miss Rate and Hit Rate
	MissRate := float64(CacheMisses) / float64(totalMemoryAccesses) * 100
	HitRate := float64(CacheHits) / float64(totalMemoryAccesses) * 100
	// Create the JSON object
	resultsJSON := utils.MultiprocessingSystemResults{
		Transactions:          transactions,
		PowerConsumption:      mps.Interconnect.PowerConsumption,
		CacheMisses:           CacheMisses,
		CacheHits:             CacheHits,
		MemoryAccesses:        totalMemoryAccesses,
		MissRate:              MissRate,
		HitRate:               HitRate,
		ReadRequests:          mps.Interconnect.ReadRequests,
		ReadExclusiveRequests: mps.Interconnect.ReadExclusiveRequests,
		DataResponses:         mps.Interconnect.DataResponses,
		Invalidates:           mps.Interconnect.Invalidates,
		MemoryReads:           mps.Interconnect.MemoryReads,
		MemoryWrites:          mps.Interconnect.MemoryWrites,
	}
	// Marshal the PE struct into a JSON string
	jsonData, err := json.MarshalIndent(resultsJSON, "", "    ")
	if err != nil {
		return "", err
	}
	// Convert the byte slice to a string
	jsonString := string(jsonData)
	return jsonString, nil
}

// Function to check if the Multiprocesing System has finished smoothly
func (mps *MultiprocessingSystem) AreWeFinished() bool {
	allDone := true
	for _, pe := range mps.ProcessingElements {
		if !pe.IsDone {
			allDone = false
			break
		}
	}
	return allDone
}

// Function to stop a new Multiprocessing System after initialized
func (mps *MultiprocessingSystem) Stop() {
	close(mps.Terminate)
	mps.WG.Wait() // Wait for all goroutines to finish gracefully
	// Close the log files for all PEs
	for _, pe := range mps.ProcessingElements {
		pe.Logger.Writer().(*os.File).Close()
	}
	// Close the log files for all CCs
	for _, cc := range mps.CacheControllers {
		cc.Logger.Writer().(*os.File).Close()
	}
	// Close the log file for the IC
	mps.Interconnect.Logger.Writer().(*os.File).Close()
	// Close the log file for the MM
	mps.MainMemory.Logger.Writer().(*os.File).Close()
	for i := 0; i < 3; i++ {
		close(mps.RequestChannelsM1[i])
		close(mps.ResponseChannelsM1[i])
		close(mps.RequestChannelsM2[i])
		close(mps.ResponseChannelsM2[i])
		close(mps.RequestChannelsBroadcast[i])
		close(mps.ResponseChannelsBroadcast[i])
	}
	close(mps.RequestChannelM3)
	close(mps.ResponseChannelM3)
	close(mps.Semaphore)
}

// Function to check if any file is empty
func FileIsEmpty(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		return false, err
	}

	// If the file size is 0, it's empty
	if stat.Size() == 0 {
		return true, nil
	}

	return false, nil
}
