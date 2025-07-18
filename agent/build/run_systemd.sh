#!/bin/bash
sudo mkdir -p /etc/agent
sudo cp ../config/config.yml /etc/agent/
sudo mkdir -p /bin/agent
sudo cp ./main /bin/agent/
sudo cp ../deployments/agent.service /etc/systemd/system/
echo "copy to system Done"

sudo systemctl daemon-reload
echo "daemon-reload Done"

sudo systemctl enable agent.service
echo "enable agent.service Done"

sudo systemctl start agent.service
echo "start agent.service Done"

sudo systemctl status agent.service
echo "status agent.service Done"
