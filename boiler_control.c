#include <wiringPi.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define transmitPin 17

bool onData[] = {
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 0, 1, 0, 1, 1, 1, 1, 0, 0
};

bool offData[] = {
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0
};

int dataLength = 25;
int shortDelay = 250;
int longDelay = 500;
int preambleDelay = 1000;

void sendPreamble() {
	for (int i = 0; i < 4; i++) {
		digitalWrite(transmitPin, HIGH);
		delayMicroseconds(preambleDelay);
		digitalWrite(transmitPin, LOW);
		delayMicroseconds(preambleDelay);
	}
}

void sendData(bool data[], int length) {
	for (int i = 0; i < length; i++) {
		bool b = data[i];
		if (b == 1) {
			digitalWrite(transmitPin, HIGH);
			delayMicroseconds(longDelay);
			digitalWrite(transmitPin, LOW);
			delayMicroseconds(shortDelay);
		} else {
			digitalWrite(transmitPin, HIGH);
			delayMicroseconds(shortDelay);
			digitalWrite(transmitPin, LOW);
			delayMicroseconds(longDelay);
		}
	}
	printf("\n");
}

int main(int argc, char **argv) {
	wiringPiSetupGpio();
	pinMode(transmitPin, OUTPUT);

	bool isOn = false;
	if (argc >= 2) {
		if (strcmp(argv[1], "on") == 0) {
			printf("on\n");
			isOn = true;
		} else if (strcmp(argv[1], "off") == 0) {
			printf("off\n");
			isOn = false;
		} else {
			printf("unrecognized command\n");
			exit(1);
		}
	} else {
		printf("no command provided\n");
		exit(1);
	}

	for (int i = 0; i < 2; i++) {
		sendPreamble();
		if (isOn) {
			sendData(onData, dataLength);
		} else {
			sendData(offData, dataLength);
		}
	}
}
