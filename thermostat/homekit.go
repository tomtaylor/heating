package main

import (
	"github.com/brutella/hc/hap"
	"github.com/brutella/hc/model"
	"github.com/brutella/hc/model/accessory"
	"github.com/brutella/log"
	"time"
)

type HomeKitService struct {
	thermostat   *Thermostat
	done         chan bool
	hkThermostat model.Thermostat
	transport    hap.Transport
}

func NewHomeKitService(thermostat *Thermostat) *HomeKitService {
	thermostatInfo := model.Info{
		Name: "Thermostat",
	}

	hkThermostat := accessory.NewThermostat(thermostatInfo, temp, 17, 25, 0.5)
	hkThermostat.SetTargetMode(model.HeatCoolModeHeat)
	hkThermostat.OnTargetTempChange(func(temp float64) {
		log.Println("[INFO] HomeKit requested thermostat to change to", temp)
		thermostat.targetTemp = temp
	})

	hkThermostat.OnTargetModeChange(func(mode model.HeatCoolModeType) {
		log.Println("[WARN] HomeKit requested thermostat to change to", mode, "but this does nothing yet")
	})

	transport, err := hap.NewIPTransport("24282428", hkThermostat.Accessory)
	if err != nil {
		log.Fatal(err)
	}

	t := HomeKitService{
		thermostat:   thermostat,
		done:         make(chan bool),
		hkThermostat: hkThermostat,
		transport:    transport,
	}

	return &t
}

func (hk *HomeKitService) RunLoop() {
	hk.updateState()

	go func() {

		for {
			select {
			case <-time.After(10 * time.Second):
				hk.updateState()
			case <-hk.done:
				break
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
	hk.hkThermostat.SetTemperature(temperature)

	if hk.thermostat.IsOn() {
		hk.hkThermostat.SetMode(model.HeatCoolModeHeat)
	} else {
		hk.hkThermostat.SetMode(model.HeatCoolModeOff)
	}
}
