package simulators

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/wade-rees-me/go-blackjack/cmd/striker/constants"
	"github.com/wade-rees-me/go-blackjack/cmd/striker/database"
	"github.com/wade-rees-me/go-blackjack/cmd/striker/queues"
)

func SimulatorRunOnce() {
	sim := NewSimulation()
	sim.PrintSimulation()
	database.ProcessingInsert(sim.Parameters.Target, sim.Parameters.Guid, sim.Hostname, sim.Parameters.Timestamp, sim.getParameters())
	sim.RunSimulation()
	database.SimulationInsert(sim.Parameters.Target, sim.Parameters.Guid, sim.Hostname, sim.Parameters.Strategy, sim.Parameters.Rules, sim.Parameters.Decks, sim.Parameters.Timestamp, sim.getReport())
	database.ProcessingDelete(sim.Parameters.Guid)
}

func SimulatorRunQueue() {
	fmt.Println("Starting queue")
	for {
		msgResult, err := queues.Receive(constants.QueueName, constants.QueueTimeout)
		if err != nil {
			fmt.Println("Warning: cannot load the request message: ", err)
			continue
		}

		fmt.Println("  Reading queue: ", fmt.Sprint(len(msgResult.Messages)))
		if len(msgResult.Messages) == 0 {
			fmt.Println("  Sleeping")
			time.Sleep(constants.QueueSleepTime * time.Second)
			continue
		}

		sim := NewSimulation()
		request := *msgResult.Messages[0].Body
		fmt.Println("Request: " + request)

		err = json.Unmarshal([]byte(request), &sim.Parameters)
		if err != nil {
			fmt.Println("Warning: cannot parse the request message: ", err)
			continue
		}

		if sim.Parameters.Target != constants.StrikerWhoAmI {
			continue
		}
		database.ProcessingInsert(sim.Parameters.Target, sim.Parameters.Guid, sim.Hostname, sim.Parameters.Timestamp, *msgResult.Messages[0].Body)
		if *msgResult.Messages[0].ReceiptHandle != "" {
			fmt.Println("Delete from queue: " + *msgResult.Messages[0].ReceiptHandle)
			err = queues.Delete(constants.QueueName, *msgResult.Messages[0].ReceiptHandle)
			if err != nil {
				fmt.Println("Warning: cannot delete the request message: ", err)
				database.ProcessingDelete(sim.Parameters.Guid)
				continue
			}
		}
		sim.RunSimulation() // Add results to the processed table in the database
		database.SimulationInsert(sim.Parameters.Target, sim.Parameters.Guid, sim.Hostname, sim.Parameters.Strategy, sim.Parameters.Rules, sim.Parameters.Decks, sim.Parameters.Timestamp, sim.getReport())
		database.ProcessingDelete(sim.Parameters.Guid)
	}
}
