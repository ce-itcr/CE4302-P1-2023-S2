package main


import (
	"fmt"
	"cache_controller.com/cache_controller"
	"cache_controller.com/utils"
)

func main(){
	RequestChannel := make(chan utils.RequestM1)
	ResponseChannel := make(chan utils.ResponseM1)
	RequestMemChannel := make(chan utils.RequestM2)
	ResponseMemChannel := make(chan utils.ResponseM2)
	RequestMemChannelToCC := make(chan utils.RequestM3)
	ResponseMemChannelToCC := make(chan utils.ResponseM3)

	var cc1 *cache_controller.Controller = cache_controller.NewController(
	RequestChannel,
	ResponseChannel,
	RequestMemChannel,
	ResponseMemChannel,
	RequestMemChannelToCC,
	ResponseMemChannelToCC,
	)


	cc1.Write(0, 2, 3, "E")




}