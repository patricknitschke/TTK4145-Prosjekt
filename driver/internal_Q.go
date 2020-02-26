package main

import (
	"fmt"
	"strings"

	"./elevio"
)

/* Constants and variables */

// NFloors represents the number of floors - make sure to match elevio file
const NFloors = 4

// NButtonTypes are "var order elevio.ButtonType : HallUp , HallDown, Cab" - from elevio
const NButtonTypes = 3

// OrderType is an enum for internalQueue - matches elevio.ButtonType
type OrderType int

const (
	HallUp OrderType = 0
	HallDn           = 1
	Cab              = 2
)

// Order makes things simpler for other modules
type Order struct {
	orderT OrderType
	floor  int
}

// OrderToButtonTypesMap solves the trouble of having two enums
var OrderToButtonTypesMap = map[OrderType]elevio.ButtonType{
	HallUp: elevio.BT_HallUp,
	HallDn: elevio.BT_HallDown,
	Cab:    elevio.BT_Cab,
}

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

// Recieve an order (from elevio - Later recieved from DECISION)
func internalQRecieveOrder(order elevio.ButtonEvent) {
	orderT := int(order.Button)
	internalQSet(order.Floor, orderT)
	fmt.Println("Order recieved:")
	internalQPrint()
}

// Remove an order
func internalQRemoveForDir(floor int, currentDirection ElevDir) {
	internalQPop(floor, int(Cab))
	switch currentDirection {
	case Up:
		internalQPop(floor, int(HallUp))
	case Down:
		internalQPop(floor, int(HallDn))
	case Stop:
		internalQPop(floor, int(HallUp))
		internalQPop(floor, int(HallDn))
	}
}

// Set an order in internalQueue
func internalQSet(floor int, buttonType int) {
	internalQueue[floor][buttonType] = true
}

// Remove an order from internalQueue
func internalQPop(floor int, buttonType int) {
	internalQueue[floor][buttonType] = false
}

func internalQGet(floor int, buttonType int) bool {
	return internalQueue[floor][buttonType]
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

// Returns true if there exists an order on current floor
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
		if internalQCheckAbove(currentFloor) == true {
			return Up
		} else if internalQCheckBelow(currentFloor) == true {
			return Down
		}
	}
	return Stop // No orders for exist despite directions
}

// Prints out the queue in a nice table
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
