package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	tempRange float64       = 0.5
	interval  time.Duration = 10 * time.Second
)

type Thermostat struct {
	boiler     *Boiler
	device     string
	temps      chan float64
	done       chan bool
	targetTemp float64
}

func NewThermostat(boiler *Boiler, device string, targetTemp float64) *Thermostat {
	return &Thermostat{
		temps:      make(chan float64, 1),
		done:       make(chan bool),
		boiler:     boiler,
		device:     device,
		targetTemp: targetTemp,
	}
}

func (t *Thermostat) RunLoop() {
	go t.boiler.RunLoop()
	t.temps <- t.getTemp()

	for {
		select {
		case temp := <-t.temps:
			log.Println("Temperature is", temp)
			if temp >= t.targetTemp+tempRange {
				if t.boiler.GetCurrentCommand() == true {
					log.Println("Over temperature, turning boiler off")
					t.boiler.SetCurrentCommand(false)
				}
			} else if temp <= t.targetTemp-tempRange {
				if t.boiler.GetCurrentCommand() == false {
					log.Println("Under temperature, turning boiler on")
					t.boiler.SetCurrentCommand(true)
				}
			}
		case <-time.After(interval):
			t.temps <- t.getTemp()
		case <-t.done:
			break
		}
	}

	t.boiler.Stop()
}

func (t *Thermostat) Stop() {
	t.done <- true
}

func (t *Thermostat) SetTemp(temp float64) {
	t.targetTemp = temp
}

func (t *Thermostat) getTemp() float64 {
	path := device + "/w1_slave"

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if lines[0][len(lines[0])-3:] != "YES" {
		log.Fatal("Temperature doesn't seem to be valid")
	}

	line := lines[1]
	l := strings.LastIndexAny(line, "t=")
	celsius, err := strconv.ParseFloat(line[l+1:], 64)
	if err != nil {
		log.Fatal("Could not convert temperature from device")
	}

	return (celsius / 1000.0)
}
