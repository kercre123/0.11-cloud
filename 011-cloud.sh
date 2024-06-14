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

if [[ $1 == "-d" ]]; then
	if [[ -f main.go ]]; then
		echo "Removing cloud program"
		rm -rf 011-cloud
	else
		echo "Only use this flag in the source directory."
	fi
fi

if [[ ! -f 011-cloud ]]; then
	echo "building the program..."
	CGO_LDFLAGS=-L"$(pwd)/vosk" CGO_CFLAGS=-I"$(pwd)/vosk" go build -o 011-cloud main.go
fi

LD_LIBRARY_PATH=$(pwd)/vosk:$(pwd)/opus/lib ./011-cloud
