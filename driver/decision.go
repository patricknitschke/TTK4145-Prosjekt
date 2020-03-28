package main

import (
	"fmt"

	"./elevio"
)

/* Constants and variables */

/* Decision module functions */

// Initialise decision module and variables
func decisionInit() {
	fmt.Println("Initialised decision")
}

// Samples the costs and sends an order to a node
func decisionDecide(newOrder Order, decLocalOrders chan<- Order, decSharedQOrders chan<- Order) {
	deliveryNode := decisionCostFunction(newOrder)

	if deliveryNode == LocalID {
		for {
			decLocalOrders <- newOrder
			return
		}
	} else {
		for {
			// Am I supposed to send to SharedQ though? I know the best elevator
			decSharedQOrders <- newOrder

			// maybe decNetOrders with new struct, given they're online
			// decNetOrders <- NetworkPacket{ElevID, newOrder}
			return
		}
	}
}

// Continuosly sends orders to decision, which decides between internalQ or other nodes
func decisionPollButtonRequests(drvButtons <-chan elevio.ButtonEvent, decLocalOrders chan<- Order, decSharedQOrders chan<- Order) {
	for {
		a := <-drvButtons
		orderT := OrderType(a.Button)
		order := Order{orderT, a.Floor}

		decisionDecide(order, decLocalOrders, decSharedQOrders)
	}
}

// Calculates the cost of an order based on the Elevator positions and state, returns optimal Node
func decisionCostFunction(newOrder Order) int {

	// NNodes is 0 now, remember to change in config
	var elevatorList [NNodes]Elevator
	var elevatorOnline [NNodes]bool

	// TODO Get elevatorList and onlineList somehow
	elevatorList[LocalID] = elevator // My state
	elevatorOnline[LocalID] = true
	// TODO ///////////////////////////////////////

	if newOrder.orderT == Cab {
		return LocalID
	}
	var CostArray [NNodes]int
	for elev := 0; elev < NNodes; elev++ {
		cost := newOrder.floor - elevatorList[elev].currentFloor
		if cost == 0 && elevatorOnline[elev] && !elevatorList[elev].MotorFailure {
			return elev
		}
		if cost < 0 {
			cost = -cost
			if elevatorList[elev].dir == Up {
				cost += 3
			}
		} else if cost > 0 {
			if elevatorList[elev].dir == Down {
				cost += 3
			}
		}
		if elevatorList[elev].doorState == true {
			cost++
		}
		CostArray[elev] = cost
	}
	fmt.Print("Cost to Serve: ")
	fmt.Println(CostArray)
	maxCost := 1000
	var bestElev int
	for elev := 0; elev < NNodes; elev++ {
		if CostArray[elev] < maxCost && elevatorOnline[elev] && !elevatorList[elev].MotorFailure {
			bestElev = elev
			maxCost = CostArray[elev]
		}
	}
	return bestElev
}

// Listens for an order confirmation from a node
func decisionConfirm() {

}

// Sends an order to internalQ when receiving an order from another node (Or sharedQ???)
func decisionRecieveOrder() {

}

// Acknowledges an order to sender node (or all nodes) once light switched on
func decisionAckOrder() {

}

// Redistributes internalQ to other nodes under some motor failure
func decisionRedistributeQ() {

}
