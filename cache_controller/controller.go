package cache_controller

import (
	"math/rand"
	"sync"
	"cache_controller.com/utils"
)

type Controller struct{
	cache	*Cache
	RequestChannel chan utils.RequestM1
	ResponseChannel chan utils.ResponseM1
	RequestMemChannel chan utils.RequestM2
	ResponseMemChannel chan utils.ResponseM2
	RequestMemChannelToCC chan utils.RequestM3
	ResponseMemChannelToCC chan utils.ResponseM3
		Quit chan struct{}
}


func NewController(rq1 utils.RequestM1, rs1 utils.ResponseM1,rq2 utils.RequestM2,
				rs2 utils.ResponseM2, rq3 utils.RequestM3, rs3 utils.ResponseM3	) *Controller{
	return &Controller{
		cache: NewCache(),
	}
}

func (cc *Controller) Run(wg *sync.WaitGroup){
	for {
		select{
		case request := <- cc.RequestChannel:
			
			if(request.Type == "Write"){
				cc.Write(request.Address, request.Data) // Función Write
				response := utils.ResponseM1{
					Status: true,
					Type: request.Type,
				}
	
				cc.ResponseChannel <- response
			}

			if(request.Type == "Read"){
				Mem := cc.Read(request.Address)
				response := utils.ResponseM1{
					Status: true,
					Type: request.Type,
					Data: cc.cache.Data[Mem],
				}
	
				cc.ResponseChannel <- response
			}
		
		case request := <- cc.RequestMemChannelToCC:

			temp := cc.GetAddressStatus(&request)
			
			response := utils.ResponseM3{
					Status: temp.Status,
					Data: temp.Data,
					GoToMem: temp.GoToMem,
				}

			cc.ResponseMemChannelToCC <- response 

		case <- cc.Quit:
		}
	}
}

func (controller *Controller) Write(address int, data int){

	for pos, addr := range controller.cache.Address{	//Busca address en cache y escribe
		if(address == addr){
			controller.cache.Address[pos] = address
			controller.cache.Data[pos] = data
			controller.cache.Status[pos] = "M"
			break
		}
	}

}

func (controller *Controller) Read(address int) int{
	pos := 0
	found := false
	for _, addr := range controller.cache.Address{	//Busca address en Cache
		if(address == addr){
			found = true
			break
		}
		pos++
	}

	if(!found){	// No está en cache
		requestMem := utils.RequestM2{
			Type: "READ",
			Address: address,
		}

		controller.RequestMemChannel <- requestMem		//Pide a memoria

		response := <- controller.ResponseMemChannel
		

		pos = controller.CacheReplace(&response, address)	//Busca espacio en cache para el dato

		return pos

	}else{
		if(controller.cache.Status[pos] == "I"){		//Sí está, revisa si es invalido
			requestMem := utils.RequestM2{
				Type: "READ",
				Address: controller.cache.Address[pos],
			}
			controller.RequestMemChannel <- requestMem

			response := <- controller.ResponseMemChannel
		

			controller.cache.Data[pos] = response.Data			//Escribe dato actualizado
			controller.cache.Status[pos] = response.StatusData
			 
		}

		return pos
	}
	
}

func (controller *Controller) CacheReplace(res *utils.ResponseM2, address int) int{
	currentState := res.StatusData
	pos := 0
	for _, status := range controller.cache.Status{		//Busca vacio o Invalido para cambiar
		if(status == "I" || status == ""){
			break
		}
		pos++
	}

	if(pos == 5){		//Si no encuentra lo hace random
		pos = rand.Intn(4)
		if(controller.cache.Status[pos] == "M"){	//Si es modified lo cambia en memoria
			requestMem := utils.RequestM2{
				Type: "WRITE",
				Address: controller.cache.Address[pos],
				Data: controller.cache.Data[pos],
			}
			controller.RequestMemChannel <- requestMem

			response := <- controller.ResponseMemChannel
		}
	}

	controller.cache.Address[pos] = address		//Escribe datos nuevos
	controller.cache.Data[pos] = res.Data
	controller.cache.Status[pos] = currentState

	return pos
}

func (controller *Controller) GetAddressStatus(req *utils.RequestM3) *utils.ResponseM3{ //Revisar que hay en cache y asignar state nuevo
	goToMem := false
	respStatus := false
	for pos, addr := range controller.cache.Address{
		if(req.Address == addr && controller.cache.Status[pos] != "I"){
			if(controller.cache.Status[pos] == "M"){
				goToMem = true
			}
			respStatus := true
			controller.cache.Status[pos] = req.NewStatusData
			
			return &utils.ResponseM3{
				Status: respStatus,
				Data:   controller.cache.Data[pos],
				GoToMem: goToMem,
			}
		}
	}

	return &utils.ResponseM3{
		Status: respStatus,
		GoToMem: goToMem,
	}
}
