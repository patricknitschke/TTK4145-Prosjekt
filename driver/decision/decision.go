package decision

import (
	"../config"
)

/* Constants and variables */

/* Decision module functions */

// Decide samples the costs and sends an order to a node
func Decide(elevNewOrder <-chan config.Order, decLocalOrders chan<- config.Order, decSharedQOrders chan<- config.Order) {
	for {
		newOrder := <-elevNewOrder
		deliveryNode := costFunction(newOrder)

		if deliveryNode == config.LocalID {
			decLocalOrders <- newOrder

		} else {
			// Am I supposed to send to SharedQ though? I know the best elevator (send to network?)
			decSharedQOrders <- newOrder

			// maybe decNetOrders with new struct, given they're online
			// decNetOrders <- NetworkPacket{ElevID, newOrder}
		}
	}
}

// Confirm listens for an order confirmation from a node
func Confirm() {

}

// RecieveOrder sends an order to internalQ when receiving an order from another node (Or sharedQ???)
func RecieveOrder() {
	// Or directly from sharedQ to internalQ?
}

// AckOrder acknowledges an order to sender node (or all nodes) once light switched on
func AckOrder() {

}

// RedistributeQ redistributes internalQ to other nodes under some motor failure
func RedistributeQ() {

}
