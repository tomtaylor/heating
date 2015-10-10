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
	pin                = rpio.Pin(17)
	device             = "/sys/bus/w1/devices/28-031572e41cff"
	targetTemp float64 = 20.0
)

func main() {
	flag.Float64Var(&targetTemp, "temperature", 20.0, "Target temperature")
	flag.Parse()
	runServer()
}

func runServer() {
	if err := rpio.Open(); err != nil {
		log.Fatal(err)
	}

	defer rpio.Close()
	pin.Output()

	log.Println("Starting thermostat with target temperature:", targetTemp)

	boiler := NewBoiler(pin)
	thermostat := NewThermostat(boiler, device, targetTemp)
	go thermostat.RunLoop()

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	signal := <-ch
	log.Println("Received signal", signal, "terminating")

	thermostat.Stop()
}
