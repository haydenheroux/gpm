package main

import (
	"flag"
	"fmt"
	"os/exec"
	"time"
)

var (
	all    bool
	names  string
	ids    string
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

	procs := Processes{}
	procs.Add(&Process{Id: 0, Name: "echo", Cmd: *exec.Command("echo", "Hello world!")})
	procs.Add(&Process{Id: 1, Name: "sleep", Cmd: *exec.Command("sleep", "5")})

	go procs.Run()

	for {
		time.Sleep(1 * time.Second)
		for _, proc := range procs.procs {
			fmt.Println(proc)
		}
	}
}
