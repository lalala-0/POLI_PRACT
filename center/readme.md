# Задача: 
Реализовать информационную систему центра мониторинга хостов и контейнеров. На Linux
Описание информационной системы
Информационная система состоит из программного обеспечения центра мониторинга (далее ЦМ) и агентов центра мониторинга (далее агент ЦМ).
ЦМ должен состоять из модулей:
- ведение и учет хостов;
- мониторинг состояния добавленных хостов с возможностью отправки уведомления;
- определение мастер хоста на основе указанного при добавлении приоритета.

Агент ЦМ должен состоять из следующих модулей:
- сбор основных метрик хоста (CPU, RAM, Disk);
- сбор данных о состоянии процессов, указанных в конфигурации хоста в ЦМ;
- сбор сведений о сетевых портах ТСР и UDP;
- сбор данных о состоянии Docker контейнеров, указанных в конфигурации хоста в ЦМ.

Технологический стек:
Архитектура взаимодействия: REST API
Язык программирования: Go (Gin)
База данных: PostgreSQL, MongoDB
Инструмент генерации документации к API: Swagger (Postman)
Инструмент автоматизации сборки и развертывания: Docker


# Конфиг ЦМ
Порт, на котором запускается ЦМ, подключение к бд,  инфа про сбор метрик (интервал запроса, их время жизни, какие собирать), и то, что будет в постгрес записано (инфа про хосты, отслеживаемые на них процессы и контейнеры)
```yaml
server:
  port: "8080"
  read_timeout: 30s
  write_timeout: 30s

postgres:
  host: "build-postgres-1"
  port: "5432"
  user: "postgres"
  password: "password"
  dbname: "monitoring"
  sslmode: "disable"
  driver: "postgres"
  maxOpenConns: 25
  maxIdleConns: 5
  connMaxLifetime: 5m

mongodb:
  uri: "mongodb://mongodb:27017"
  dbname: "metrics"
  connectTimeout: 5s
  maxPoolSize: 20
  minPoolSize: 5
  serverSelectionTimeout: 10s

metrics:
  poll_interval: 60s
  metrics_ttl_days: 14
  self_check_interval: 5m

  system:
    enabled: true
    collect_cpu: true
    collect_ram: true
    collect_disks: true

  process:
    enabled: true

  network:
    enabled: true
    monitor_tcp: true
    monitor_udp: true

  container:
    enabled: true

logging:
  level: "info"
  file_path: "/app/logs/monitoring.log"
  max_size: 100
  max_backups: 5
  max_age: 30
  compress: true

initial_data:
  hosts:
    - hostname: "12345"
      ip_address: 192.168.58.129
      agent_port: 8081
      priority: 10
      is_master: true
      status: "active"
      processes:
        - "nginx"
        - "postgres"
        - "redis"
      containers:
        - "build-mongodb-1"
        - "build-postgres-1"
      alerts:
        - metric_name: "cpu_usage_percent"
          threshold_value: 90.0
          condition: ">"
          enabled: true
        - metric_name: "memory_usage_percent"
          threshold_value: 85.0
          condition: ">"
          enabled: true
    - hostname: "server-2"
      ip_address: "192.168.1.102"
      agent_port: 8081
      priority: 5
      is_master: false
      processes:
        - "nginx"
        - "redis"
      containers:
        - "cache-redis"

```

# Как запускать

1. Запуск Docker контейнера
```bash
sudo docker compose up --build -d
```

2. Просмотр логов
```bash
sudo docker compose logs monitoring-center
```

3. Завершить работу центра мониторинга
```bash
sudo docker compose down -v
```




