package main

import (
	"fmt"
	"time"

	"./elevio"
)

// Check internalQueue and serve order if it exists
func fsmOnNewFloor(newFloor int) {
	fmt.Print("Arrived at new floor : ")
	fmt.Printf("%+v\n", newFloor)

	elevatorSetFloor(newFloor)
	if newFloor == 0 || newFloor == NFloors-1 {
		elevatorSetDir(Stop)
	}
	elevatorPrint()

	if internalQCheckThisFloorSameDir(newFloor, elevatorGetDir()) {
		elevatorSetMotorDir(Stop)
		timer1 := time.NewTimer(2 * time.Second)
		<-timer1.C
		internalQRemoveForDir(newFloor, elevatorGetDir())
	}
	elevatorSetDir(internalQReturnElevDir(newFloor, elevatorGetDir()))
}

// Recieves a ButtonEvent and places it into internalQueue (later send to Decision module instead of internalQ)
func fsmOnButtonRequest(a elevio.ButtonEvent) {
	fmt.Printf("%+v\n", a)
	internalQRecieveOrder(a)
	elevio.SetButtonLamp(a.Button, a.Floor, true)
	dirToGo := internalQReturnElevDir(elevatorGetFloor(), elevatorGetDir())
	elevatorSetDir(dirToGo)
	fmt.Println(dirToGo)
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
