package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"stock-exchange/internal/common"
)

var (
	stockRegex   = regexp.MustCompile(`^([\w_]+):(\d+)$`)
	processRegex = regexp.MustCompile(`^([\w_]+):\(([^)]*)\):\(([^)]*)\):(\d+)$`)
	optimizeRegex = regexp.MustCompile(`^optimize:\(([^)]+)\)$`)
)

// ParseConfig reads and parses the entire configuration file.
func ParseConfig(filePath string) (*common.Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &common.Config{
		Stocks:    make(common.Resource),
		Processes: []common.Process{},
		Optimize:  []string{},
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		if err := parseLine(line, config); err != nil {
			return nil, err
		}
	}
	
	if len(config.Processes) == 0 {
		return nil, fmt.Errorf("Missing processes")
	}
	if len(config.Optimize) == 0 {
		return nil, fmt.Errorf("Missing optimize directive")
	}

	return config, scanner.Err()
}

// parseLine determines the type of a line and calls the appropriate parser.
func parseLine(line string, config *common.Config) error {
	switch {
	case stockRegex.MatchString(line):
		return parseStock(line, config)
	case processRegex.MatchString(line):
		return parseProcess(line, config)
	case optimizeRegex.MatchString(line):
		return parseOptimize(line, config)
	default:
		return fmt.Errorf("Error while parsing `%s`", line)
	}
}

// parseResource is a helper to parse semi-colon delimited resource strings.
func parseResource(s string) (common.Resource, error) {
	res := make(common.Resource)
	if s == "" {
		return res, nil
	}
	parts := strings.Split(s, ";")
	for _, part := range parts {
		kv := strings.Split(strings.TrimSpace(part), ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid resource format: %s", part)
		}
		qty, err := strconv.Atoi(kv[1])
		if err != nil {
			return nil, fmt.Errorf("invalid quantity for %s", kv[0])
		}
		res[kv[0]] = qty
	}
	return res, nil
}

func parseStock(line string, config *common.Config) error {
	parts := stockRegex.FindStringSubmatch(line)
	qty, _ := strconv.Atoi(parts[2])
	config.Stocks[parts[1]] = qty
	return nil
}

func parseProcess(line string, config *common.Config) error {
	parts := processRegex.FindStringSubmatch(line)
	name := parts[1]
	needsStr := parts[2]
	resultsStr := parts[3]
	cycles, _ := strconv.Atoi(parts[4])

	needs, err := parseResource(needsStr)
	if err != nil {
		return fmt.Errorf("error parsing needs for process %s: %v", name, err)
	}
	results, err := parseResource(resultsStr)
	if err != nil {
		return fmt.Errorf("error parsing results for process %s: %v", name, err)
	}

	config.Processes = append(config.Processes, common.Process{
		Name: name, Needs: needs, Results: results, Cycles: cycles,
	})
	return nil
}

func parseOptimize(line string, config *common.Config) error {
	parts := optimizeRegex.FindStringSubmatch(line)
	config.Optimize = strings.Split(parts[1], ";")
	return nil
}