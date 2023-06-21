package main

import (
	"flag"
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
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

	p := tea.NewProgram(procs)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

func (procs Processes) Init() tea.Cmd {
	return nil
}

func (procs Processes) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return procs, tea.Quit

		}

	}

	return procs, nil
}

func (procs Processes) View() string {
	s := ""

	for _, proc := range procs.procs {
		s += fmt.Sprintf("[%d] %s (%s) | pid: %d %v %s\n", proc.Id, proc.Name, proc.Version, proc.pid, proc.Alive, proc.startTime)
	}

	return s
}
