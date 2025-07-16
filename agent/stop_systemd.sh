#!/bin/bash
sudo systemctl stop agent.service
sudo systemctl disable agent.service
sudo systemctl status agent.service

