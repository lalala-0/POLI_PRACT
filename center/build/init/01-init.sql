-- Создание таблицы хостов
CREATE TABLE hosts (
    id SERIAL PRIMARY KEY,
    hostname VARCHAR(255) NOT NULL,
    ip_address VARCHAR(50) NOT NULL,
    agent_port INTEGER DEFAULT 8081,
    priority INTEGER DEFAULT 0,
    is_master BOOLEAN DEFAULT FALSE,
    status VARCHAR(50) DEFAULT 'unknown',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Создание таблицы для хранения списка отслеживаемых процессов
CREATE TABLE host_processes (
    id SERIAL PRIMARY KEY,
    host_id INTEGER REFERENCES hosts(id) ON DELETE CASCADE,
    process_name VARCHAR(255) NOT NULL
);

-- Создание таблицы для хранения списка отслеживаемых контейнеров
CREATE TABLE host_containers (
    id SERIAL PRIMARY KEY,
    host_id INTEGER REFERENCES hosts(id) ON DELETE CASCADE,
    container_name VARCHAR(255) NOT NULL
);

CREATE TABLE alert_rules (
    id SERIAL PRIMARY KEY,
    host_id INT REFERENCES hosts(id) ON DELETE CASCADE,
    metric_name VARCHAR(100) NOT NULL,
    threshold_value FLOAT NOT NULL,
    condition VARCHAR(10) NOT NULL,
    enabled BOOLEAN DEFAULT TRUE
);