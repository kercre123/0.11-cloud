#!/bin/bash

if [[ $EUID -ne 0 ]]; then
    echo "This script must be run as root. sudo ./start.sh"
    exit 1
fi

if [[ ! -f source.sh ]]; then
	echo "You must create a source.sh file with the env vars in the readme. They all need an export statement"
	exit 1
fi

source source.sh

if [[ ! -f main ]]; then
	echo "building the program... (if there are errors, make sure you have wire-pod installed. it will put vosk in the right place)"
	CGO_LDFLAGS=-L/root/.vosk/libvosk CGO_CFLAGS=-I/root/.vosk/libvosk go build -o 011-cloud main.go
fi

LD_LIBRARY_PATH=/root/.vosk/libvosk ./011-cloud
