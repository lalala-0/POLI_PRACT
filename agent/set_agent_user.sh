#!/bin/bash
#Создание специального пользователя для агента
sudo useradd --system --no-create-home --shell /bin/false agentuser
echo "agentuser created"
#Добавление пользователя в группу docker
sudo usermod -aG docker agentuser
echo "agentuser added in docker group"

#Изменение прав на файлы
sudo chown -R agentuser:docker /bin/agent
sudo chmod 750 /bin/agent/main
sudo chown -R agentuser:docker /etc/agent
sudo chmod 640 /etc/agent/config.yml
