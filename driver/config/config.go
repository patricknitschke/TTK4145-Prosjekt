package config

import (
	"../elevio"
)

// LocalID is this node's ID, where should it though?
const LocalID = 0

// NNodes represents the no. of elevator nodes in our network
const NNodes = 1

// NFloors represents the number of floors - make sure to match elevio file
const NFloors = 9

// NButtonTypes are "var order elevio.ButtonType : HallUp , HallDown, Cab" - from elevio
const NButtonTypes = 3

// DoorTimerSec represents door timer time
const DoorTimerSec = 1

// DriverChans contains the channels used when the Elevator moves.
type DriverChans struct {
	DrvButtons chan elevio.ButtonEvent
	DrvFloors  chan int
	DrvObstr   chan bool
	DrvStop    chan bool
}

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
	OrderT OrderType
	Floor  int
	// ? ID int
}
