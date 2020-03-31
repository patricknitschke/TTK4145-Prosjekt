package elevator

import (
	"fmt"
	"strings"

	"../config"
)

/* Constants and variables */

// InternalQueue maintains orders for its node
var internalQueue [config.NFloors][config.NButtonTypes]bool

/* HallUp   HallDn    Cab
-|-------||-------||-------|
3| _____ || _____ || _____ |
2| _____ || _____ || _____ |
1| _____ || _____ || _____ |
0| _____ || _____ || _____ |
---------------------------- */

/* InternalQ module functions */

/* Unexported functions */

// Set up internalQueue
func internalQInit() {
	for floor := 0; floor < config.NFloors; floor++ {
		for button := 0; button < config.NButtonTypes; button++ {
			internalQPop(floor, button)
		}
	}
	fmt.Println("Initialised internalQ")
}

// Receive an order (from elevio - Later recieved from DECISION)
func internalQRecieveOrder(order config.Order) {
	internalQSet(order.Floor, int(order.OrderT))
	fmt.Println("Order recieved:")
	internalQPrint()
}

// Remove an order while handling the direction
func internalQRemoveOrder(floor int, currentDirection ElevDir) {
	internalQPop(floor, int(config.Cab))

	// Get both orders regardless of dir if no other orders
	if !(internalQCheckBelow(floor) || internalQCheckAbove(floor)) {
		internalQPop(floor, int(config.HallUp))
		internalQPop(floor, int(config.HallDn))
		return
	}

	switch currentDirection {
	case Up:
		internalQPop(floor, int(config.HallUp))
		if internalQCheckAbove(floor) == false {
			internalQPop(floor, int(config.HallDn))
		}
		break
	case Down:
		internalQPop(floor, int(config.HallDn))
		if internalQCheckBelow(floor) == false {
			internalQPop(floor, int(config.HallUp))
		}
		break
	case Stop:
		internalQPop(floor, int(config.HallUp))
		internalQPop(floor, int(config.HallDn))
	}
}

// Returns an elevator direction after checking current direction and orders
func internalQReturnElevDir(currentFloor int, currentDirection ElevDir) ElevDir {
	switch currentDirection {
	case Up:
		if internalQCheckAbove(currentFloor) == true {
			return currentDirection
		} else if internalQCheckBelow(currentFloor) == true {
			return Down
		}
	case Down:
		if internalQCheckBelow(currentFloor) == true {
			return currentDirection
		} else if internalQCheckAbove(currentFloor) == true {
			return Up
		}
	case Stop:
		if internalQCheckAbove(currentFloor) == true { // Order of check functions not important
			return Up
		} else if internalQCheckBelow(currentFloor) == true {
			return Down
		}
	}
	return Stop // No orders for exist for current direction
}

// Returns true if there exists an order on current floor with SAME direction OR if no other orders beyond
func internalQCheckThisFloorThisDir(currentFloor int, currentDirection ElevDir) bool {
	// Check current floor for same dir
	if internalQueue[currentFloor][config.Cab] {
		return true
	} else if (currentDirection == Up || currentDirection == Stop) && internalQueue[currentFloor][config.HallUp] {
		return true
	} else if (currentDirection == Down || currentDirection == Stop) && internalQueue[currentFloor][config.HallDn] {
		return true
	}

	// Check current floor for no orders beyond
	if currentDirection == Up && internalQCheckAbove(currentFloor) == false {
		return true
	} else if currentDirection == Down && internalQCheckBelow(currentFloor) == false {
		return true
	}
	return false
}

// Returns true if new order matches current direction
func internalCheckThisOrderThisDir(orderT config.OrderType, currentDirection ElevDir) bool {
	if orderT == config.Cab {
		return true
	} else if (currentDirection == Up || currentDirection == Stop) && orderT == config.HallUp {
		return true
	} else if (currentDirection == Down || currentDirection == Stop) && orderT == config.HallDn {
		return true
	} else {
		return false
	}
}

// Returns true if there exists an order strictly above current floor
func internalQCheckAbove(currentFloor int) bool {
	for floor := currentFloor + 1; floor < config.NFloors; floor++ {
		for button := 0; button < config.NButtonTypes; button++ {
			if internalQueue[floor][button] == true {
				return true
			}
		}
	}
	return false
}

// Returns true if there exists an order strictly below current floor
func internalQCheckBelow(currentFloor int) bool {
	for floor := currentFloor - 1; floor > -1; floor-- {
		for button := 0; button < config.NButtonTypes; button++ {
			if internalQueue[floor][button] == true {
				return true
			}
		}
	}
	return false
}

// Fill up internalQueue by iterating through SharedQ
func internalQCheckSharedQ(sharedQ []int) {
	/*
		for i, order := range sharedQ {
			if order.id == myNodeID {
				fmt.Println("adding order ", i, ": ", order)
				}
			}
	*/
}

// Set an order in internalQueue
func internalQSet(floor int, buttonType int) {
	internalQueue[floor][buttonType] = true
}

// Remove an order from internalQueue
func internalQPop(floor int, buttonType int) {
	internalQueue[floor][buttonType] = false
}

// Return an element of intenalQueue - true if order exists
func internalQGet(floor int, buttonType int) bool {
	return internalQueue[floor][buttonType]
}

// Prints out the queue in a cool table
func internalQPrint() {
	fmt.Println("\n   HallUp   HallDn    Cab  ")
	fmt.Println("-" + strings.Repeat("|-------|", config.NButtonTypes))
	for floor := config.NFloors - 1; floor > -1; floor-- {
		fmt.Print(floor)
		for button := 0; button < config.NButtonTypes; button++ {
			i := internalQueue[floor][button]
			if i {
				fmt.Print("| ", "true ", " |")
			} else {
				fmt.Print("| ", "_____", " |")
			}
		}
		fmt.Println()
	}
	fmt.Print("-"+strings.Repeat("---------", config.NButtonTypes), "\n\n")
}
