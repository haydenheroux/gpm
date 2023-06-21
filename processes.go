package main

import (
	"fmt"
	"os/exec"
	"time"
)

type Process struct {
	Alive     bool
	Id        int
	Name      string
	pid       int
	startTime time.Time
	Version   string
	Cmd       exec.Cmd
}

func (proc *Process) Start(queue chan<- *Process) (err error) {
	defer func() {
		proc.Alive = false
		queue <- proc
	}()

	proc.Alive = true

	proc.Cmd.Start()

	proc.pid = proc.Cmd.Process.Pid
	proc.startTime = time.Now()

	fmt.Println(proc.pid)

	return proc.Cmd.Wait()
}

type Processes struct {
	procs []*Process
}

func (procs *Processes) Add(proc *Process) {
	procs.procs = append(procs.procs, proc)
}

func (procs *Processes) Run() {
	n := len(procs.procs)

	queue := make(chan *Process, n)

	for _, proc := range procs.procs {
		go proc.Start(queue)
	}

	for proc := range queue {
		go proc.Start(queue)
	}
}
