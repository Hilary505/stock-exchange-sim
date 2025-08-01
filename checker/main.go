package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"stock-exchange/internal/common"
	"stock-exchange/internal/common/parser"
	simulator "stock-exchange/internal/common/stimulator" 
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: go run ./checker <config_file> <log_file>")
		os.Exit(1)
	}

	configFile := os.Args[1]
	logFile := os.Args[2]

	config, err := parser.ParseConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing config: %v\n", err)
		os.Exit(1)
	}

	schedule, err := readLogFile(logFile, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading log file: %v\n", err)
		os.Exit(1)
	}

	// Call the dedicated simulator function.
	finalCycle, err := simulator.SimulateSchedule(config, schedule)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detected\n%v\nExiting...\n", err)
		os.Exit(1)
	}

	fmt.Printf("Trace completed, no error detected.\n")
	fmt.Printf("Final cycle: %d\n", finalCycle) 
}

// readLogFile reads the schedule from a .log file.
func readLogFile(path string, config *common.Config) ([]common.ScheduleEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var schedule []common.ScheduleEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid log format: %s", line)
		}
		cycle, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid cycle in log: %s", parts[0])
		}
		procName := parts[1]

		var foundProc *common.Process
		for i := range config.Processes {
			if config.Processes[i].Name == procName {
				foundProc = &config.Processes[i]
				break
			}
		}
		if foundProc == nil {
			return nil, fmt.Errorf("process '%s' from log not found in config", procName)
		}
		schedule = append(schedule, common.ScheduleEntry{Process: *foundProc, Cycle: cycle})
	}
	return schedule, scanner.Err()
}
