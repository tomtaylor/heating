package main

import (
	"bufio"
	"github.com/RobinUS2/golang-moving-average"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	thermometerInterval = 10 * time.Second
	averageTempDuration = 3 * time.Minute
)

type Thermometer struct {
	device      string
	temps       chan float64
	temperature float64
	done        chan bool
}

func NewThermometer(device string) *Thermometer {
	return &Thermometer{
		device: device,
		temps:  make(chan float64, 1),
		done:   make(chan bool),
	}
}

func (t *Thermometer) RunLoop() {
	t.temps <- t.getTemperature()
	samples := int(averageTempDuration / thermometerInterval)
	ma := movingaverage.New(samples)

	for {
		select {
		case temp := <-t.temps:
			ma.Add(temp)
			t.temperature = ma.Avg()
			log.Printf("[INFO] Temperature over last 3 minutes is %.2f°C\n", t.temperature)
		case <-time.After(thermometerInterval):
			t.temps <- t.getTemperature()
		case <-t.done:
			break
		}
	}
}

func (t *Thermometer) Stop() {
	t.done <- true
}

func (t *Thermometer) Temperature() float64 {
	return t.temperature
}

func (t *Thermometer) getTemperature() float64 {
	path := t.device + "/w1_slave"

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

	statusLine := lines[0]
	statusLineLen := len(statusLine)

	if statusLineLen < 3 || statusLine[statusLineLen-3:] != "YES" {
		log.Fatal("Temperature doesn't seem to be valid")
	}

	dataLine := lines[1]
	l := strings.LastIndexAny(dataLine, "t=")
	celsius, err := strconv.ParseFloat(dataLine[l+1:], 64)
	if err != nil {
		log.Fatal("Could not convert temperature from device")
	}

	return (celsius / 1000.0)
}
