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
	borderStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
	aliveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	deadStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	cellColor = lipgloss.Color("255")
	cellStyle = lipgloss.NewStyle().Foreground(cellColor)
)

func (procs Processes) View() string {
	columns := []table.Column{
		{Title: "id", Width: 16},
		{Title: "name", Width: 22},
		// TODO This width is super sensitive because of the styling
		{Title: "alive", Width: 18},
		{Title: "pid", Width: 16},
		{Title: "uptime", Width: 20},
	}

	rows := []table.Row{}

	for _, proc := range procs.procs {
		id := cellStyle.Render(fmt.Sprint(proc.Id))
		name := proc.Name
		var alive, pid, uptime string

		if proc.alive {
			alive = aliveStyle.Render("alive")
			pid = fmt.Sprint(proc.pid)
			uptime = time.Since(proc.startTime).String()
		} else {
			alive = deadStyle.Render("dead")
			pid = ""
			uptime = ""
		}

		rows = append(rows, table.Row{id, name, alive, pid, uptime})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(len(procs.procs)),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true).Bold(true)
	s.Selected  = s.Selected.Foreground(cellColor).Bold(false)
	t.SetStyles(s)

	return borderStyle.Render(t.View()) + "\n"
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(50 * time.Millisecond, func (t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
