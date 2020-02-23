package main

import (
	"fmt"
)

/* Constants and variables */

// NFloors represents the number of floors - make sure to match elevio file
const NFloors = 4

// NButtonTypes are "var order elevio.ButtonType : HallUp, HallDown, Cab" - from elevio
const NButtonTypes = 3

var internalQueue [NFloors][NButtonTypes]bool

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

// Set an order in internalQueue
func internalQSet(floor int, buttonType int) {
	internalQueue[floor][buttonType] = true
}

// Remove an order from internalQueue
func internalQPop(floor int, buttonType int) {
	internalQueue[floor][buttonType] = false
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

// Returns true if there exists an order above current floor
func internalQCheckAbove(currentFloor int) bool {
	for floor := currentFloor; floor < NFloors; floor++ {
		for button := 0; button < NButtonTypes; button++ {
			if internalQueue[floor][button] == true {
				return true
			}
		}
	}
	return false
}

// Returns true if there exists an order below current floor
func internalQCheckBelow(currentFloor int) bool {
	for floor := currentFloor; floor > -1; floor-- {
		for button := 0; button < NButtonTypes; button++ {
			if internalQueue[floor][button] == true {
				return true
			}
		}
	}
	return false
}

// Returns an elevator direction after checking current direction and orders
func internalQReturnElevDir(currentFloor int, currentDirection int) int {
	switch currentDirection {
	case 1: // Going up
		if internalQCheckAbove(currentFloor) == true {
			return currentDirection
		} else if internalQCheckBelow(currentFloor) == true {
			return -1
		}

	case -1: // Going down
		if internalQCheckBelow(currentFloor) == true {
			return currentDirection
		} else if internalQCheckAbove(currentFloor) == true {
			return 1
		}

	default: // Stopped, case 0:
		if internalQCheckAbove(currentFloor) == true {
			return currentDirection
		} else if internalQCheckBelow(currentFloor) == true {
			return -1
		}
	}
	return 0 // No orders for exists despite directions - stop
}
