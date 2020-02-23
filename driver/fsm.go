package main

import (
	"fmt"

	"./elevio"
)

// Check internalQueue and serve order if it exists
func fsmOnNewFloor(newFloor int) {
	fmt.Print("Arrived at new floor : ")
	fmt.Printf("%+v\n", newFloor)
	elevatorSetFloor(newFloor)
}

// Recieves a ButtonEvent and places it into internalQueue (later send to Decision module instead of internalQ)
func fsmOnButtonRequest(a elevio.ButtonEvent) {
	fmt.Printf("%+v\n", a)
	internalQRecieveOrder(a)
	elevio.SetButtonLamp(a.Button, a.Floor, true)
}

func fsmOnObstruction(a bool) {
	fmt.Printf("%+v\n", a)
	if a {
		elevatorSetDir(Stop)
	} else {
		elevatorInit()
		internalQInit()
	}
}

func fsmOnStop(a bool) {
	fmt.Printf("%+v\n", a)
	for f := 0; f < NFloors; f++ {
		for b := elevio.ButtonType(0); b < 3; b++ {
			elevio.SetButtonLamp(b, f, false)
		}
	}
}
