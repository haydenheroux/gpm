package main

import (
	"flag"
	"fmt"
)

var (
	all bool
	names string
	ids string
	config string
)

func init() {
	flag.BoolVar(&all, "all", false, "Start and monitor all configured processes.")
	flag.StringVar(&names, "name", "", "Start only the process with this name.")
	flag.StringVar(&ids, "id", "", "Start only the process with this identifier.")
	flag.StringVar(&config, "config", "", "Load the configuration from this path.")
}

func main() {
	flag.Parse()

	allProcesses := loadAll()

	// Set containing processes to be started
	var processes map[interface{}]struct{}

	if all {
		processes = allProcesses
	} else {
		processes = make(map[interface{}]struct{}, 0)
	}

	// Add processes to set if matching name
	if names != "" && !all {
		processesMatchingNames := allProcesses
		for process := range processesMatchingNames {
			processes[process] = struct{}{}
		}
	}

	// Add processes to set if matching identifier
	if ids != "" && !all {
		processesMatchingIds := allProcesses
		for process := range processesMatchingIds {
			processes[process] = struct{}{}
		}
	}

	// Start all processes in the set
	for process := range processes {
		go start(process)
	}

	// TODO Render GUI
	fmt.Println("Monitoring processes...")
}

func loadAll() map[interface{}]struct{} {
	return make(map[interface{}]struct{}, 0)
}

func start(process interface{}) {
	fmt.Println("Started a process")
}
