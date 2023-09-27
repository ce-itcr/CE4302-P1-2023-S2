// cache_controller/cache_controller.go

package CacheController

import (
	"time"
	"sync"
	"log"
	"os"
	"strconv"
	"MultiprocessingSystem/utils"
)

type CacheController struct {
	ID int
    RequestChannelM1 chan utils.RequestM1 
    ResponseChannelM1 chan utils.ResponseM1
	RequestChannelM2 chan utils.RequestM2 
    ResponseChannelM2 chan utils.ResponseM2  
	Quit chan struct{}
	Logger *log.Logger
}


func New(id int, RequestChannelPE chan utils.RequestM1, 
					ResponseChannelPE chan utils.ResponseM1,
					RequestChannelIC chan utils.RequestM2,
					ResponseChannelIC chan utils.ResponseM2,
					quit chan struct{}) (*CacheController, error) {

    // Create the log file
    logFile, err := os.Create("logs/CC/CC" + strconv.Itoa(id) + ".log")
    if err != nil {
        log.Fatalf("Error creating log file for CC%d: %v", id, err)
    }

    // Initialize logger for the PE using its respective log file
    logger1 := log.New(logFile, "CC" + strconv.Itoa(id) + "_", log.Ldate|log.Ltime)
	
	return &CacheController{
		ID: id,
		Logger: logger1,
        RequestChannelM1: RequestChannelPE,
        ResponseChannelM1: ResponseChannelPE,
		RequestChannelM2: RequestChannelIC,
		ResponseChannelM2: ResponseChannelIC,
		Quit: quit,
    }, nil
}

func (cc *CacheController) Run(wg *sync.WaitGroup) {
	cc.Logger.Printf(" - CC%d is running.\n", cc.ID)
    for {
        // Listen to the PE for a request
		select {
			case request := <- cc.RequestChannelM1:
				cc.Logger.Printf(" - CC%d is Received request from PE%d.\n", cc.ID, cc.ID)
	
				switch request.Type {
				case "READ":
					cc.Logger.Printf(" - CC%d is processing a READ request.\n", cc.ID)
					cc.Logger.Printf(" - Address: %d.\n", request.Address)
					time.Sleep(3 * time.Second)
				
				case "WRITE":	
					cc.Logger.Printf(" - CC%d is processing a WRITE request.\n", cc.ID)
					cc.Logger.Printf(" - Address: %d, Data: %d.\n", request.Address, request.Data)
					time.Sleep(5 * time.Second)

					// Send a request to the Interconnect
					// Prepare the request message
					icRequest := utils.RequestM2 {
						Type: request.Type,
						Address: request.Address,                    
						Data: 0,                      
					}

					// Send the request to the Interconnect
					cc.Logger.Printf(" - CC%d is about to send a request to the Interconnect.\n", cc.ID)
					cc.RequestChannelM2 <- icRequest
					cc.Logger.Printf(" - CC%d sent (Type: %s, Address: %d) to the Interconnect.\n", cc.ID, icRequest.Type, icRequest.Address)

					// Wait for the response from the CacheController
					cc.Logger.Printf(" - CC%d is waiting for a response from the Interconnect....\n", cc.ID)
					icResponse := <- cc.ResponseChannelM2

					// Process the response
					cc.Logger.Printf(" - CC%d received --> Status: %v, Type: %s, Data: %d.\n", cc.ID, icResponse.Status, icResponse.Type, icResponse.Data)
				}
				
				// Create a struct to pack the response
				cc.Logger.Printf(" - CC%d is making a response object.\n", cc.ID)
				response:= utils.ResponseM1 {
					Status: true, 
					Type:   request.Type,
					Data:   12,     
				}
				
				// Enviar respuesta al PE correspondiente
				cc.Logger.Printf(" - CC%d is about to send the response to PE%d.\n", cc.ID, cc.ID)
				cc.ResponseChannelM1 <- response
				cc.Logger.Printf(" - CC%d has sent the response back to PE%d.\n", cc.ID, cc.ID)

			case request := <- cc.RequestChannelM2:
				cc.Logger.Printf(" - CC%d is Received request from Interconnect.\n", cc.ID)

				// Process the request from Interconnect
				cc.Logger.Printf(" - The request is %s.\n", request.Type)

			case <- cc.Quit:
				cc.Logger.Printf(" - CC%d has received an external signal to terminate.\n", cc.ID)
                return
		}
      
    }
}
