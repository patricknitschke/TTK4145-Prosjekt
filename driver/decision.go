package main

import (
	"fmt"
	"time"

	"./elevio"
)

/* Constants and variables */

const _pollRate = 20 * time.Millisecond

// NNodes represents the no. of elevator nodes in our network
const NNodes = 3

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

// An order buffer that allows internalQ to pick up orders
var orderBuffer chan Order

// ButtonTypeToOrderTypeMap solves the trouble of having two enums
var ButtonTypeToOrderTypeMap = map[elevio.ButtonType]OrderType{
	elevio.BT_HallUp:   HallUp,
	elevio.BT_HallDown: HallDn,
	elevio.BT_Cab:      Cab,
}

//
//
//

/* Decision module functions */

// Initialise decision module and variables
func decisionInit() {
	orderBuffer = make(chan Order)
	fmt.Println("Initialised decision")
}

// Samples the costs and sends an order to a node
func decisionDecide(order Order) {
	deliveryNode := decisionSampleCosts(order)
	fmt.Print("Order to 'Node ")
	fmt.Print(deliveryNode)
	fmt.Println("' decided.")

	// Send to sharedQ
	myNode := 1
	if deliveryNode == myNode {
		go func() {
			for {
				fmt.Print("Order sending...")
				orderBuffer <- order
				fmt.Print("Order sent.")
				return
			}
		}()
	}
}

// Delivers an order to internalQ once channel has order
func decisionPollOrders(receiver chan<- Order) {
	for {
		select {
		case o := <-orderBuffer:
			receiver <- o
		}
	}
}

// Obtains the cost to serve an order from each node and returns the minimum node ID
func decisionSampleCosts(order Order) int {

	// Sample costs using sharedQ
	costsToServe := [NNodes]int{5, 3, 7}

	// Return minimum node
	var minIndex int
	var minCost int
	for index, cost := range costsToServe {
		if index == 0 || cost < minCost {
			minCost = cost
			minIndex = index
		}
	}
	return minIndex
}

// Calculates the cost to serve an order based
func decisionCostFunction(order Order) {

}

// Listens for an order confirmation from a node
func decisionConfirm() {

}

// Sends an order to internalQ when receiving an order from another node
func decisionRecieveOrder() {

}

// Acknowledges an order to sender node (or all nodes) once light switched on
func decisionAckOrder() {

}

// Redistributes internalQ to other nodes under some motor failure
func decisionRedistributeQ() {

}
