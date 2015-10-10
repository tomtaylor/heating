package main

import (
	"flag"
	"github.com/stianeikeland/go-rpio"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	temp     float64 = 20.0
	tempPath string
	device   = ""
	pin      = rpio.Pin(17)
)

func main() {
	flag.Float64Var(&temp, "temp", 20.0, "Target temperature")
	flag.StringVar(&tempPath, "tempPath", "", "Target temperature path")
	flag.StringVar(&device, "device", "", "Thermometer device")
	flag.Parse()

	if device == "" {
		log.Fatal("No device specified")
	}

	if tempPath != "" {
		_, err := os.Stat(tempPath)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Starting thermostat with temperature from path:", tempPath)
	} else {
		log.Println("Starting thermostat with target temperature:", temp)
	}

	if err := rpio.Open(); err != nil {
		log.Fatal(err)
	}

	defer rpio.Close()
	pin.Output()

	log.Println("Using thermometer at", device)

	boiler := NewBoiler(pin)
	thermostat := NewThermostat(boiler, device, temp)

	if tempPath != "" {
		tempPoller := NewTempPoller(thermostat, tempPath)
		go tempPoller.RunLoop()
	}

	go thermostat.RunLoop()

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	signal := <-ch
	log.Println("Received signal", signal, "terminating")

	thermostat.Stop()
}
