package workers

import (
	"log"
	"sort"
	"time"

	"github.com/hibiken/asynq"
	"urbangrid.com/constants"
	"urbangrid.com/functions"
	"urbangrid.com/queues"
)

const (
	vehicleClearanceTime = 2  // Time for each vehicle to clear the signal in seconds
	minGreenDuration     = 10 // Minimum duration for green light in seconds
	maxGreenDuration     = 60 // Maximum duration for green light in seconds
)

func TrafficController() {
	// Get the current timer
	currentTimer, err := functions.GetSignalTimer()
	if err != nil {
		log.Printf("Error getting signal timer: %v", err)

		// Set the timer to the current time
		err = functions.SetSignalTimer(time.Now(), "")
		if err != nil {
			log.Printf("Error setting signal timer: %v", err)
		}
		return
	}

	now := time.Now()

	// Check if the current time is past the timer
	timer, err := time.Parse(time.RFC3339, currentTimer.Timer)
	if err != nil {
		log.Printf("Error parsing timer: %v", err)
		return
	}

	if now.After(timer) {
		signals, err := functions.GetSignalData()
		if err != nil {
			log.Printf("Error getting signal data: %v", err)
			return
		}

		if len(signals) == 0 {
			log.Println("No signals found.")
			return
		}

		// expect the last active signal
		lastActiveSignal := currentTimer.LastActiveSignal

		// Get the last active signal
		for i, signal := range signals {
			if signal.Id == lastActiveSignal {
				// Move the last active signal to the front of the slice
				signals = append(signals[:i], signals[i+1:]...)
				break
			}
		}

		// Sort signals by vehicle count in descending order
		sort.Slice(signals, func(i, j int) bool {
			return signals[i].Count > signals[j].Count
		})

		// Get the signal with the highest count
		signal := signals[0]

		// Calculate the predicted clearance time
		clearanceTime := predictClearanceTime(signal.Count)

		// Set the new timer for the next signal change
		nextChangeTime := now.Add(time.Duration(clearanceTime) * time.Second)

		// Turn on the signal with the highest count and set duration
		err = functions.TurnOnSignal(signal)
		if err != nil {
			log.Printf("Error turning on signal %s: %v", signal.Id, err)
		}
		err = functions.SetSignalTimer(nextChangeTime, signal.Id)

		if err != nil {
			log.Printf("Error setting signal timer: %v", err)
		}

		defer queues.EnqueueTask(asynq.NewTask(constants.SIGNAL_CHANGE, []byte("")))

		// Turn off all other signals
		for _, d := range signals[1:] {
			err = functions.TurnOffSignal(d)
			if err != nil {
				log.Printf("Error turning off signal %s: %v", d.Id, err)
			}
		}
	}
}

// PredictClearanceTime calculates the estimated time to clear the signal based on vehicle count
func predictClearanceTime(vehicleCount int64) int64 {
	clearanceTime := vehicleCount * vehicleClearanceTime
	if clearanceTime < minGreenDuration {
		clearanceTime = minGreenDuration
	}
	if clearanceTime > maxGreenDuration {
		clearanceTime = maxGreenDuration
	}
	return clearanceTime
}
