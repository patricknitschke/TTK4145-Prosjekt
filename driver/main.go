package main

import (
	"./config"
	"./decision"
	"./elevator"
	"./elevio"
)

func main() {
	elevio.Init("localhost:15657", config.NFloors) // Connect to Elevator Server

	// Initialise modules here
	elevator.Init()

	// Channels
	drvChans := config.DriverChans{
		DrvButtons: make(chan elevio.ButtonEvent),
		DrvFloors:  make(chan int),
		DrvObstr:   make(chan bool),
		DrvStop:    make(chan bool),
	}

	// Orders to Local Elevator and SharedQ
	elevNewOrder := make(chan config.Order)
	decLocalOrders := make(chan config.Order)
	decSharedQOrders := make(chan config.Order)

	// Update driver channels and take in new Orders
	go elevator.PollDriverChannels(drvChans, elevNewOrder)

	// Local FSM responds to events on the channels
	go elevator.FsmRun(drvChans, decLocalOrders)

	// Continuously handles new Orders
	go decision.Decide(elevNewOrder, decLocalOrders, decSharedQOrders)

	select {}
}
