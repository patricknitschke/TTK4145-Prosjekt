package decision

import (
	"fmt"

	"../config"
	"../elevator"
)

// Calculates the cost of an order based on the Elevator positions and state, returns optimal Node
func costFunction(newOrder config.Order) int {

	// NNodes is 0 now, remember to change in config
	var elevatorList [config.NNodes]elevator.Elevator
	var elevatorOnline [config.NNodes]bool

	// TODO Get elevatorList and onlineList somehow
	elevatorList[config.LocalID] = elevator.GetState() // My state
	elevatorOnline[config.LocalID] = true
	// TODO ///////////////////////////////////////

	if newOrder.OrderT == config.Cab {
		return config.LocalID
	}
	var CostArray [config.NNodes]int
	for elev := 0; elev < config.NNodes; elev++ {
		cost := newOrder.Floor - elevatorList[elev].CurrentFloor
		if cost == 0 && elevatorOnline[elev] && !elevatorList[elev].MotorFailure {
			return elev
		}
		if cost < 0 {
			cost = -cost
			if elevatorList[elev].Dir == elevator.Up {
				cost += 3
			}
		} else if cost > 0 {
			if elevatorList[elev].Dir == elevator.Down {
				cost += 3
			}
		}
		if elevatorList[elev].DoorState == true {
			cost++
		}
		CostArray[elev] = cost
	}
	maxCost := 1000
	var bestElev int
	for elev := 0; elev < config.NNodes; elev++ {
		if CostArray[elev] < maxCost && elevatorOnline[elev] && !elevatorList[elev].MotorFailure {
			bestElev = elev
			maxCost = CostArray[elev]
		}
	}
	fmt.Print("Cost to Serve: ")
	fmt.Println(CostArray)
	return bestElev
}
