-- Создание таблицы хостов
CREATE TABLE hosts (
    id SERIAL PRIMARY KEY,
    hostname VARCHAR(255) NOT NULL,
    ip_address VARCHAR(50) NOT NULL,
    agent_port INTEGER NOT NULL DEFAULT 8081,
    priority INTEGER NOT NULL DEFAULT 0,
    is_master BOOLEAN NOT NULL DEFAULT FALSE,
    status VARCHAR(50) NOT NULL DEFAULT 'unknown',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Создание таблицы для хранения списка отслеживаемых процессов
CREATE TABLE host_processes (
    id SERIAL PRIMARY KEY,
    host_id INTEGER NOT NULL REFERENCES hosts(id) ON DELETE CASCADE,
    process_name VARCHAR(255) NOT NULL
);

-- Создание таблицы для хранения списка отслеживаемых контейнеров
CREATE TABLE host_containers (
    id SERIAL PRIMARY KEY,
    host_id INTEGER NOT NULL REFERENCES hosts(id) ON DELETE CASCADE,
    container_name VARCHAR(255) NOT NULL
);

CREATE TABLE alert_rules (
    id SERIAL PRIMARY KEY,
    host_id INTEGER NOT NULL REFERENCES hosts(id) ON DELETE CASCADE,
    metric_name VARCHAR(100) NOT NULL,
    threshold_value FLOAT NOT NULL,
    condition VARCHAR(10) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE
);