package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Thermometer struct {
	device string
}

func NewThermometer(device string) *Thermometer {
	return &Thermometer{
		device: device,
	}
}

func (t *Thermometer) Temperature() float64 {
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
