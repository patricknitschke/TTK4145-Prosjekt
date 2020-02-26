package main

import (
	"fmt"

	"./elevio"
)

// ElevDir is an enum for elevator direction - matches elevio.MotorDirection
type ElevDir int

const (
	Up   ElevDir = 1
	Down         = -1
	Stop         = 0
)

// Elevator contains current state
type Elevator struct {
	currentFloor int
	dir          ElevDir
}

var elevator Elevator

// Move elevator to first floor and set elevator state
func elevatorInit() {
	if elevio.GetFloor() == -1 {
		elevatorSetDir(Down)
	}
	for elevio.GetFloor() == -1 { // Busy wait until elevator initialised
	}
	elevatorSetDir(Stop)
	elevatorSetFloor(elevio.GetFloor())
	elevatorPrint()
	fmt.Println("Initialised elevator.")
}

// Remove orders from queue and ...?
func elevatorServeOrder(floor int, dir ElevDir) {
	internalQRemoveForDir(floor, dir)
}

// Set a new direction and update state
func elevatorSetDir(newDirection ElevDir) {
	elevator.dir = newDirection
	elevatorSetMotorDir(newDirection)
}

// Set a new direction and update state
func elevatorSetMotorDir(newDirection ElevDir) {
	elevio.SetMotorDirection(elevio.MotorDirection(newDirection))
}

// Set a new floor state
func elevatorSetFloor(newFloor int) {
	elevator.currentFloor = newFloor
}

// Return the elevator direction
func elevatorGetDir() ElevDir {
	return elevator.dir
}

// Return the elevator current floor
func elevatorGetFloor() int {
	return elevator.currentFloor
}

// Print out current state of the elevator and queue
func elevatorPrint() {
	fmt.Print("\n-------Elevator state-------\n\n")
	fmt.Println("Current floor: ", elevator.currentFloor)
	fmt.Println("Direction: ", elevator.dir)
	internalQPrint()
}
