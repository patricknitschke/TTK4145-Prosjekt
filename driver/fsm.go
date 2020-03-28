package main

import (
	"fmt"

	"./elevio"
)

/* Constants and variables */

/* FSM module functions */

// Initialise variables in FSM
func fsmInit() {
	fmt.Println("Initialised fsm")
}

// Controls the Elevator behaviour, start and when in motion
func fsmDriving(driverChans DriverChans) {
	// Sets the elevator in motion if not already
	go func() {
		for {
			if elevatorGetDriveFlag() == true {
				fsmOnNewFloor(elevatorGetFloor()) // Start == arrive on current floor
				elevatorSetDriveFlag(false)
			}
		}
	}()

	// Handles elevator events
	go func() {
		for {
			select {
			case a := <-driverChans.drvFloors:
				fsmOnNewFloor(a)
			case a := <-driverChans.drvObstr:
				fsmOnObstruction(a)
			case a := <-driverChans.drvStop:
				fsmOnStop(a)
			default:
				// Run continous code here - other FSM stuff
			}
		}
	}()
	select {}
}

// Check internalQueue and serve order if it exists
func fsmOnNewFloor(newFloor int) {
	elevatorSetNewFloor(newFloor)

	// Handles orders in same dir or if none beyond
	if internalQCheckThisFloorThisDir(newFloor, elevatorGetDir()) {
		elevatorSetMotorDir(Stop) // Stop motor, maintain ElevDir
		elevatorDoorState()
		internalQRemoveOrder(newFloor, elevatorGetDir())
		elevatorLightsMatchQueue()
	}

	// New direction after internalQ updated
	elevatorSetDir(internalQReturnElevDir(newFloor, elevatorGetDir()))
	elevatorPrint()
}

// Receive an order from decision module and set internalQ
func fsmPollOrders(decOrders <-chan Order) {
	for {
		o := <-decOrders
		fmt.Print("Received new order: ")
		fmt.Printf("%+v\n", o)

		internalQRecieveOrder(o)
		elevatorLightsMatchQueue()

		// Start elevator if off
		if elevatorGetDir() == Stop && elevatorGetDriveFlag() == false {
			elevatorSetDriveFlag(true)
		}

		// Stall the elevator if you want to get in
		if elevio.GetFloor() == o.floor && internalQCheckThisFloorThisDir(o.floor, elevatorGetDir()) {
			elevatorEnterDoorState()
		}
	}
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
