package main

import "./elevio"

// LocalID is this nodes ID, where should it though?
const LocalID = 0

// NNodes represents the no. of elevator nodes in our network
const NNodes = 1

// NFloors represents the number of floors - make sure to match elevio file
const NFloors = 4

// NButtonTypes are "var order elevio.ButtonType : HallUp , HallDown, Cab" - from elevio
const NButtonTypes = 3

// Number of seconds for door closed
const doorTimerSec = 2

// OrderType is an enum for decision - matches elevio.ButtonType
type OrderType int

// HallUp, HallDn and Cab are the OrderType symbols
const (
	HallUp OrderType = 0
	HallDn           = 1
	Cab              = 2
)

// Order makes things simpler for other modules
type Order struct {
	orderT OrderType
	floor  int
}

// DriverChans contains the channels used when the Elevator moves.
type DriverChans struct {
	drvButtons chan elevio.ButtonEvent
	drvFloors  chan int
	drvObstr   chan bool
	drvStop    chan bool
}
