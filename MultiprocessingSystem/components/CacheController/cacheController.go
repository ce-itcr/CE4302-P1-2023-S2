// cache_controller/cache_controller.go

package CacheController

import (
	"time"
	"fmt"
	"sync"
	"MultiprocessingSystem/utils"
)

type CacheController struct {
    RequestChannel chan utils.Request
    ResponseChannel chan utils.Response
	Quit chan struct{} 

}


func New(RequestChannelPE chan utils.Request, ResponseChannelPE chan utils.Response, quit chan struct{}) (*CacheController, error) {
    return &CacheController{
        RequestChannel: RequestChannelPE,
        ResponseChannel: ResponseChannelPE,
		Quit: quit,
    }, nil
}

func (cc *CacheController) Run(wg *sync.WaitGroup) {
    for {
        // Listen to the PE for a request
		select {
			case request := <- cc.RequestChannel:
				// Simulate that the Cache Controller is processing the data
				time.Sleep(10 * time.Second)

				// Create a struct to pack the response
				response:= utils.Response{
					Status: true, 
					Type:   request.Type,
					Data:   12,     
				}
				
				// Enviar respuesta al PE correspondiente
				cc.ResponseChannel <- response

			case <- cc.Quit:
				fmt.Printf("Cache Received termination signal and is exiting gracefully.\n")
                return
		}
      
    }
}
