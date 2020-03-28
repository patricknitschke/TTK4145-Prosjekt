package main

import (
	"./elevio"
)

func main() {
	elevio.Init("localhost:15657", NFloors) // Connect to Elevator Server

	// Initialise modules here
	elevatorInit()
	internalQInit()
	decisionInit()
	fsmInit()

	// Channels
	drvChans := DriverChans{
		make(chan elevio.ButtonEvent),
		make(chan int),
		make(chan bool),
		make(chan bool)}

	// Orders to Local Elevator and SharedQ
	decLocalOrders := make(chan Order)
	decSharedQOrders := make(chan Order)

	// Each goroutine updates the channels when they're available
	go elevio.PollButtons(drvChans.drvButtons)
	go elevio.PollFloorSensor(drvChans.drvFloors)
	go elevio.PollObstructionSwitch(drvChans.drvObstr)
	go elevio.PollStopButton(drvChans.drvStop)

	go decisionPollButtonRequests(drvChans.drvButtons, decLocalOrders, decSharedQOrders) // Continuously handles the pressing of buttons

	// Local FSM responds to events on the channels
	go fsmPollOrders(decLocalOrders) // Handles new orders from decision (decOrders)
	go fsmDriving(drvChans)          // Starts elev if stopped and incoming order

	select {}
}
