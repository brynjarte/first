package Queue

import(
	"FileHandler"
	"driver" // MÅ FINNA EI LØYSING, TRENGER KUN BUTTONMESSAGE
)

var NumOfElevs int
var NumOfFloors int
const UP int = 0
const DOWN int = 1

var allQueues[5][] FileHandler.Directions
var queue[] FileHandler.Directions
var numberInQueue[5] int
var ordersInDirection[2] int 


func queueInit(){
	
	queueList := FileHandler.Read(&NumOfElevs, &NumOfFloors)
	
	for j:=0;j<len(queueList);j+=2{
		q := FileHandler.Directions{queueList[j], queueList[j+1]}
		queue = append(queue, q)
		//internalOrders = append(internalOrders,q)
	}
	//AddQueue(queue, ourID)
	ordersInDirection[UP] = 0
	ordersInDirection[DOWN] = 0
	for i,_ := range queue {
		ordersInDirection[UP] += queue[i].UP
		ordersInDirection[DOWN] = queue[i].DOWN
	}
	
	/*
	for i:=0;i<NumOfElevs;i++ {
		numInQueue = append(numInQueue,0}
	}
	*/
	addQueue(1, queue)
	addQueue(0, queue)
}

func addQueue(elevatorID int, queue []FileHandler.Directions) {
	allQueues[elevatorID] = queue
}



func Queue(addOrderChannel chan driver.ButtonMessage, removeOrderChannel chan int, setDirectionChannel chan int, checkDirectionChannel chan int, checkOrdersChannel chan int, stopChannel chan int){//, findBestElevator chan ButtonMessage ){
	newDirection := -1
	direction := 0
	currentFloor := 0
	queueInit()
	for{
		println("Current Floor: ", currentFloor, "\nDirection: ", direction)
		println("UPS\tDOWNS")
		for i:= 0; i < NumOfFloors; i++{
			println(allQueues[1][i].UP, "\t", allQueues[1][i].DOWN)		
		}
		println()
		println()
		select{
			case newOrder := <- addOrderChannel:
				addOrder(1, newOrder, currentFloor, direction)

			case floor := <- removeOrderChannel:
				removeOrder(1, floor, direction)
				
			case movingDirection := <- checkDirectionChannel:
				direction = movingDirection
				if(checkIfOrdersInDirection(1, currentFloor, direction) || newDirection != -1){ // ELEVATORID = 1
					setDirectionChannel <- direction
				} else{
					newDirection := setDirection()
					if(newDirection != -1){	
						direction = newDirection
						setDirectionChannel <- newDirection
					}else{
						checkDirectionChannel <- 1
					}			
				}
			case floor := <- checkOrdersChannel:
				currentFloor = floor
				//fmt.Println(currentFloor,direction)
				if(checkOrders(1, currentFloor, direction)){
					stopChannel <- 1
				}
			//case <- findBestElevatorChannel:

		}
	}
}

func setDirection() int{
    if(ordersInDirection[0] != 0){
        return 0
    }else if(ordersInDirection[1] != 0){
        return 1
    }else{
        return -1
    }
}

func checkIfOrdersInDirection(elevatorID int, currentFloor int, direction int ) bool{
	if(numberInQueue [elevatorID] == 0 || direction == -1){
		return false
	}
	if(direction == UP){
		if(currentFloor == NumOfFloors-1){
			return false
		}
		for floor := currentFloor ;floor < NumOfFloors; floor++{
			if(allQueues[elevatorID][floor].UP == 1){
				return true
			}
		}
	}else if (direction == DOWN){
		if(currentFloor == 0){
			return false
		}
		for floor := currentFloor ;floor > 0; floor--{
			if(allQueues[elevatorID][floor].DOWN == 1){
				return true
			}
		}
	}
	return false
}

func checkOrders(elevatorID int, currentFloor int, direction int) bool{
	if (direction == UP) {	
		if(allQueues[elevatorID][currentFloor].UP == 1){
			return true
		}
		return false
	} else if (direction == DOWN) {
		if(allQueues[elevatorID][currentFloor].DOWN == 1){
			return true
		}
		return false
	}
	return false
}



func addOrder(elevatorID int, order driver.ButtonMessage, currentFloor int, movingDirection int) {
	//Antar at vi vet hvilken heis vi skal bruk
	var directionOfOrder = -1
	
	if ( currentFloor - order.Floor > 0 ) {
        directionOfOrder = DOWN
    } else if ( currentFloor - order.Floor < 0 ) {
        directionOfOrder = UP
    } else {
		directionOfOrder = movingDirection
	}
	
    if( order.Button == driver.BUTTON_CALL_UP && allQueues[elevatorID][order.Floor].UP==0 ) {
        allQueues[elevatorID][order.Floor].UP = 1
        numberInQueue[elevatorID] ++
        ordersInDirection[UP]++
    } else if( order.Button == driver.BUTTON_CALL_DOWN && allQueues[elevatorID][order.Floor].DOWN==0 ) {
        allQueues[elevatorID][order.Floor].DOWN = 1
        numberInQueue[elevatorID] ++
        ordersInDirection[DOWN] ++
    } else if(order.Button == driver.BUTTON_COMMAND) {
    	
    	if directionOfOrder == UP {
			if( allQueues[elevatorID][order.Floor].UP==0){
            	numberInQueue[elevatorID] ++
            	ordersInDirection[directionOfOrder] ++
            	allQueues[elevatorID][order.Floor].UP = 1
			}	
		} else if directionOfOrder == DOWN {
			if( allQueues[elevatorID][order.Floor].DOWN == 0){
            	numberInQueue[elevatorID] ++
            	ordersInDirection[directionOfOrder] ++
            	allQueues[elevatorID][order.Floor].DOWN = 1
			}
		}	
   	}
	
}
	
func removeOrder(elevatorID int, floor int, movingDirection int) {
	
	if (movingDirection == 0) {
		allQueues[elevatorID][floor].UP = 0
	}else {
		allQueues[elevatorID][floor].DOWN = 0 
	}
}
	

func findBestElevator(elevatorID int, order driver.ButtonMessage) int{
	min := 0
	for elev:= 0; elev < NumOfElevs; elev++{
			if(numberInQueue [elev] <= min){
				min = elev
			}
	}
	return min
}


/*func ClearAllExternalOrders(elevatorID int) {
	numberInQueue [elevatorID] = 0
	for floor := 0; floor < NumOfFloors; floor++ {
		if(internalOrders[floor][UP]){
			numberInQueue [elevatorID]++
		}
		if(internalOrders[floor][DOWN]){
			numberInQueue [elevatorID]++
		}
	}
	queue = internalOrders
}
*/


