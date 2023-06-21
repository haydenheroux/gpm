package main

import (
	"os/exec"
	"time"
)

type Process struct {
	alive     bool
	Id        int
	Name      string
	pid       int
	startTime time.Time
	Cmd       exec.Cmd
}

func (proc *Process) Start(queue chan<- *Process) (err error) {
	defer func() {
		proc.alive = false
		queue <- proc
	}()

	proc.alive = true

	cmd := proc.Cmd

	cmd.Start()

	proc.pid = cmd.Process.Pid
	proc.startTime = time.Now()

	return cmd.Wait()
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
