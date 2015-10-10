all: thermostat boiler_control

.PHONY: thermostat
thermostat:
	mkdir -p build
	GOOS=linux GOARCH=arm GOARM=7 go build -o build/thermostat ./thermostat

.PHONY: boiler_control
boiler_control:
	mkdir -p build
	gcc -Wall -o build/boiler_control boiler_control.c -std=c99 -lwiringPi
