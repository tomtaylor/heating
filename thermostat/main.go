package main

import (
	"flag"
	"github.com/brutella/log"
	"os"
	"os/signal"
	"syscall"
)

var (
	temp     float64 = 20.0
	tempPath string
	device   = ""
	verbose  = false
)

func main() {
	flag.Float64Var(&temp, "temp", 20.0, "Target temperature")
	flag.StringVar(&tempPath, "tempPath", "", "Target temperature path")
	flag.StringVar(&device, "device", "", "Thermometer device")
	flag.BoolVar(&verbose, "verbose", false, "Verbose logging")
	flag.Parse()

	log.Verbose = verbose

	if device == "" {
		log.Fatal("No device specified")
	}

	if tempPath != "" {
		_, err := os.Stat(tempPath)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("[INFO] Starting thermostat with temperature from path:", tempPath)
	} else {
		log.Println("[INFO] Starting thermostat with target temperature:", temp)
	}

	log.Println("[VERB] Using thermometer at", device)

	boiler := NewBoiler()
	thermometer := NewThermometer(device)
	thermostat := NewThermostat(boiler, thermometer, temp)
	api := NewAPI(thermostat)

	if tempPath != "" {
		tempPoller := NewTempPoller(thermostat, tempPath)
		temp = tempPoller.PollTemp()
		go tempPoller.RunLoop()
	}

	go thermostat.RunLoop()
	go api.RunLoop()

	homekitService := NewHomeKitService(thermostat)
	go homekitService.RunLoop()

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	signal := <-ch
	log.Println("[INFO] Received signal", signal, "terminating")

	homekitService.Stop()
	thermostat.Stop()
}
