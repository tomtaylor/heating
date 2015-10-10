package main

import (
	"github.com/stianeikeland/go-rpio"
	"log"
	"time"
)

var (
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

type Boiler struct {
	commands       chan bool
	done           chan bool
	pin            rpio.Pin
	currentCommand bool
}

func NewBoiler(pin rpio.Pin) *Boiler {
	return &Boiler{
		commands:       make(chan bool, 1),
		done:           make(chan bool),
		pin:            pin,
		currentCommand: false,
	}
}

func (b *Boiler) SetCurrentCommand(isOn bool) {
	b.commands <- isOn
}

func (b *Boiler) GetCurrentCommand() bool {
	return b.currentCommand
}

func (b *Boiler) Stop() {
	b.done <- true
}

func (b *Boiler) RunLoop() {
	for {
		select {
		case command := <-b.commands:
			b.currentCommand = command
			b.sendCommand(command)
		case <-time.After(30 * time.Second):
			b.commands <- b.currentCommand
		case <-b.done:
			break
		}
	}
}

func (b *Boiler) sendCommand(fire bool) {
	log.Println("Boiler %s", fire)
	//for i := 0; i < 2; i++ {
	//b.sendPreamble()
	//if fire {
	//b.sendData(onData)
	//} else {
	//b.sendData(offData)
	//}
	//}
}

func (b *Boiler) sendPreamble() {
	for i := 0; i < 4; i++ {
		b.pin.High()
		time.Sleep(preambleDelay)
		b.pin.Low()
		time.Sleep(preambleDelay)
	}
}

func (b *Boiler) sendData(data []bool) {
	for _, bit := range data {
		if bit {
			b.pin.High()
			time.Sleep(longDelay)
			b.pin.Low()
			time.Sleep(shortDelay)
		} else {
			b.pin.High()
			time.Sleep(shortDelay)
			b.pin.Low()
			time.Sleep(longDelay)
		}
	}
}
