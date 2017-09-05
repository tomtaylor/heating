package main

import (
	"bufio"
	"github.com/brutella/log"
	"os"
	"strconv"
	"time"
)

var (
	pollerInterval = 5 * time.Second
)

type TempPoller struct {
	path       string
	done       chan bool
	lastTemp   float64
	thermostat *Thermostat
}

func NewTempPoller(thermostat *Thermostat, path string) *TempPoller {
	return &TempPoller{
		path:       path,
		thermostat: thermostat,
		done:       make(chan bool),
	}
}

func (tp *TempPoller) RunLoop() {
	tp.updateTemp()

Loop:
	for {
		select {
		case <-time.After(pollerInterval):
			tp.updateTemp()
		case <-tp.done:
			break Loop
		}
	}
}

func (tp *TempPoller) Stop() {
	tp.done <- true
}

func (tp *TempPoller) updateTemp() {
	temp := tp.PollTemp()
	if temp != tp.lastTemp {
		log.Println("[INFO] Setting thermostat target temperature to", temp)
		tp.thermostat.targetTemp = temp
		tp.lastTemp = temp
	}
}

func (tp *TempPoller) PollTemp() float64 {
	file, err := os.Open(tp.path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	f, err := strconv.ParseFloat(lines[0], 64)
	if err != nil {
		log.Fatal(err)
	}

	return f
}
