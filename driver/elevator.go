package main

import (
	"fmt"
	"time"

	"./elevio"
)

/* Constants and variables */

// ElevDir is an enum for elevator direction - matches elevio.MotorDirection
type ElevDir int

// Up, Down and Stop are ElevDir symbols
const (
	Up   ElevDir = 1
	Down         = -1
	Stop         = 0
)

// Elevator contains current state
type Elevator struct {
	currentFloor     int
	dir              ElevDir
	startDrivingFlag bool
	doorTimer        *time.Timer
	doorState        bool
	MotorFailure     bool
}

var elevator Elevator

/* Elevator module functions */

// Move elevator to first floor and set elevator state
func elevatorInit() {
	if elevio.GetFloor() == -1 {
		elevatorSetDir(Down)
	}
	for elevio.GetFloor() == -1 { // Busy wait until elevator initialised
	}
	elevatorSetDir(Stop)
	elevatorSetFloor(elevio.GetFloor())
	elevatorSetDriveFlag(false)
	elevator.doorTimer = time.NewTimer(0)
	elevator.doorState = false
	elevator.MotorFailure = false
	elevatorPrint()
	fmt.Println("Initialised elevator.")
}

// Sets the elevator state when we arrive to a new floor (Handles edge cases)
func elevatorSetNewFloor(newFloor int) {
	fmt.Print("Arrived at new floor : ")
	fmt.Printf("%+v\n", newFloor)

	elevatorSetFloor(newFloor)
	switch newFloor {
	case NFloors - 1:
		elevatorSetDir(Down)
		break
	case 0:
		elevatorSetDir(Up)
		break
	}
	elevatorPrint()
}

// Matches the elevator lights with the current internalQueue
func elevatorLightsMatchQueue() {
	for floor := 0; floor < NFloors; floor++ {
		for button := 0; button < NButtonTypes; button++ {
			if internalQGet(floor, button) == true {
				elevatorSetLight(floor, button, true)
			} else {
				elevatorSetLight(floor, button, false)
			}
		}
	}
}

// Set a new direction and update state
func elevatorSetDir(newDirection ElevDir) {
	elevator.dir = newDirection
	elevatorSetMotorDir(newDirection)
}

// Set a new direction while maintaining state
func elevatorSetMotorDir(newDirection ElevDir) {
	elevio.SetMotorDirection(elevio.MotorDirection(newDirection))
}

// Set a new floor state
func elevatorSetFloor(newFloor int) {
	elevator.currentFloor = newFloor
}

// Sets the elevator door state in fsmDoor
func elevatorSetDriveFlag(startDrivingFlag bool) {
	elevator.startDrivingFlag = startDrivingFlag
}

// Sets a button lamp to state
func elevatorSetLight(floor int, button int, state bool) {
	elevio.SetButtonLamp(elevio.ButtonType(button), floor, state)
}

// Return the elevator direction
func elevatorGetDir() ElevDir {
	return elevator.dir
}

// Return the elevator current floor
func elevatorGetFloor() int {
	return elevator.currentFloor
}

// Return the elevator door state
func elevatorGetDriveFlag() bool {
	return elevator.startDrivingFlag
}

// Set the door lamp and start timer
func elevatorEnterDoorState() {
	fmt.Print("Enter door state.")
	elevio.SetDoorOpenLamp(true)
	elevator.doorState = true

	// Door - Stop, empty and reset timer
	elevator.doorTimer.Stop()
	select {
	case <-elevator.doorTimer.C:
	default:
	}
	elevator.doorTimer.Reset(doorTimerSec * time.Second)
}

// DoorState is a 2 second timer with a lamp
func elevatorDoorState() {
	defer fmt.Print("Exit door state.")
	defer elevio.SetDoorOpenLamp(false)

	elevatorEnterDoorState()
	<-elevator.doorTimer.C
	elevator.doorState = false
}

// Print out current state of the elevator and queue
func elevatorPrint() {
	fmt.Print("\n-------Elevator state-------\n\n")
	fmt.Println("Current floor: ", elevator.currentFloor)
	fmt.Println("Direction: ", elevator.dir)
	internalQPrint()
}
