package interconnect

import (
    "log"
    "os"
    "sync"
	"time"
    "MultiprocessingSystem/utils"
)

type Interconnect struct {
    ID              int
    RequestChannels []chan utils.RequestM2     // Request channels from CacheControllers
    ResponseChannels []chan utils.ResponseM2   // Response channels to CacheControllers
	RequestChannelsToCC []chan utils.RequestM3
	ResponseChannelsToCC []chan utils.ResponseM3
    Quit            chan struct{}
    Logger          *log.Logger
}

func New(requestChannels []chan utils.RequestM2, responseChannels []chan utils.ResponseM2,
	requestChannelsToCC []chan utils.RequestM3, responseChannelsToCC []chan utils.ResponseM3, quit chan struct{}) (*Interconnect, error) {
    // Create the log file for the Interconnect
    logFile, err := os.Create("logs/IC/IC.log")
    if err != nil {
        log.Fatalf("Error creating log file for Interconnect: %v", err)
    }

    // Initialize logger for the Interconnect using its respective log file
    logger := log.New(logFile, "IC" +"_", log.Ldate|log.Ltime)

    return &Interconnect{
        RequestChannels: requestChannels,
        ResponseChannels: responseChannels,
		RequestChannelsToCC: requestChannelsToCC,
		ResponseChannelsToCC: responseChannelsToCC,
        Quit:            quit,
        Logger:          logger,
    }, nil
}

func (ic *Interconnect) Run(wg *sync.WaitGroup) {
    ic.Logger.Printf(" - IC is running.\n")

	for {
		select {
		case request := <- ic.RequestChannels[0]:
			ic.Logger.Printf(" - IC is handling a request from CC0.\n")
			if(request.Type == "Write"){
				ic.Logger.Printf(" - IC is processing the write request.\n")
				for i := 0; i < 4; i++{
					if(i != 0){
						requestCC := utils.RequestM3{
							Address:   request.Address, // Your response data
							NewStatusData: "I",
						}

						ic.RequestChannelsToCC[i] <- requestCC

						response := <- ic.ResponseChannelsToCC[i]
						if(response.Status){
							ic.Logger.Printf(" - IC found address in Cache %d.\n", i)
						}
					}
				}
				
				ic.Logger.Printf(" - IC is about to send a response to CC0\n")

				response := utils.ResponseM2{
					Status: true,
				}

				ic.ResponseChannels[0] <- response

				ic.Logger.Printf(" - IC sent a response to CC0,\n")
			}

			if(request.Type == "Read"){
				goToMem := true
				value := 0
				status := "E"
				ic.Logger.Printf(" - IC is processing the read request.\n")
				for i := 0; i < 4; i++{
					if(i != 0){
						requestCC := utils.RequestM3{
							Address:   request.Address, // Your response data
							NewStatusData: "S",
						}


						ic.RequestChannelsToCC[i] <- requestCC

						response := <- ic.ResponseChannelsToCC[i]

						if(response.Status){
							goToMem = false
							value = response.Data
							ic.Logger.Printf(" - IC found address in Cache %d.\n", i)
							status = "S"
						}
					}
				}

				// Ir a mem

				ic.Logger.Printf(" - IC is about to send a response to CC0\n")

				response := utils.ResponseM2{
					Status: true,
					Data: value,
					StatusData: status,
				}

				ic.ResponseChannels[0] <- response

				ic.Logger.Printf(" - IC sent a response to CC0,\n")

			}

		case request := <- ic.RequestChannels[1]:
			ic.Logger.Printf(" - IC is handling a request from CC1.\n")
			
			ic.Logger.Printf(" - IC is processing the request.\n")
			time.Sleep(2 * time.Second)
			// Create a response for the CacheController
			response := utils.ResponseM2{
				Status: true,
				Type:   request.Type,
				Data:   42, // Your response data
			}

			ic.Logger.Printf(" - IC is about to send a response to CC1\n")
			// Send the response back to the CacheController 1
			ic.ResponseChannels[1] <- response
			ic.Logger.Printf(" - IC sent a response to CC1,\n")

		case request := <- ic.RequestChannels[2]:l
			ic.Logger.Printf(" - IC is handling a request from CC2.\n")

			ic.Logger.Printf(" - IC is processing the request.\n")
			time.Sleep(2 * time.Second)
			// Create a response for the CacheController
			response := utils.ResponseM2{
				Status: true,
				Type:   request.Type,
				Data:   42, // Your response data
			}

			ic.Logger.Printf(" - IC is about to send a response to CC2.\n")
			// Send the response back to the CacheController 2
			ic.ResponseChannels[2] <- response
			ic.Logger.Printf(" - IC sent a response to CC2,\n")

		case <- ic.Quit:
			ic.Logger.Printf(" - IC has received an external signal to terminate.\n")
			return
		}
	}
}
