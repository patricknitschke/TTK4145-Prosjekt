package elevator

import (
	"fmt"
	"time"

	"../config"
	"../elevio"
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
	CurrentFloor     int
	Dir              ElevDir
	startDrivingFlag bool
	doorTimer        *time.Timer
	DoorState        bool
	MotorFailure     bool
}

var elevator Elevator

/* Elevator module functions */

/* Exported Functions */

// PollDriverChannels updates the driverChannels
func PollDriverChannels(drvChans config.DriverChans, elevNewOrder chan<- config.Order) {
	go elevio.PollButtons(drvChans.DrvButtons)
	go elevio.PollFloorSensor(drvChans.DrvFloors)
	go elevio.PollObstructionSwitch(drvChans.DrvObstr)
	go elevio.PollStopButton(drvChans.DrvStop)

	for {
		select {
		case newButtonPress := <-drvChans.DrvButtons:
			newOrder := config.Order{
				OrderT: config.OrderType(newButtonPress.Button),
				Floor:  newButtonPress.Floor,
			}
			elevNewOrder <- newOrder
		default:
			/* Any other conversions to match our elevator config */
		}
	}
}

// Init the elevator submodules
func Init() {
	elevatorInit()
	internalQInit()
}

//GetState returns our local elevator variable
func GetState() Elevator {
	return elevator
}

/* Unexported Functions */

// Moves elevator to first floor and set elevator state
func elevatorInit() {
	if elevio.GetFloor() == -1 {
		elevatorSetDir(Down)
	}
	for elevio.GetFloor() == -1 { // Busy wait until elevator initialised
	}
	elevatorSetDir(Stop)
	elevator.CurrentFloor = elevio.GetFloor()
	elevator.doorTimer = time.NewTimer(0)
	elevator.startDrivingFlag = false
	elevator.DoorState = false
	elevator.MotorFailure = false
	fmt.Println("Initialised elevator.")
}

// Sets the elevator state when we arrive to a new floor (Handles edge cases)
func elevatorSetNewFloor(newFloor int) {
	fmt.Print("Arrived at new floor : ")
	fmt.Printf("%+v\n", newFloor)

	elevator.CurrentFloor = newFloor
	elevio.SetFloorIndicator(newFloor)
	switch newFloor {
	case config.NFloors - 1:
		elevatorSetDir(Down)
		break
	case 0:
		elevatorSetDir(Up)
		break
	}
}

// Matches the elevator lights with the current internalQueue
func elevatorLightsMatchQueue() {
	for floor := 0; floor < config.NFloors; floor++ {
		for button := 0; button < config.NButtonTypes; button++ {
			if internalQGet(floor, button) == true {
				elevio.SetButtonLamp(elevio.ButtonType(button), floor, true)
			} else {
				elevio.SetButtonLamp(elevio.ButtonType(button), floor, false)
			}
		}
	}
}

// Set a new direction and update state
func elevatorSetDir(newDirection ElevDir) {
	elevator.Dir = newDirection
	elevatorSetMotorDir(newDirection)
}

// Set a new direction while maintaining state
func elevatorSetMotorDir(newDirection ElevDir) {
	elevio.SetMotorDirection(elevio.MotorDirection(newDirection))
}

// Set the door lamp and start timer
func elevatorEnterDoorState() {
	fmt.Print("Enter door state.")
	elevio.SetDoorOpenLamp(true)
	elevator.DoorState = true

	// Door - Stop, empty and reset timer
	elevator.doorTimer.Stop()
	select {
	case <-elevator.doorTimer.C:
	default:
	}
	elevator.doorTimer.Reset(config.DoorTimerSec * time.Second)
}

// DoorState is a 2 second timer with a lamp
func elevatorDoorState() {
	defer fmt.Print("Exit door state.")
	defer elevio.SetDoorOpenLamp(false)

	elevatorEnterDoorState()
	<-elevator.doorTimer.C
	elevator.DoorState = false
}

// Print out current state of the elevator and queue
func elevatorPrint() {
	fmt.Print("\n-------Elevator state-------\n\n")
	fmt.Println("Current floor: ", elevator.CurrentFloor)
	fmt.Println("Direction: ", elevator.Dir)
	internalQPrint()
}
