package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"stock-exchange/internal/common"
	"stock-exchange/internal/common/parser"
	"stock-exchange/internal/common/scheduler"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: go run . <file> <waiting_time_seconds>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	timeoutSec, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid waiting time: %v\n", err)
		os.Exit(1)
	}
	timeout := time.Duration(timeoutSec * float64(time.Second))

	config, err := parser.ParseConfig(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "Exiting...")
		os.Exit(1)
	}

	schedule, finalStocks := scheduler.Run(config, timeout)

	// Print results to stdout
	fmt.Println("Main Processes:")
	for _, entry := range schedule {
		fmt.Printf(" %d:%s\n", entry.Cycle, entry.Process.Name)
	}

	fmt.Println("Stock:")
	printStocks(finalStocks, os.Stdout)

	// Write schedule to log file
	logFilePath := strings.TrimSuffix(filePath, filepath.Ext(filePath)) + ".log"
	writeLogFile(logFilePath, schedule)
}

func writeLogFile(path string, schedule []common.ScheduleEntry) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log file: %v\n", err)
		return
	}
	defer f.Close()

	for _, entry := range schedule {
		fmt.Fprintf(f, "%d:%s\n", entry.Cycle, entry.Process.Name)
	}
}

func printStocks(stocks common.Resource, writer *os.File) {
	// Sort stock names for consistent output
	keys := make([]string, 0, len(stocks))
	for k := range stocks {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(writer, " %s => %d\n", k, stocks[k])
	}
}
