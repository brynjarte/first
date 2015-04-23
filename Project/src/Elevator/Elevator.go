
package Elevator

import(
	"Network"
	"driver"
	"Queue"	
	"Source"

)

func errorHandling(elevatorID int){
	defer Elevator(elevatorID)
	for{
		select{
			case err := <- Source.ErrorChannel:
				if( err != nil){
					//SKRIV TIL LOGG
					panic(err)
					return
				}
		}
	}
}

func Elevator(elevatorID int){

	elevatorStatus := Source.ElevatorInfo{elevatorID, -1, -1}

	//STATES
	//EventHandler
	updateElevatorInfoChannel := make(chan Source.ElevatorInfo,1)

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
	initFinished := make(chan int)
	//Queue
	addOrderChannel := make(chan Source.ButtonMessage,1)
	removeOrderChannel := make(chan int, 1)
	nextOrderChannel := make(chan int, 1)
	checkOrdersChannel := make(chan int, 1)
	finishedRemoving := make(chan int, 1)
	fromElevToQueue := make(chan Source.Message, 1)
	orderRemovedChannel := make(chan Source.ButtonMessage, 1)
	
	requestQueueChannel := make(chan int, 1)
	receiveQueueChannel := make(chan Source.Message, 10)
	//UDP

	completedOrderChannel := make(chan Source.ButtonMessage, 1)
	externalOrderChannel := make(chan Source.ButtonMessage, 1)
	handleOrderChannel := make(chan Source.Message, 1)
	removeExternalOrderChannel := make(chan Source.Message, 1)
	bestElevatorChannel := make(chan Source.Message, 1)
	removeElevatorChannel := make(chan int, 1)
	nextOrderedFloor := 100

	go errorHandling(elevatorStatus.ID)
	go Source.SourceInit()
	go driver.Drivers(newOrderChannel, floorReachedChannel, setSpeedChannel, stopChannel, stoppedChannel, setButtonLightChannel, initFinished)
	<- initFinished
	println("Finished inting")
  	go Queue.Queue(elevatorStatus, addOrderChannel, removeOrderChannel, nextOrderChannel, checkOrdersChannel, orderInEmptyQueue, finishedRemoving, fromElevToQueue, bestElevatorChannel, removeElevatorChannel, completedOrderChannel, orderRemovedChannel, requestQueueChannel, receiveQueueChannel)
   	go handleOrders(elevatorStatus.ID, addOrderChannel , setButtonLightChannel, newOrderChannel, externalOrderChannel, handleOrderChannel, fromElevToQueue, orderRemovedChannel, completedOrderChannel)
	go Network.Slave(elevatorStatus, externalOrderChannel, updateElevatorInfoChannel, handleOrderChannel, removeExternalOrderChannel, bestElevatorChannel, removeElevatorChannel, completedOrderChannel, requestQueueChannel, receiveQueueChannel)
	

	//run <- 1
	prevFloor := 10

	
	for{
		select{
			case arrivedAtFloor := <- floorReachedChannel:// FLOOR REACHED		
				println("FLOORREACHED")
				prevFloor = elevatorStatus.CurrentFloor
				elevatorStatus.CurrentFloor = arrivedAtFloor

				checkOrdersChannel <- elevatorStatus.CurrentFloor
				nextOrder := <- nextOrderChannel
				nextOrderedFloor = nextOrder
				direction := prevFloor - elevatorStatus.CurrentFloor

				if(direction < 0){
					elevatorStatus.Direction = Source.UP
				}else if (direction > 0) {
					elevatorStatus.Direction = Source.DOWN
				}
				
				if(elevatorStatus.CurrentFloor < nextOrderedFloor && elevatorStatus.Direction == Source.DOWN){
					run <- 1
				} else if(elevatorStatus.CurrentFloor > nextOrderedFloor && elevatorStatus.Direction == Source.UP){
					run <- 1
				}
					
				updateElevatorInfoChannel <- elevatorStatus
				//println("currentFloor:", elevatorStatus.CurrentFloor, "\nOrderedFloor", nextOrderedFloor, "Direccion:", elevatorStatus.Direction)
				//println("NEXTORDEEER", nextOrderedFloor, "Current floor " , currentFloor)
				if(elevatorStatus.CurrentFloor == nextOrderedFloor ){
					stop <- elevatorStatus.Direction
				}
												 	
			case <- stop:
                println("ELEVATOR: StopChannel")
				stopChannel <- elevatorStatus.Direction
				wait <- 1
                break

			case <- wait:
				
                println("ELEVATOR :WAIT")
			 	wait:
				for{
					select{
						case <- stoppedChannel:
							removeOrderChannel <- elevatorStatus.CurrentFloor
							println("Removing")
							<- finishedRemoving
							println("order removed")
							run <- 1
							break wait
					}
				}
							
			case <- run:
        		println("ELEVATOR :RUN")
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

		}
	}
}


func handleOrders(elevatorID int, addOrderChannel chan Source.ButtonMessage, setButtonLightChannel chan Source.ButtonMessage, newOrderChannel chan Source.ButtonMessage, externalOrderChannel chan Source.ButtonMessage, handleOrderChannel chan Source.Message, fromElevToQueue chan Source.Message, orderRemovedChannel chan Source.ButtonMessage, completedOrderChannel chan Source.ButtonMessage){
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
			case newExternalOrder := <- handleOrderChannel:	
				println("New ext update is", newExternalOrder.UpdatedElevInfo)
				fromElevToQueue <- newExternalOrder  
				if(newExternalOrder.CompletedOrder && elevatorID != newExternalOrder.MessageTo){
					newExternalOrder.Button.Value = 0
					setButtonLightChannel <- newExternalOrder.Button
				} else if(newExternalOrder.AcceptedOrder){
				 	newExternalOrder.Button.Value = 1
					setButtonLightChannel <- newExternalOrder.Button
				}
			case orderRemoved := <- orderRemovedChannel:
				orderRemoved.Value = 0
				if (orderRemoved.Button != Source.BUTTON_COMMAND) {
					completedOrderChannel <- orderRemoved
					fromElevToQueue <- Source.Message{false, false, false, true, false, elevatorID, -1, Source.ElevatorInfo{elevatorID, -1, -1}, orderRemoved}
				}
				setButtonLightChannel <- orderRemoved
		}
	}
}



