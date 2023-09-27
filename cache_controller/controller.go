package cache_controller

import (
	"math/rand"
	"sync"
	"cache_controller.com/utils"
	"time"
)

type Controller struct{
	cache	*Cache
	RequestChannel chan utils.Request
	ResponseChannel chan utils.Response
	RequestMemChannel chan utils.RequestMem
	ResponseMemChannel chan utils.ResponseMem
		Quit chan struct{}
}


func NewController() *Controller{
	return &Controller{
		cache: NewCache(),
	}
}

func (cc *Controller) Run(wg *sync.WaitGroup){
	for {
		select{
		case request := <- cc.RequestChannel:
			time.Sleep(10 * time.Second)

			response := utils.Response{
				Status: true,
				Type: request.Type,
				Data: 12,
			}

			cc.ResponseChannel <- response

		case <- cc.Quit:
		}
	}
}

func (controller *Controller) Write(req *utils.RequestMem){
	utils.RequestMem{
		Type: "Write",
		Address: req.Address,
		Data: req.Data,
	}
	// Hacer algo en Memoria
}

//write cambiar status

func (controller *Controller) Read(address int) *utils.ResponseMem{
	pos := 0
	found := false
	for _, addr := range controller.cache.address{
		if(address == addr){
			found = true
			break
		}
		pos++
	}

	if(found == false){
		utils.RequestMem{
			Type: "READ",
			Address: address,
		}
		// Make request


		//Check cache
		return &utils.ResponseMem{
			Status: true,
			Address: 3,
			Data: 0,
			StatusData: "E",
		}
	}else{
		if(controller.cache.status[pos] == "I"){
			utils.RequestMem{
				Type: "READ",
				Address: controller.cache.address[pos],
			}
			// Make request
			 
		}

		return &utils.ResponseMem{
			Status: true,
			Address: 3,
			Data: 0,
			StatusData: "E",
		}
	}
	
}

func (controller *Controller) CacheReplace (res *utils.ResponseMem){
	pos := 0
	for _, status := range controller.cache.status{
		if(status == "I" || status == ""){
			break
		}
		pos++
	}

	if(pos == 5){
		pos = rand.Intn(4)
		if(controller.cache.status[pos] == "M"){
			utils.RequestMem{
				Type: "WRITE",
				Address: controller.cache.address[pos],
				Data: controller.cache.data[pos],
			}
			// Make request
		}
	}

	controller.cache.address[pos] = res.Address
	controller.cache.data[pos] = res.Data
	controller.cache.status[pos] = res.StatusData
}

func (controller *Controller) GetAddressStatus(address int) string{
	pos := 0
	for _, addr := range controller.cache.address{
		if(address == addr){
			break
		}
		pos++
	}

	if(pos == 5){
		return ""
	}

	return controller.cache.status[pos]

}


//Buscar dato 
//revisar estados de coerencia
//transiciones

//direct mapping
//revisar lo que hay en cada registro