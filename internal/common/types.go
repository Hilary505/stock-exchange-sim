package common


// Resource represents a map of item names to their quantities.
// Used for both stocks and the needs/results of processes.
type Resource map[string]int

// Process defines a single task with its requirements and outcomes.
type Process struct {
	Name    string
	Needs   Resource
	Results Resource
	Cycles  int
}

// Config holds all the parsed information from a configuration file.
type Config struct {
	Stocks    Resource
	Processes []Process
	Optimize  []string
}

// ScheduleEntry represents a single scheduled task.
type ScheduleEntry struct {
	Cycle int
	Process
}