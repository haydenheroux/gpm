package main

import (
	"flag"
	"fmt"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	return tickCmd()
}

func (procs Processes) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
		case tea.KeyMsg:
					return procs, tea.Quit

		case tickMsg:
			return procs, tickCmd()

		default:
			return procs, nil
	}
}

var (
	border = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
)

func (procs Processes) View() string {
	columns := []table.Column{
		{Title: "id", Width: 6},
		{Title: "name", Width: 22},
		{Title: "version", Width: 16},
		{Title: "pid", Width: 20},
		{Title: "uptime", Width: 20},
		{Title: "alive", Width: 8},
	}

	rows := []table.Row{}

	for _, proc := range procs.procs {
		id := fmt.Sprint(proc.Id)
		name := proc.Name
		version := proc.Version
		pid := fmt.Sprint(proc.pid)
		uptime := time.Since(proc.startTime).String()
		alive := fmt.Sprint(proc.Alive)
		rows = append(rows, table.Row{id, name, version, pid, uptime, alive})
	}

	table := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	return border.Render(table.View()) + "\n"
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(50 * time.Millisecond, func (t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
