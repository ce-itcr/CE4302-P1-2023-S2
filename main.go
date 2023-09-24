package main


import (
	"fmt"
	"cache_controller.com/cache_controller"
)

func main(){
	var cc1 *cache_controller.Controller = cache_controller.NewController()
	var cc2 *cache_controller.Controller = cache_controller.NewController()
	var cc3 *cache_controller.Controller = cache_controller.NewController()


	cc1.Write(0, 2, 3, "E")
	fmt.Println(cc1.Read(0)[0])
	fmt.Println(cc1.Read(0)[1])
	cc2.Write(0, 5, 5, "S")
	fmt.Println(cc1.Read(0)[2])
	fmt.Println(cc2.Read(0)[0])
	fmt.Println(cc2.Read(0)[1])
	cc3.Write(1, 5, 5, "S")
	fmt.Println(cc2.Read(0)[2])
	fmt.Println(cc3.Read(1)[0])
	fmt.Println(cc3.Read(1)[1])
	fmt.Println(cc3.Read(1)[2])



}