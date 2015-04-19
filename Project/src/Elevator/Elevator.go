
package Elevator

import(
	"UDP"
	"driver"
	"Queue"	
	"Source"
   	//"fmt"
)

const (
	UP = 0
	DOWN = 1
)

func Elevator(){
	elevatorStatus := Source.Elevator{0, -1, -1}
	//elev2 := Source.Elevator{1, -1, -1}
	//elev3 := Source.Elevator{2, -1, -1}
	
	//STATES
	//EventHandler
	
	
	
	updateElevatorInfoChannel := make(chan Source.Elevator,1)// FÅR INN OPPDATERING FRÅ NETTVERKET*/

	//ELEV
	wait := make(chan int, 1)
	run := make(chan int, 1)
	stop := make(chan int, 1)
	orderInEmptyQueue := make(chan int, 1)

	//driver
	newOrderChannel := make(chan Source.ButtonMessage,1)
	floorReachedChannel := make(chan int, 1)
	setSpeedChannel := make(chan int, 1)
	stopChannel := make(chan int, 1)
	stoppedChannel := make(chan int, 1)
	setButtonLightChannel := make(chan Source.ButtonMessage,1) 

	//Queue
	addOrderChannel := make(chan Source.ButtonMessage,1)
	removeOrderChannel := make(chan int, 1)
	nextOrderChannel := make(chan int, 1)
	checkOrdersChannel := make(chan int, 1)
	orderRemovedChannel := make(chan int, 1)
	newElevInfoChannel := make(chan Source.Elevator,1)
	fromUdpToQueue := make(chan Source.Message, 1)
	//UDP

	completedOrderChannel := make(chan Source.Message, 1)
	externalOrderChannel := make(chan Source.ButtonMessage, 1)
	UDPaddOrderChannel := make(chan Source.Message, 1)
	removeExternalOrderChannel := make(chan Source.Message, 1)

	nextOrderedFloor := 100


	
	
	go driver.Drivers(newOrderChannel, floorReachedChannel, setSpeedChannel, stopChannel, stoppedChannel, setButtonLightChannel)
  	go Queue.Queue(addOrderChannel, removeOrderChannel, nextOrderChannel, checkOrdersChannel, orderInEmptyQueue, orderRemovedChannel, newElevInfoChannel, fromUdpToQueue)
   	go handleOrders(2, addOrderChannel , setButtonLightChannel, newOrderChannel, externalOrderChannel, UDPaddOrderChannel)
	//SKAL FROMUDPTOQUEUE tAKAST INN I SLAVE???????????????????????????
	go UDP.Slave(completedOrderChannel, externalOrderChannel, updateElevatorInfoChannel, UDPaddOrderChannel , removeExternalOrderChannel, newElevInfoChannel, fromUdpToQueue)

	run <- 1



	for{
		select{
			case arrivedAtFloor := <- floorReachedChannel:// FLOOR REACHED		
	
				elevatorStatus.CurrentFloor = arrivedAtFloor
				
				/*elev1.CurrentFloor = currentFloor
				elev1.Direction = movingDirection
				elevatorInfoChannel <- elev1*/

				checkOrdersChannel <- elevatorStatus.CurrentFloor
				nextOrder := <- nextOrderChannel
				nextOrderedFloor = nextOrder
				direction := nextOrderedFloor - elevatorStatus.CurrentFloor

				if(direction > 0){
					elevatorStatus.Direction = UP
				}else{
					elevatorStatus.Direction = DOWN
				}

				updateElevatorInfoChannel <- elevatorStatus
				newElevInfoChannel <- elevatorStatus
				//println("currentFloor:", elevatorStatus.CurrentFloor, "\nOrderedFloor", nextOrderedFloor, "Direccion:", elevatorStatus.Direction)
				//println("NEXTORDEEER", nextOrderedFloor, "Current floor " , currentFloor)
				if(elevatorStatus.CurrentFloor == nextOrderedFloor ){
					stop <- elevatorStatus.Direction
				}
												 	
			case <- stop:

               /// println("ELEVATOR: StopChannel")
				stopChannel <- elevatorStatus.Direction
				wait <- 1
                break

			case <- wait:
				
                //println("ELEVATOR :WAIT")
			 	wait:
				for{
					select{
						case <- stoppedChannel:
							removeOrderChannel <- elevatorStatus.CurrentFloor
							<- orderRemovedChannel	
							//println("order removed")
							run <- 1
							break wait
					}
				}
							
			case <- run:
        		//println("ELEVATOR :RUN")
				checkOrdersChannel <- elevatorStatus.CurrentFloor
				orderedFloor := <- nextOrderChannel
				nextOrderedFloor = orderedFloor
				//println("ELEVATOR :RUN ORDER: ", nextOrderedFloor)
				if(nextOrderedFloor == -1){
					break
				}else{
					if(nextOrderedFloor > elevatorStatus.CurrentFloor){
					 	setSpeedChannel <- 0
					}else if(nextOrderedFloor < elevatorStatus.CurrentFloor){
						setSpeedChannel <- 1
					}else{
						stop <- elevatorStatus.Direction
					}
				}

            case <- orderInEmptyQueue:
				run <- 1

/*
		  	case updatedElevInfo := <- UpdateElevatorInfoChannel:

				if(elev1.ID == updatedElevInfo.ID){
					elev1 = updatedElevInfo
				} else if (elev2.ID == updatedElevInfo.ID){
					elev2 = updatedElevInfo
				} else if (elev3.ID == updatedElevInfo.ID){
					elev3 = updatedElevInfo
				}
				*/
/*			case newReceivedOrder := <- addOrderChannel:

				if(newReceivedOrder.MessageTo == elev1.ID){
					Queue.AddOrder( newReceivedOrder.Button, newReceivedOrder.ID, currentFloor, movingDirection)
  			  	} 
				driver.Elev_set_button_lamp(newReceivedOrder.Button) 
              	
		*/	
		}
	}
}


func handleOrders(elevatorID int, addOrderChannel chan Source.ButtonMessage, setButtonLightChannel chan Source.ButtonMessage, newOrderChannel chan Source.ButtonMessage, externalOrderChannel chan Source.ButtonMessage, UDPaddOrderChannel chan Source.Message, fromUdpToQueue chan Source.Message){
	for{
		select{
			case newOrder := <- newOrderChannel:
				newOrder.Value = 1
				if(newOrder.Button == Source.BUTTON_COMMAND){
					addOrderChannel <- newOrder
					setButtonLightChannel <- newOrder
				} else{
					externalOrderChannel <- newOrder
					setButtonLightChannel <- newOrder
				}
			case newExternalOrder := <- UDPaddOrderChannel:	
				fromUdpToQueue <- newExternalOrder  
				if(newExternalOrder.CompletedOrder && elevatorID != newExternalOrder.MessageTo){
					newExternalOrder.Button.Value = 0
					setButtonLightChannel <- newExternalOrder.Button
				} else if(newExternalOrder.NewOrder){
				 	newExternalOrder.Button.Value = 1
					setButtonLightChannel <- newExternalOrder.Button
				}
	
		}
	}
}



