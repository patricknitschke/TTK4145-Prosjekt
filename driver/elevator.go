package main

import (
	"fmt"
)

// Elevator maintains current state
type Elevator struct {
	dir          int
	currentFloor int
}

var elevator Elevator

func elevatorInit() {
	fmt.Println("Initialised elevator module")
	elevator.currentFloor = 0
	elevator.dir = 0
}
