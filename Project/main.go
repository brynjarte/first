package main  

import (
	"driver"
)

func main() {
	
	someChannel := make(chan int,1)
	driver.Elev_init()
	go driver.ReadButtons()
	
	go func () {
		for {
		read := <- driver.ReadButtonsChannel
		println(read.Button)
		println(read.Floor)
		}	
	}()
	<- someChannel
	println()
}

