#!/bin/bash

set -e

if [[ -f main.go ]]; then
    echo "This script is only meant to be run in a server deployment."
    exit 0
fi

if [[ $1 == "-u" ]]; then
	if [[ $2 == "" ]]; then
		echo "You must also provide the location of the current install."
		exit 0
	fi
	ORIGDIR="$(pwd)"
	cd $2
    if [[ ! -f 011-cloud ]]; then
        echo "This is the wrong directory."
    fi
	sudo systemctl stop 011-cloud
    sync
	sudo rm -rf webroot install.sh 011-cloud 011-cloud.sh stttest.pcm intent-data weather-map.json model vosk opus
    sudo cp -r ${ORIGDIR}/webroot .
    sudo cp ${ORIGDIR}/011-cloud .
    sudo cp ${ORIGDIR}/install.sh .
    sudo cp ${ORIGDIR}/011-cloud.sh .
    sudo cp ${ORIGDIR}/stttest.pcm .
    sudo cp -r ${ORIGDIR}/intent-data .
    sudo cp ${ORIGDIR}/weather-map.json .
    sudo cp -r ${ORIGDIR}/model .
    sudo cp -r ${ORIGDIR}/vosk .
    sudo cp -r ${ORIGDIR}/opus .
    sync
    sudo systemctl start 011-cloud
fi
	

echo "[Unit]" >011-cloud.service
echo "Description=Vector Cloud for 0.11" >>011-cloud.service
echo "StartLimitIntervalSec=500" >>011-cloud.service
echo "StartLimitBurst=5" >>011-cloud.service
echo >>011-cloud.service
echo "[Service]" >>011-cloud.service
echo "Type=simple" >>011-cloud.service
echo "Restart=on-failure" >>011-cloud.service
echo "RestartSec=5s" >>011-cloud.service
echo "WorkingDirectory=$(readlink -f ./)" >>011-cloud.service
echo "ExecStart=$(readlink -f ./011-cloud.sh)" >>011-cloud.service
echo >>011-cloud.service
echo "[Install]" >>011-cloud.service
echo "WantedBy=multi-user.target" >>011-cloud.service

cat 011-cloud.service

echo

sudo mv 011-cloud.service /lib/systemd/system/

sudo systemctl daemon-reload

echo "Daemon installed. Create a source.sh file with your desired settings and run systemctl start 011-cloud."
