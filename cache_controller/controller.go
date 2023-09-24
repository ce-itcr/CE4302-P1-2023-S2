package cache_controller

import (
	"math/rand"
)

type Controller struct{
	cache	*Cache
}

// Structure for read/write requests
type Request struct {
    Type string                     // WRITE or READ operation
    Address int                     // The address to READ or WRITE from
    Data int                        // (Only for WRITE) The data to store
}

type Response struct {
    Status bool                     // Status to know if the request was successfull
    Type string                     // WRITE or READ operation
    Data int                        // (Only for READ) The data to store in the register
	Address int
	Status_Data string
}

func NewController() *Controller{
	return &Controller{
		cache: NewCache(),
	}
}

func (controller *Controller) Write(pos int, data int, address int, state string){
	if(data >= 0){
		controller.cache.setData(pos, data);
	}

	if(address >= 0){
		controller.cache.setAddress(pos, address);
	}

	if(state != ""){
		controller.cache.setState(pos, state);
	}
}

func (controller *Controller) Write(pos int, data int, address int, state string){
	if(data >= 0){
		controller.cache.setData(pos, data);
	}

	if(address >= 0){
		controller.cache.setAddress(pos, address);
	}

	if(state != ""){
		controller.cache.setState(pos, state);
	}
}

//write cambiar status

func (controller *Controller) Read(address int) *Response{
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
		Request{
			Type: "READ",
			Address: address,
			Data: 0,
		}
		// Make request


		//Check cache
		return &Response{
			Status: true,
			Type: "READ",
			Data: "Response.Data",
		}
	}else{
		if(controller.cache.status[pos] == "I"){
			Request{
				Type: "READ",
				Address: controller.cache.address[pos],
				Data: 0,
			}
			// Make request
			 
		}

		return &Response{
			Status: true,
			Type: "READ",
			Data: "Response.Data",
		} 
	}
	
}

func (controller *Controller) CacheReplace (res *Response){
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
			Request{
				Type: "WRITE",
				Address: controller.cache.address[pos],
				Data: controller.cache.data[pos],
			}
			// Make request
		}
	}

	controller.cache.address[pos] = res.Address
	controller.cache.data[pos] = res.Data
	controller.cache.status[pos] = res.Status_Data
}


//Buscar dato 
//revisar estados de coerencia
//transiciones

//direct mapping
//revisar lo que hay en cada registro