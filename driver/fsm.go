package main

import (
	"fmt"
	"time"

	"./elevio"
)

// Check internalQueue and serve order if it exists
func fsmOnNewFloor(newFloor int) {
	elevatorSetNewFloor(newFloor)

	// Handle orders in the same direction, if not switch and handle again (Redundancy acceptable - 100% cases accounted for)
	for i := 0; i < 2; i++ {
		if internalQCheckThisFloorSameDir(newFloor, elevatorGetDir()) {
			elevatorSetMotorDir(Stop) // Stop motor, maintain ElevDir
			fsmDoorState()
			internalQRemoveOrder(newFloor, elevatorGetDir())
			elevatorLightsMatchQueue()
		}

		// New direction after internalQ updated
		elevatorSetDir(internalQReturnElevDir(newFloor, elevatorGetDir()))
		elevatorPrint()
	}
}

// Runs fsmOnButtonRequest forever (used in go routine)
func fsmPollButtonRequest(drvButtons chan elevio.ButtonEvent) {
	for {
		fsmOnButtonRequest(<-drvButtons)
	}
}

// Recieves a ButtonEvent and places it into internalQueue (later send to Decision module instead of internalQ)
func fsmOnButtonRequest(a elevio.ButtonEvent) {
	fmt.Print("Received new order: ")
	fmt.Printf("%+v\n", a)

	// Handle order for current floor directly (if stopped), then internalQueue unnecessary
	if a.Floor == elevatorGetFloor() && elevatorGetDir() == Stop {
		fsmDoorState()
		return
	}

	internalQRecieveOrder(a)
	elevatorLightsMatchQueue()

	// Set the elevator in motion if not already
	if elevatorGetDir() == Stop {
		elevatorSetDir(internalQReturnElevDir(elevatorGetFloor(), elevatorGetDir()))
	}
}

// Set the door lamp and pause for 2 seconds
func fsmDoorState() {
	fmt.Print("Enter door state.")
	elevio.SetDoorOpenLamp(true)
	timer1 := time.NewTimer(2 * time.Second)
	<-timer1.C
	elevio.SetDoorOpenLamp(false)
}

// Restart elevator on obstruction
func fsmOnObstruction(a bool) {
	fmt.Print("Obstruction switch active: ")
	fmt.Printf("%+v\n", a)
	if a {
		elevatorSetDir(Stop)
	} else {
		elevatorInit()
		internalQInit()
		elevatorLightsMatchQueue()
	}
}

// Just restarts the elevator
func fsmOnStop(a bool) {
	fmt.Print("Stop button pressed:")
	fmt.Printf("%+v\n", a)
	elevatorInit()
	internalQInit()
	elevatorLightsMatchQueue()
}
