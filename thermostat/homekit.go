package main

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/log"
	"time"
)

var (
	defaultOnTemp  = 20.0
	defaultOffTemp = 17.0
)

type HomeKitService struct {
	thermostat *Thermostat
	done       chan bool
	accessory  accessory.Thermostat
	transport  hc.Transport
}

func NewHomeKitService(thermostat *Thermostat) *HomeKitService {
	info := accessory.Info{
		Name: "Thermostat",
	}

	acc := accessory.NewThermostat(info, temp, 17, 25, 0.5)
	acc.Thermostat.TargetHeatingCoolingState.UpdateValue(characteristic.TargetHeatingCoolingStateHeat)
	acc.Thermostat.TargetTemperature.OnValueRemoteUpdate(func(temp float64) {
		log.Println("[INFO] HomeKit requested thermostat temperature change to", temp)
		thermostat.targetTemp = temp
	})

	acc.Thermostat.TargetHeatingCoolingState.OnValueRemoteUpdate(func(mode int) {
		log.Println("[INFO] HomeKit requested thermostat mode change to", mode)

		switch mode {
		case characteristic.TargetHeatingCoolingStateHeat:
			log.Println("[INFO] HomeKit setting thermostat to default on temp of", defaultOnTemp)
			thermostat.targetTemp = defaultOnTemp
		case characteristic.TargetHeatingCoolingStateOff:
			log.Println("[INFO] HomeKit setting thermostat to default off temp of", defaultOffTemp)
			thermostat.targetTemp = defaultOffTemp
		case characteristic.TargetHeatingCoolingStateAuto, characteristic.TargetHeatingCoolingStateCool:
			acc.Thermostat.TargetHeatingCoolingState.UpdateValue(characteristic.TargetHeatingCoolingStateHeat)
		}

	})

	config := hc.Config{Pin: "24282428"}
	transport, err := hc.NewIPTransport(config, acc.Accessory)
	if err != nil {
		log.Fatal(err)
	}

	t := HomeKitService{
		thermostat: thermostat,
		done:       make(chan bool),
		accessory:  *acc,
		transport:  transport,
	}

	return &t
}

func (hk *HomeKitService) RunLoop() {
	hk.updateState()

	go func() {
	Loop:
		for {
			select {
			case <-time.After(10 * time.Second):
				hk.updateState()
			case <-hk.done:
				break Loop
			}
		}

		hk.transport.Stop()
	}()

	hk.transport.Start()
}

func (hk *HomeKitService) Stop() {
	hk.done <- true
}

func (hk *HomeKitService) updateState() {
	temperature := hk.thermostat.Temperature()
	hk.accessory.Thermostat.CurrentTemperature.UpdateValue(temperature)

	targetTemperature := hk.thermostat.TargetTemperature()
	hk.accessory.Thermostat.TargetTemperature.UpdateValue(targetTemperature)

	if hk.thermostat.IsOn() {
		hk.accessory.Thermostat.CurrentHeatingCoolingState.UpdateValue(characteristic.CurrentHeatingCoolingStateHeat)
	} else {
		hk.accessory.Thermostat.CurrentHeatingCoolingState.UpdateValue(characteristic.CurrentHeatingCoolingStateOff)
	}
}
