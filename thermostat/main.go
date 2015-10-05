package main

import (
	"github.com/stianeikeland/go-rpio"
	"log"
	"math/rand"
	"time"
)

var (
	pin           = rpio.Pin(4)
	shortDelay    = 250 * time.Microsecond
	longDelay     = 500 * time.Microsecond
	preambleDelay = 1000 * time.Microsecond
	onData        = []bool{
		false, false, false, false, false, false, false, false, false, false, false, false, false, false, false,
		true, false, true, false, true, true, true, true, false, false,
	}
	offData = []bool{
		false, false, false, false, false, false, false, false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false, false, false,
	}
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatal(err)
	}

	defer rpio.Close()
	pin.Output()

	commandChan := make(chan bool, 1)

	// TODO: Replace this with a thermometer driven loop rather than random off/on
	go func() {
		current := false
		for {
			current = !current
			commandChan <- current
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		}
	}()

	sendCommandsInLoop(commandChan)
}

func sendCommand(fire bool) {
	for i := 0; i < 2; i++ {
		sendPreamble()
		if fire {
			sendData(onData)
		} else {
			sendData(offData)
		}
	}
}

func sendPreamble() {
	for i := 0; i < 4; i++ {
		pin.High()
		time.Sleep(preambleDelay)
		pin.Low()
		time.Sleep(preambleDelay)
	}
}

func sendData(data []bool) {
	for _, bit := range data {
		if bit {
			pin.High()
			time.Sleep(longDelay)
			pin.Low()
			time.Sleep(shortDelay)
		} else {
			pin.High()
			time.Sleep(shortDelay)
			pin.Low()
			time.Sleep(longDelay)
		}
	}
}

func sendCommandsInLoop(c chan bool) {
	currentCommand := false
	for {
		select {
		case command := <-c:
			currentCommand = command
			sendCommand(command)
		case <-time.After(30 * time.Second):
			c <- currentCommand
		}
	}
}
