package main

import (
	"github.com/brutella/log"
	"time"
)

var (
	minUnder           = 0.2
	maxOver            = 0.0
	thermostatInterval = 60 * time.Second
)

type Thermostat struct {
	boiler      *Boiler
	thermometer *Thermometer
	done        chan bool
	targetTemp  float64
	isOn        bool
}

func NewThermostat(boiler *Boiler, thermometer *Thermometer, targetTemp float64) *Thermostat {
	return &Thermostat{
		done:        make(chan bool),
		boiler:      boiler,
		thermometer: thermometer,
		targetTemp:  targetTemp,
	}
}

func (t *Thermostat) RunLoop() {
	go t.boiler.RunLoop()

Loop:
	for {
		select {
		case <-time.After(thermostatInterval):
			temp = t.thermometer.Temperature()

			if temp >= t.targetTemp+maxOver {
				if t.boiler.GetCurrentCommand() == true {
					log.Println("[INFO] Over temperature, turning boiler off")
					t.boiler.SetCurrentCommand(false)
					t.isOn = false
				}
			} else if temp <= t.targetTemp-minUnder {
				if t.boiler.GetCurrentCommand() == false {
					log.Println("[INFO] Under temperature, turning boiler on")
					t.boiler.SetCurrentCommand(true)
					t.isOn = true
				}
			}
		case <-t.done:
			break Loop
		}
	}

	t.boiler.Stop()
}

func (t *Thermostat) Stop() {
	t.done <- true
}

func (t *Thermostat) SetTargetTemperature(temp float64) {
	t.targetTemp = temp
}

func (t *Thermostat) TargetTemperature() float64 {
	return t.targetTemp
}

func (t *Thermostat) Temperature() float64 {
	return t.thermometer.Temperature()
}

func (t *Thermostat) IsOn() bool {
	return t.isOn
}
