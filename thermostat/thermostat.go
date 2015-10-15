package main

import (
	"github.com/RobinUS2/golang-moving-average"
	"github.com/brutella/log"
	"time"
)

var (
	minUnder            = 0.3
	maxOver             = 0.0
	interval            = 10 * time.Second
	averageTempDuration = 3 * time.Minute
)

type Thermostat struct {
	boiler      *Boiler
	thermometer *Thermometer
	temps       chan float64
	done        chan bool
	targetTemp  float64
	isOn        bool
}

func NewThermostat(boiler *Boiler, thermometer *Thermometer, targetTemp float64) *Thermostat {
	return &Thermostat{
		temps:       make(chan float64, 1),
		done:        make(chan bool),
		boiler:      boiler,
		thermometer: thermometer,
		targetTemp:  targetTemp,
	}
}

func (t *Thermostat) RunLoop() {
	go t.boiler.RunLoop()
	t.temps <- t.Temperature()
	samples := int(averageTempDuration / interval)
	ma := movingaverage.New(samples)

	for {
		select {
		case temp := <-t.temps:
			ma.Add(temp)
			averageTemp := ma.Avg()

			log.Printf("[INFO] Temperature over last 3 minutes is %.2f°C\n", averageTemp)
			if averageTemp >= t.targetTemp+maxOver {
				if t.boiler.GetCurrentCommand() == true {
					log.Println("[INFO] Over temperature, turning boiler off")
					t.boiler.SetCurrentCommand(false)
					t.isOn = false
				}
			} else if averageTemp <= t.targetTemp-minUnder {
				if t.boiler.GetCurrentCommand() == false {
					log.Println("[INFO] Under temperature, turning boiler on")
					t.boiler.SetCurrentCommand(true)
					t.isOn = true
				}
			}
		case <-time.After(interval):
			t.temps <- t.Temperature()
		case <-t.done:
			break
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
