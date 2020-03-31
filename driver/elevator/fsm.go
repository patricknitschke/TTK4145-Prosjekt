package elevator

import (
	"fmt"
	"time"

	"../config"
	"../elevio"
)

/* FSM module functions */

/* Exported Functions */

// FsmRun controls the Elevator behaviour, start and when in motion
func FsmRun(driverChans config.DriverChans, decLocalOrders <-chan config.Order) {
	go fsmPollOrders(decLocalOrders) // Receives new orders from decision

	// Sets the elevator in motion if not already
	go fsmPollStart()

	// Handles elevator events
	for {
		select {
		case a := <-driverChans.DrvFloors:
			fsmOnNewFloor(a)
		case a := <-driverChans.DrvObstr:
			fsmOnObstruction(a)
		case a := <-driverChans.DrvStop:
			fsmOnStop(a)
		default:
			// Run continous code here - other FSM stuff
		}
	}
}

/* Unexported functions */

// Check internalQueue and serve order if it exists
func fsmOnNewFloor(newFloor int) {
	elevatorSetNewFloor(newFloor)

	// Handles orders in same dir or if none beyond
	if internalQCheckThisFloorThisDir(newFloor, elevator.Dir) {
		elevatorSetMotorDir(Stop) // Stop motor, maintain ElevDir
		elevatorDoorState()
		internalQRemoveOrder(newFloor, elevator.Dir)
		elevatorLightsMatchQueue()
	}

	// New direction after internalQ updated
	elevatorSetDir(internalQReturnElevDir(newFloor, elevator.Dir))
	elevatorPrint()
}

// Receive an order from decision module and set internalQ
func fsmPollOrders(decLocalOrders <-chan config.Order) {
	for {
		o := <-decLocalOrders
		fmt.Print("Received new order: ")
		fmt.Printf("%+v\n", o)

		internalQRecieveOrder(o)
		elevatorLightsMatchQueue()

		// Start elevator if off
		if elevator.Dir == Stop && elevator.startDrivingFlag == false {
			elevator.startDrivingFlag = true
		}

		// Stall the elevator if you want to get in, given same direction
		if elevio.GetFloor() == o.Floor && internalCheckThisOrderThisDir(o.OrderT, elevator.Dir) {
			elevatorEnterDoorState()
		}
	}
}

// Sets the elevator in motion based on a flag
func fsmPollStart() {
	for {
		if elevator.startDrivingFlag == true {
			fsmOnNewFloor(elevator.CurrentFloor) // Start == arrive on current floor
			elevator.startDrivingFlag = false
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
	elevio.SetStopLamp(true)
	elevatorInit()
	internalQInit()
	elevatorLightsMatchQueue()
	stopTimer := time.NewTimer(500 * time.Millisecond)
	<-stopTimer.C
	elevio.SetStopLamp(false)
}
