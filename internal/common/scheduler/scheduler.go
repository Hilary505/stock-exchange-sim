package scheduler

import (
	"fmt"
	"math"
	"os"
	"time"

	"stock-exchange/internal/common"
)

// ActiveProcess tracks a process that is currently running.
type ActiveProcess struct {
	Process  common.Process
	EndCycle int
}

// Run performs the optimization and scheduling.
func Run(config *common.Config, timeout time.Duration) ([]common.ScheduleEntry, common.Resource) {
	startTime := time.Now()
	stocks := config.Stocks
	schedule := []common.ScheduleEntry{}
	activeProcesses := []ActiveProcess{}
	currentCycle := 0

	for {
		// Check for timeout, which is the shutdown condition for infinite loops
		if time.Since(startTime) >= timeout {
			break
		}

		// Step 1: Complete any finished processes
		newlyFinished := true
		for newlyFinished {
			newlyFinished = false
			nextActive := []ActiveProcess{}
			for _, ap := range activeProcesses {
				if ap.EndCycle <= currentCycle {
					for name, qty := range ap.Process.Results {
						stocks[name] += qty
					}
					newlyFinished = true
				} else {
					nextActive = append(nextActive, ap)
				}
			}
			activeProcesses = nextActive
		}

		// Step 2: Try to start new processes in the current cycle (greedy approach)
		startedNewProcess := true
		for startedNewProcess {
			startedNewProcess = false
			bestProcess := findBestDoableProcess(config.Processes, stocks)
			if bestProcess != nil {
				// Consume resources
				for name, qty := range bestProcess.Needs {
					stocks[name] -= qty
				}
				// Add to schedule and active list
				schedule = append(schedule, common.ScheduleEntry{Cycle: currentCycle, Process: *bestProcess})
				activeProcesses = append(activeProcesses, ActiveProcess{Process: *bestProcess, EndCycle: currentCycle + bestProcess.Cycles})
				startedNewProcess = true
			}
		}

		// Step 3: Advance time
		if len(activeProcesses) == 0 {
			// No running processes and nothing new to start means we are done.
			break
		}

		// Find the soonest end cycle to jump time forward
		minEndCycle := math.MaxInt32
		for _, ap := range activeProcesses {
			if ap.EndCycle < minEndCycle {
				minEndCycle = ap.EndCycle
			}
		}

		if minEndCycle <= currentCycle {
			// This can happen if a process takes 0 cycles. Advance by 1 to avoid infinite loop.
			currentCycle++
		} else {
			currentCycle = minEndCycle
		}
	}

	fmt.Fprintf(os.Stderr, "No more process doable at cycle %d\n", currentCycle)
	return schedule, stocks
}

// findBestDoableProcess finds the first process that can be started.
// A more complex heuristic could be implemented here to better satisfy the "optimize" directive.
func findBestDoableProcess(processes []common.Process, stocks common.Resource) *common.Process {
	for i, p := range processes {
		canDo := true
		for name, neededQty := range p.Needs {
			if stocks[name] < neededQty {
				canDo = false
				break
			}
		}
		if canDo {
			return &processes[i]
		}
	}
	return nil
}
