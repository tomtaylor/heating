package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var (
	transmitInterval = time.Second * 15
)

type Boiler struct {
	commands       chan bool
	done           chan bool
	currentCommand bool
}

func NewBoiler() *Boiler {
	return &Boiler{
		commands:       make(chan bool, 1),
		done:           make(chan bool),
		currentCommand: false,
	}
}

func (b *Boiler) SetCurrentCommand(isOn bool) {
	if isOn != b.currentCommand {
		b.commands <- isOn
	}
}

func (b *Boiler) GetCurrentCommand() bool {
	return b.currentCommand
}

func (b *Boiler) Stop() {
	b.done <- true
}

func (b *Boiler) RunLoop() {
	b.commands <- b.currentCommand

	for {
		select {
		case command := <-b.commands:
			b.currentCommand = command
			b.sendCommand(command)
		case <-time.After(transmitInterval):
			b.commands <- b.currentCommand
		case <-b.done:
			break
		}
	}
}

func (b *Boiler) sendCommand(fire bool) {
	var arg string
	if fire {
		arg = "on"
	} else {
		arg = "off"
	}

	wd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	bin := filepath.Join(wd, "boiler_control")

	args := []string{arg}
	cmd := exec.Command(bin, args...)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
