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

	if lines[0][len(lines[0])-3:] != "YES" {
		log.Fatal("Temperature doesn't seem to be valid")
	}

	line := lines[1]
	l := strings.LastIndexAny(line, "t=")
	celsius, err := strconv.ParseFloat(line[l+1:], 64)
	if err != nil {
		log.Fatal("Could not convert temperature from device")
	}

	return (celsius / 1000.0)
}
