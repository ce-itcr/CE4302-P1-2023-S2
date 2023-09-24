package cache_controller

import (
	"strconv"
	"sync"
)

type Controller struct{
	cache	*Cache

}

func NewController() *Controller{
	return &Controller{
		cache: NewCache(),
	}
}


func Read(pos int, controller *Controller) [3]string{
	res:= [3]string{"", "", ""}
	res[0] = strconv.Itoa(controller.cache.getData(pos));
	res[1] = strconv.Itoa(controller.cache.getAddress(pos));
	res[2] = controller.cache.getState(pos);
	return res
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