package simulator

import (
	"fmt"
	"stock-exchange/internal/common"
)

// ActiveProcess tracks a process that is currently running during simulation.
type ActiveProcess struct {
	Process  common.Process
	EndCycle int
}

// SimulateSchedule runs a step-by-step simulation of a given schedule against a config.
// It verifies that all resource constraints are met at each step.
func SimulateSchedule(config *common.Config, schedule []common.ScheduleEntry) (int, error) {
	// Create a copy of the initial stocks to avoid modifying the original config.
	stocks := make(common.Resource)
	for k, v := range config.Stocks {
		stocks[k] = v
	}

	activeProcesses := []ActiveProcess{}
	currentCycle := 0
	scheduleIndex := 0

	// The simulation continues as long as there are scheduled processes to start
	// or active processes still running.
	for scheduleIndex < len(schedule) || len(activeProcesses) > 0 {
		// Determine the time of the next event. An event is either a process starting
		// (from the schedule) or a process finishing (from the active list).
		nextEventCycle := -1

		// Check the start time of the next scheduled process.
		if scheduleIndex < len(schedule) {
			nextEventCycle = schedule[scheduleIndex].Cycle
		}

		// Check the end time of the soonest-to-finish active process.
		for _, ap := range activeProcesses {
			if nextEventCycle == -1 || ap.EndCycle < nextEventCycle {
				nextEventCycle = ap.EndCycle
			}
		}

		// If there are no more events, the simulation is over.
		if nextEventCycle == -1 {
			break
		}

		// Jump the simulation time to the next event.
		currentCycle = nextEventCycle

		// --- Process Events at the Current Cycle ---

		// 1. Complete any processes that have finished by this cycle.
		// We filter the list, keeping only the ones that are still active.
		var stillActive []ActiveProcess
		for _, ap := range activeProcesses {
			if ap.EndCycle <= currentCycle {
				// Process finished. Add its results back to the stocks.
				for name, qty := range ap.Process.Results {
					stocks[name] += qty
				}
			} else {
				stillActive = append(stillActive, ap)
			}
		}
		activeProcesses = stillActive

		// 2. Start any new processes that are scheduled for this exact cycle.
		for scheduleIndex < len(schedule) && schedule[scheduleIndex].Cycle == currentCycle {
			entry := schedule[scheduleIndex]
			fmt.Printf("Evaluating: %d:%s\n", entry.Cycle, entry.Process.Name)

			// AUDIT POINT: Check if resources are sufficient before starting.
			for name, neededQty := range entry.Process.Needs {
				if stocks[name] < neededQty {
					// This is the error condition required by the audit.
					return currentCycle, fmt.Errorf("at %d:%s stock insufficient (need %d of %s, have %d)",
						entry.Cycle, entry.Process.Name, neededQty, name, stocks[name])
				}
			}

			// If checks pass, consume the resources.
			for name, neededQty := range entry.Process.Needs {
				stocks[name] -= neededQty
			}

			// Add the process to the active list.
			activeProcesses = append(activeProcesses, ActiveProcess{
				Process:  entry.Process,
				EndCycle: currentCycle + entry.Process.Cycles,
			})
			scheduleIndex++
		}
	}

	// Simulation completed without errors.
	return currentCycle, nil
}
