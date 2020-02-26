package main

import (
	"./elevio"
)

// NNodes is the number of nodes in our network
const NNodes = 3

func main() {
	elevio.Init("localhost:15657", NFloors) // Connect to Elevator Server

	// Initialise modules here
	elevatorInit()
	internalQInit()

	// Channels
	drvButtons := make(chan elevio.ButtonEvent)
	drvFloors := make(chan int)
	drvObstr := make(chan bool)
	drvStop := make(chan bool)

	// Each goroutine updates the channels when they're available
	go elevio.PollButtons(drvButtons)
	go elevio.PollFloorSensor(drvFloors)
	go elevio.PollObstructionSwitch(drvObstr)
	go elevio.PollStopButton(drvStop)

	// FSM responds to events on the channels
	go fsmPollButtonRequest(drvButtons) // Continuously handles the pressing of buttons

	for {
		select {
		case a := <-drvFloors:
			fsmOnNewFloor(a)

		case a := <-drvObstr:
			fsmOnObstruction(a)

		case a := <-drvStop:
			fsmOnStop(a)

		default:
			// Run continous code here - other FSM stuff
		}
	}

}
