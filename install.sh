#!/bin/bash

if [[ -f main.go ]]; then
    echo "This script is only meant to be run in a server deployment."
    exit 0
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
sudo systemctl enable 011-cloud

echo "Daemon installed. Create a source.sh file with your desired settings and run systemctl start 011-cloud."