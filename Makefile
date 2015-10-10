all:
	mkdir -p build
	gcc -Wall -o build/boiler_control boiler_control.c -std=c99 -lwiringPi
	GOOS=linux GOARCH=arm GOARM=7 go build -o build/thermostat ./thermostat
