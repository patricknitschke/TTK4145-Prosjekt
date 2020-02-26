package main

import (
	"fmt"
	"time"
)

const queueInterval = 15    // seconds
const imAliveInterval = 100 // microseconds

// Maybe variables could be structs instead of *Timers?
var aliveTimer [NNodes]*time.Timer
var queueTimer [NFloors][NButtonTypes]*time.Timer

// Forever watches for aliveTimeouts -> calls watchAlert
func watchPollAliveTimeout() {
	for {
		select {
		case <-aliveTimer[1].C:

			//	case ... : how to select on N nodes?

		}
	}
}

// Forever watches for queueTimeouts -> calls watchAlert
func watchPollQueueTimeout() {

}

// Function that resets internalQ, asks for help and stuff
func watchAlert() {

}

// Restarts the imAlive timer
func watchSetAlive(elevatorID int) {
	aliveTimer[elevatorID] = time.NewTimer(imAliveInterval * time.Millisecond)
}

// Starts a timer for a particular order
func watchSetQueue(order Order) {
	queueTimer[order.floor][int(order.orderT)] = time.NewTimer(queueInterval * time.Second)
}

// Disables a timer for a particular order
func watchSPopQueue(order Order) {

}

// Timer for regular interval ImAlive broadcasts
func watchBcastImAlive() {
	for now := range time.Tick(imAliveInterval * time.Millisecond) {
		//newtorkBcastImalive()
		fmt.Println(now)
	}
}

// Returns an bool array ([]bool) of alive nodes - could be used to update Decision module's static var?
func watchCheckWhoAlive() {

}
