#!/bin/bash

# monitoring-agent.sh - Агент мониторинга системы на Bash

# Настройки по умолчанию
PORT=8080
COLLECT_INTERVAL=30
HOST_ID=$(hostname)
DATA_DIR="/tmp/monitoring-agent"
CONFIG_FILE="$DATA_DIR/config.json"
METRICS_FILE="$DATA_DIR/metrics.json"
PID_FILE="/var/run/monitoring-agent.pid"
LOG_FILE="/var/log/monitoring-agent.log"

# Флаг для завершения работы
RUNNING=true

# Проверка и установка зависимостей
check_dependencies() {
  for cmd in jq curl netstat ps df free; do
    if ! command -v $cmd &>/dev/null; then
      echo "Ошибка: команда '$cmd' не найдена. Установите необходимые зависимости."
      exit 1
    fi
  done
}

# Инициализация
initialize() {
  mkdir -p "$DATA_DIR"

  # Создаем конфигурацию если не существует
  if [ ! -f "$CONFIG_FILE" ]; then
    echo '{
      "processes": [],
      "containers": [],
      "collect_interval": '$COLLECT_INTERVAL'
    }' > "$CONFIG_FILE"
  fi

  # Создаем пустой файл метрик
  echo '{
    "host_id": "'$HOST_ID'",
    "timestamp": "'$(date -Iseconds)'",
    "system": {
      "cpu": {"usage_percent": 0},
      "ram": {"total": 0, "used": 0, "free": 0, "usage_percent": 0},
      "disk": {"total": 0, "used": 0, "free": 0, "usage_percent": 0}
    },
    "processes": [],
    "ports": []
  }' > "$METRICS_FILE"

  # Сохраняем PID для контроля
  echo $$ > "$PID_FILE"
}

# Обработка сигналов завершения
trap_handler() {
  echo "Получен сигнал завершения, останавливаем агент..."
  RUNNING=false
}

trap 'trap_handler' INT TERM

# Сбор метрик CPU
collect_cpu_metrics() {
  # Используем top для получения загрузки CPU
  cpu_usage=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | sed 's/,/./')

  # Обновляем метрики в JSON
  tmp=$(mktemp)
  jq '.system.cpu.usage_percent = '"$cpu_usage" "$METRICS_FILE" > "$tmp" && mv "$tmp" "$METRICS_FILE"
}

# Сбор метрик RAM
collect_ram_metrics() {
  # Получаем данные о памяти
  mem_data=$(free -b | grep Mem)
  total=$(echo "$mem_data" | awk '{print $2}')
  used=$(echo "$mem_data" | awk '{print $3}')
  free=$(echo "$mem_data" | awk '{print $4}')
  usage_percent=$(echo "scale=2; $used * 100 / $total" | bc)

  # Обновляем метрики в JSON
  tmp=$(mktemp)
  jq '.system.ram = {
    "total": '"$total"',
    "used": '"$used"',
    "free": '"$free"',
    "usage_percent": '"$usage_percent"'
  }' "$METRICS_FILE" > "$tmp" && mv "$tmp" "$METRICS_FILE"
}

# Сбор метрик диска
collect_disk_metrics() {
  # Получаем данные о диске (корневой раздел)
  disk_data=$(df -B1 / | tail -1)
  total=$(echo "$disk_data" | awk '{print $2}')
  used=$(echo "$disk_data" | awk '{print $3}')
  free=$(echo "$disk_data" | awk '{print $4}')
  usage_percent=$(echo "$disk_data" | awk '{print $5}' | sed 's/%//')

  # Обновляем метрики в JSON
  tmp=$(mktemp)
  jq '.system.disk = {
    "total": '"$total"',
    "used": '"$used"',
    "free": '"$free"',
    "usage_percent": '"$usage_percent"'
  }' "$METRICS_FILE" > "$tmp" && mv "$tmp" "$METRICS_FILE"
}

# Сбор метрик процессов
collect_process_metrics() {
  # Получаем список процессов для мониторинга
  processes_list=$(jq -r '.processes[]' "$CONFIG_FILE" 2>/dev/null || echo "")

  # Если список пуст, мониторим все процессы
  if [ -z "$processes_list" ]; then
    process_data=$(ps aux | tail -n +2 | head -10) # Ограничиваем 10 процессами для демонстрации
  else
    process_filter=$(echo "$processes_list" | tr '\n' '|' | sed 's/|$//')
    process_data=$(ps aux | grep -E "$process_filter" | grep -v grep)
  fi

  # Формируем JSON-массив процессов
  process_json="["
  while read -r line; do
    user=$(echo "$line" | awk '{print $1}')
    pid=$(echo "$line" | awk '{print $2}')
    cpu=$(echo "$line" | awk '{print $3}')
    mem=$(echo "$line" | awk '{print $4}')
    cmd=$(echo "$line" | awk '{for(i=11;i<=NF;++i)print $i}' | head -c 50)

    process_json+='{"pid":'$pid',"name":"'$cmd'","cpu_percent":'$cpu',"mem_percent":'$mem'},'
  done <<< "$process_data"
  process_json=${process_json%,}  # Удаляем последнюю запятую
  process_json+="]"

  # Обновляем метрики в JSON
  tmp=$(mktemp)
  jq '.processes = '"$process_json" "$METRICS_FILE" > "$tmp" && mv "$tmp" "$METRICS_FILE"
}

# Сбор сетевых метрик (открытые порты)
collect_network_metrics() {
  # Получаем список открытых портов
  port_data=$(netstat -tuln | grep LISTEN)

  # Формируем JSON-массив портов
  port_json="["
  while read -r line; do
    proto=$(echo "$line" | awk '{print $1}')
    addr=$(echo "$line" | awk '{print $4}')
    port=$(echo "$addr" | awk -F: '{print $NF}')

    port_json+='{"port":'$port',"protocol":"'$proto'","state":"LISTEN"},'
  done <<< "$port_data"
  port_json=${port_json%,}  # Удаляем последнюю запятую
  port_json+="]"

  # Обновляем метрики в JSON
  tmp=$(mktemp)
  jq '.ports = '"$port_json" "$METRICS_FILE" > "$tmp" && mv "$tmp" "$METRICS_FILE"
}

# Функция для обновления временной метки
update_timestamp() {
  tmp=$(mktemp)
  jq '.timestamp = "'$(date -Iseconds)'"' "$METRICS_FILE" > "$tmp" && mv "$tmp" "$METRICS_FILE"
}

# Сборщик метрик (вторая горутина)
metrics_collector() {
  echo "Запуск сборщика метрик с интервалом ${COLLECT_INTERVAL}с"

  while $RUNNING; do
    # Обновляем интервал из конфигурации
    COLLECT_INTERVAL=$(jq -r '.collect_interval' "$CONFIG_FILE")

    # Собираем метрики
    collect_cpu_metrics
    collect_ram_metrics
    collect_disk_metrics
    collect_process_metrics
    collect_network_metrics
    update_timestamp

    # Ждем следующего цикла
    sleep $COLLECT_INTERVAL
  done

  echo "Сборщик метрик остановлен"
}

# HTTP-сервер (первая горутина)
start_http_server() {
  echo "Запуск HTTP-сервера на порту $PORT"

  # Используем временный файл сокета
  SOCKET=$(mktemp)
  rm "$SOCKET" # netcat создаст сам

  # Запуск в бесконечном цикле
  while $RUNNING; do
    # Принимаем одно соединение с помощью netcat
    { echo -ne "HTTP/1.1 200 OK\r\nContent-Type: application/json\r\n\r\n"; cat "$METRICS_FILE"; } | nc -l -p $PORT >/dev/null
  done

  echo "HTTP-сервер остановлен"
}

# Функция для обновления конфигурации (используется внешним клиентом)
update_config() {
  local key=$1
  local value=$2

  case $key in
    processes)
      tmp=$(mktemp)
      jq '.processes = '"$value" "$CONFIG_FILE" > "$tmp" && mv "$tmp" "$CONFIG_FILE"
      ;;
    containers)
      tmp=$(mktemp)
      jq '.containers = '"$value" "$CONFIG_FILE" > "$tmp" && mv "$tmp" "$CONFIG_FILE"
      ;;
    interval)
      tmp=$(mktemp)
      jq '.collect_interval = '"$value" "$CONFIG_FILE" > "$tmp" && mv "$tmp" "$CONFIG_FILE"
      ;;
    *)
      echo "Неизвестный параметр конфигурации: $key"
      return 1
      ;;
  esac

  return 0
}

# Запуск агента
main() {
  echo "Запуск агента мониторинга..."
  check_dependencies
  initialize

  # Запуск сборщика метрик в фоновом режиме (вторая горутина)
  metrics_collector &
  COLLECTOR_PID=$!

  # Запуск HTTP-сервера (первая горутина)
  start_http_server &
  SERVER_PID=$!

  # Ждем сигнала завершения
  wait $COLLECTOR_PID
  wait $SERVER_PID

  echo "Агент остановлен"
  rm -f "$PID_FILE"
}

# Запуск основной функции
main "$@"