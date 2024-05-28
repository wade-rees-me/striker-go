package simulators

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/wade-rees-me/striker-go/constants"
	"github.com/wade-rees-me/striker-go/database"
	"github.com/wade-rees-me/striker-go/logger"
	"github.com/wade-rees-me/striker-go/queues"
)

func SimulatorRunOnce() {
	if err := SimulatorProcess(NewSimulation()); err != nil {
		logger.Log.Error(fmt.Sprintf("Simulation failed: %s", err))
	}
}

func SimulatorRunQueue() {
	logger.Log.Info(fmt.Sprintf("Starting the queue: %s:", constants.QueueName))
	for {
		logger.Log.Info(fmt.Sprintf("Reading from the queue: %s:", constants.QueueName))
		msgResult, err := queues.Receive(constants.QueueName, constants.QueueTimeout)
		if err != nil {
			logger.Log.Warning(fmt.Sprintf("Cannot load the request message from queue: %s: %s", constants.QueueName, err))
			time.Sleep(constants.QueueSleepTime * time.Second)
			continue
		}
		if len(msgResult.Messages) == 0 || *msgResult.Messages[0].ReceiptHandle == "" {
			time.Sleep(constants.QueueSleepTime * time.Second)
			continue
		}

		parameters := new(SimulationParameters)
		request := *msgResult.Messages[0].Body
		logger.Log.Debug(fmt.Sprintf("Request message loaded from queue: %s: %s", constants.QueueName, request))
		if err = json.Unmarshal([]byte(request), parameters); err != nil {
			logger.Log.Warning(fmt.Sprintf("Cannot parse the request message from queue: %s: %s", request, err))
			continue
		}
		if parameters.Target != constants.StrikerWhoAmI {
			continue
		}

		var wg sync.WaitGroup
		done := make(chan bool)
		sim := NewSimulation()
		sim.Parameters = *parameters

		wg.Add(2)
		go SimulatorProcessRun(done, &wg, sim)
		go monitorQueue(done, &wg, *msgResult.Messages[0].ReceiptHandle)

		logger.Log.Debug(fmt.Sprintf("Waiting for goroutines to finish..."))
		wg.Wait()
		logger.Log.Debug(fmt.Sprintf("Goroutines have finished..."))

		if err := queues.Delete(constants.QueueName, *msgResult.Messages[0].ReceiptHandle); err != nil {
			logger.Log.Error(fmt.Sprintf("Request message not deleted from queue: %s: %s", constants.QueueName, err))
		} else {
			logger.Log.Info(fmt.Sprintf("Request delete from the queue: %s:", constants.QueueName))
		}
	}
}

func SimulatorProcessRun(done chan bool, wg *sync.WaitGroup, s *Simulation) error {
	defer wg.Done()

	logger.Log.Debug(fmt.Sprintf("Starting Queue process..."))
	if err := SimulatorProcess(s); err != nil {
		logger.Log.Error(fmt.Sprintf("Queue Simulation failed: %s", err))
		return err
	}
	done <- true
	logger.Log.Debug(fmt.Sprintf("Finished Queue process..."))
	return nil
}

func monitorQueue(done chan bool, wg *sync.WaitGroup, receiptHandle string) {
	defer wg.Done()

	logger.Log.Debug(fmt.Sprintf("Starting the monitor process..."))
	for {
		select {
		case <-done:
			logger.Log.Debug(fmt.Sprintf("Stopping the monitor process..."))
			return
		case <-time.After(constants.QueueSleepTime / 2 * time.Second):
			logger.Log.Debug(fmt.Sprintf("Updating the Queue timeout..."))
			queues.Update(constants.QueueName, constants.QueueSleepTime, receiptHandle)
		}
	}
}

func SimulatorProcess(s *Simulation) error {
	s.PrintSimulation()
	if err := database.ProcessingInsert(s.Parameters.Target, s.Parameters.Guid, s.Hostname, s.Parameters.Timestamp, s.getParameters()); err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to insert into Processing table: %s", err))
		return err
	}
	s.RunSimulation()
	if err := database.SimulationInsert(s.Parameters.Target, s.Parameters.Guid, s.Hostname, s.Parameters.Strategy, s.Parameters.Rules, s.Parameters.Decks, s.Parameters.Timestamp, s.Duration, s.getReport()); err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to insert into Simulation table: %s", err))
		return err
	}
	if err := database.ProcessingDelete(s.Parameters.Target, s.Parameters.Guid); err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to delete from Processing table: %s", err))
		return err
	}
	return nil
}
