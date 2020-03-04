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

// Receive an order from decision module
func fsmOnOrderReceived(o Order) {
	fmt.Print("Received new order: ")
	fmt.Printf("%+v\n", o)

	// Handle order for current floor directly (if stopped), then internalQueue unnecessary
	if o.floor == elevatorGetFloor() && elevatorGetDir() == Stop {
		elevatorSetLight(o.floor, int(o.orderT), true)
		fsmDoorState()
		elevatorSetLight(o.floor, int(o.orderT), false)
		return
	}

	internalQRecieveOrder(o)
	elevatorLightsMatchQueue()

	// Set the elevator in motion if not already
	if elevatorGetDir() == Stop {
		elevatorSetDir(internalQReturnElevDir(elevatorGetFloor(), elevatorGetDir()))
	}
}

// Recieves a ButtonEvent and sends it to decision module
func fsmOnButtonRequest(a elevio.ButtonEvent) {
	orderT := ButtonTypeToOrderTypeMap[a.Button]
	order := Order{orderT, a.Floor}

	decisionDecide(order)
}

// Runs fsmOnButtonRequest forever (in go routine)
func fsmPollButtonRequest(drvButtons chan elevio.ButtonEvent) {
	for {
		fsmOnButtonRequest(<-drvButtons)
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
		decisionInit()
		elevatorLightsMatchQueue()
	}
}

// Just restarts the elevator
func fsmOnStop(a bool) {
	fmt.Print("Stop button pressed:")
	fmt.Printf("%+v\n", a)
	elevatorInit()
	internalQInit()
	decisionInit()
	elevatorLightsMatchQueue()
}
