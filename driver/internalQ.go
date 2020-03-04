package main

import (
	"fmt"
	"strings"
)

/* Constants and variables */

// NFloors represents the number of floors - make sure to match elevio file
const NFloors = 4

// NButtonTypes are "var order elevio.ButtonType : HallUp , HallDown, Cab" - from elevio
const NButtonTypes = 3

// InternalQueue maintains orders for its node
var internalQueue [NFloors][NButtonTypes]bool

/* HallUp   HallDn    Cab
-|-------||-------||-------|
3| _____ || _____ || _____ |
2| _____ || _____ || _____ |
1| _____ || _____ || _____ |
0| _____ || _____ || _____ |
---------------------------- */

/* InternalQ module functions */

// Set up internalQueue
func internalQInit() {
	for floor := 0; floor < NFloors; floor++ {
		for button := 0; button < NButtonTypes; button++ {
			internalQPop(floor, button)
		}
	}
	fmt.Println("Initialised internalQ")
}

// Receive an order (from elevio - Later recieved from DECISION)
func internalQRecieveOrder(order Order) {
	internalQSet(order.floor, int(order.orderT))
	fmt.Println("Order recieved:")
	internalQPrint()
}

// Remove an order while handling the direction
func internalQRemoveOrder(floor int, currentDirection ElevDir) {
	internalQPop(floor, int(Cab))

	// Get both orders regardless of dir if no other orders
	if !(internalQCheckBelow(floor) || internalQCheckAbove(floor)) {
		internalQPop(floor, int(HallUp))
		internalQPop(floor, int(HallDn))
		return
	}

	switch currentDirection {
	case Up:
		internalQPop(floor, int(HallUp))
		break
	case Down:
		internalQPop(floor, int(HallDn))
		break
	case Stop:
		internalQPop(floor, int(HallUp))
		internalQPop(floor, int(HallDn))
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

// Returns true if there exists an order on current floor with SAME direction
func internalQCheckThisFloorSameDir(currentFloor int, currentDirection ElevDir) bool {
	if internalQueue[currentFloor][Cab] {
		return true
	} else if (currentDirection == Up || currentDirection == Stop) && internalQueue[currentFloor][HallUp] {
		return true
	} else if (currentDirection == Down || currentDirection == Stop) && internalQueue[currentFloor][HallDn] {
		return true
	}
	return false
}

// Returns true if there exists an order strictly above current floor
func internalQCheckAbove(currentFloor int) bool {
	for floor := currentFloor + 1; floor < NFloors; floor++ {
		for button := 0; button < NButtonTypes; button++ {
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
		for button := 0; button < NButtonTypes; button++ {
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
	fmt.Println("-" + strings.Repeat("|-------|", NButtonTypes))
	for floor := NFloors - 1; floor > -1; floor-- {
		fmt.Print(floor)
		for button := 0; button < NButtonTypes; button++ {
			i := internalQueue[floor][button]
			if i {
				fmt.Print("| ", "true ", " |")
			} else {
				fmt.Print("| ", "_____", " |")
			}
		}
		fmt.Println()
	}
	fmt.Print("-"+strings.Repeat("---------", NButtonTypes), "\n\n")
}
