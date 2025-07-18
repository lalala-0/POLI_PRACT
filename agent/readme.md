# Как сие работает? Работа агента мониторинга с двумя горутинами

Агент мониторинга построен на взаимодействии двух горутин, которые разделяют между собой конфигурацию и метрики:

## Горутина 1: HTTP-сервер (транспортный слой)

**Что делает:**
- Принимает HTTP-запросы от пользователей
- Предоставляет доступ к собранным метрикам через API
- Изменяет конфигурацию по запросу пользователя
- Хранит последние полученные метрики в переменной `lastMetrics`

**API эндпоинты:**
- `GET /metrics/*` - получение метрик
- `POST /config/*` - изменение конфигурации

## Горутина 2: Сборщик метрик (коллекторы)

**Что делает:**
- Периодически (с интервалом из конфигурации) собирает метрики системы
- Использует различные коллекторы для сбора разных типов метрик
- Отправляет собранные метрики в канал для HTTP-сервера
- Обновляет настройки коллекторов на основе текущей конфигурации

## Механизм взаимодействия

1. **Обмен метриками:**
   ```
   Горутина 2 (сборщик) ---> Канал metricsCh ---> Горутина 1 (HTTP-сервер)
   ```

2. **Обмен конфигурацией:**
   ```
   Клиент ---> HTTP API ---> MetricsService ---> Коллекторы
   ```

3. **Синхронизация данных:**
    - `metricsService` использует `sync.RWMutex` для безопасного доступа
    - Метрики передаются через буферизованный канал (`metricsCh`)

## Пример потока данных

1. Клиент отправляет запрос на обновление списка процессов:
   ```
   POST /config/processes --> updateProcessConfig() --> metricsService.UpdateProcessConfig()
   ```

2. Горутина сборщика при следующей итерации:
   ```
   Проверяет изменения в конфигурации --> Обновляет коллекторы --> Собирает метрики
   ```

3. Собранные метрики передаются в HTTP-сервер:
   ```
   Сборщик --> metricsCh --> HTTP-сервер.lastMetrics
   ```

4. Клиент запрашивает метрики:
   ```
   GET /metrics --> HTTP-сервер возвращает lastMetrics
   ```

Такая архитектура обеспечивает разделение ответственности и эффективный обмен данными между компонентами системы.

# Вопросики

- Порт, на котором запускается агент, задается в конфиге или в командной строке?
- Агент должен собирать метрики только при получении get запроса от ЦМ ИЛИ собирать метрики с некоторым интервалом и отправлять при запросе ЦМ, то что он уже насобирал? (пока делается второе)
- Апи нормальное?
- Метрики нормальные? Надо чем-то дополнить?


# Про скрипт 

Установка и использование
Сохраните скрипт как monitoring-agent.sh и сделайте его исполняемым:

`chmod +x monitoring-agent.sh`

Установите зависимости (на Debian/Ubuntu):


`sudo apt update && sudo apt install jq bc netcat procps`

Запустите агент:

`sudo ./monitoring-agent.sh`

Проверка

`curl http://localhost:8080/`



# Как добавить в автозапуск через systemctl

1. Настроить users и groups
```bash
./build/set_agent_user.sh
```
Или

1.1. Создание специального пользователя для агента
```bash
sudo useradd --system --no-create-home --shell /bin/false agentuser
```
1.2. Добавление пользователя в группу docker
```bash
sudo usermod -aG docker agentuser
```
1.3. Изменение прав на файлы
```bash
sudo chown -R agentuser:docker /bin/agent
sudo chmod 750 /bin/agent/main
sudo chown -R agentuser:docker /etc/agent
sudo chmod 640 /etc/agent/config.yml
```

2. Компилируем бинарник
```bash
./build/build_agent.sh
```
Или `go build ./cmd/main.go`

3. Запуск агента как юнита system
```bash
./build/build_agent.sh
```
Или 

3.1. Или распределяем необходимые файлы по директориям
- Создаем директорию `sudo mkdir -p /etc/agent`
- Конфиг config.yaml добавляем в папку /etc/agent `sudo cp ./config/config.yml /etc/agent/`
- Создаем директорию `sudo mkdir -p /bin/agent`
- Бинарник main.exe добавляем в папку /bin/agent `sudo cp ./main /bin/agent/`

3.2. Создаём agent.service
- копируем юнит-файл `sudo cp ./deployments/agent.service /etc/systemd/system/`

3.3. Перезапускаем systemd
```bash
sudo systemctl daemon-reload
sudo systemctl enable agent.service
sudo systemctl start agent.service
sudo systemctl status agent.service
```

4. Просмотр журнала
```bash
./build/journal.sh	
```
Или `journalctl -u agent.service -f`


5 Завершение работы агента
```bash
./build/stop_systemd_agent.sh
```
Или

```bash
sudo systemctl stop agent.service
sudo systemctl disable agent.service
sudo systemctl status agent.service
```
