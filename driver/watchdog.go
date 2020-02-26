package main

import (
	"fmt"
	"time"
)

const queueInterval = 15    // seconds
const imAliveInterval = 100 // microseconds

var aliveTimer [NNodes]*time.Timer
var queueTimer [NFloors][NButtonTypes]*time.Timer

func watchPollAliveTimeout() {
	for {
		select {
		case <-aliveTimer[1].C:

			//	case ... : how to select on N nodes?

		}
	}
}

func watchPollQueueTimeout() {

}

func watchSetAlive(elevatorID int) {
	aliveTimer[elevatorID] = time.NewTimer(imAliveInterval * time.Millisecond)
}

func watchSetQueue(order Order) {
	queueTimer[order.floor][int(order.orderT)] = time.NewTimer(queueInterval * time.Second)
}

func watchBcastImAlive() {
	for now := range time.Tick(imAliveInterval * time.Millisecond) {
		//newtorkBcastImalive()
		fmt.Println(now)
	}
}

func watchAlert() {

}

func watchCheckWhoAlive() {

}
