package db

//package postgres

import (
	"center/internal/config"
	"database/sql"
	"fmt"
	"log"
	"strings"

	// Импортируем все возможные драйверы
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

//
//type PostgresConfig struct {
//	Driver          string
//	Host            string
//	Port            string
//	User            string
//	Password        string
//	DBName          string
//	SSLMode         string
//	MaxOpenConns    uint64
//	MaxIdleConns    uint64
//	ConnMaxLifetime time.Duration
//}

func InitPostgres(cfg config.PostgresConfig) error {
	connStr := generateConnectionString(cfg)
	var err error
	DB, err = sql.Open(cfg.Driver, connStr)
	if err != nil {
		return err
	}
	//fmt.Println("-------------------------")
	if err = DB.Ping(); err != nil {
		return err
	}
	configureConnectionPool(cfg)
	log.Println("Connected to PostgreSQL database")
	return nil
}

func generateConnectionString(cfg config.PostgresConfig) string {
	switch cfg.Driver {
	case "postgres", "pq":
		if cfg.SSLMode == "" {
			cfg.SSLMode = "disable"
		}
		return fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	case "pgx":
		if cfg.SSLMode == "" {
			cfg.SSLMode = "disable"
		}
		return fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	case "sqlite3":
		return cfg.DBName

	default:
		return ""
	}
}

func configureConnectionPool(cfg config.PostgresConfig) {
	//максимальное количество одновременно открытых соединений с БД
	DB.SetMaxOpenConns(int(cfg.MaxOpenConns))
	//количество неактивных соединений, которые сохраняются в пуле
	DB.SetMaxIdleConns(int(cfg.MaxIdleConns))
	//максимальное время жизни соединения
	DB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
}

// Проверяет структуру БД
func EnsurePostgresStructure() error {
	// Проверка существования таблиц
	requiredTables := []string{
		"hosts",
		"host_processes",
		"host_containers",
		"alert_rules",
	}

	for _, table := range requiredTables {
		exists, err := tableExists(table)
		if err != nil {
			return fmt.Errorf("error checking table %s: %w", table, err)
		}
		if !exists {
			return fmt.Errorf("missing required table: %s", table)
		}
	}

	//  Проверка структуры таблиц
	if err := verifyTableStructure("hosts", []ColumnDefinition{
		{Name: "id", Type: "integer", NotNull: true, PrimaryKey: true},
		{Name: "hostname", Type: "character varying(255)", NotNull: true},
		{Name: "ip_address", Type: "character varying(50)", NotNull: true},
		{Name: "agent_port", Type: "integer", Default: "8081"},
		{Name: "priority", Type: "integer", Default: "0"},
		{Name: "is_master", Type: "boolean", Default: "false"},
		{Name: "status", Type: "character varying(50)", Default: "'unknown'::character varying"},
		{Name: "created_at", Type: "timestamp without time zone", Default: "now()"},
		{Name: "updated_at", Type: "timestamp without time zone", Default: "now()"},
	}); err != nil {
		return err
	}

	if err := verifyTableStructure("host_processes", []ColumnDefinition{
		{Name: "id", Type: "integer", NotNull: true, PrimaryKey: true},
		{Name: "host_id", Type: "integer", NotNull: true},
		{Name: "process_name", Type: "character varying(255)", NotNull: true},
	}); err != nil {
		return err
	}

	if err := verifyTableStructure("host_containers", []ColumnDefinition{
		{Name: "id", Type: "integer", NotNull: true, PrimaryKey: true},
		{Name: "host_id", Type: "integer", NotNull: true},
		{Name: "container_name", Type: "character varying(255)", NotNull: true},
	}); err != nil {
		return err
	}

	if err := verifyTableStructure("alert_rules", []ColumnDefinition{
		{Name: "id", Type: "integer", NotNull: true, PrimaryKey: true},
		{Name: "host_id", Type: "integer", NotNull: true},
		{Name: "metric_name", Type: "character varying(100)", NotNull: true},
		{Name: "threshold_value", Type: "double precision", NotNull: true},
		{Name: "condition", Type: "character varying(10)", NotNull: true},
		{Name: "enabled", Type: "boolean", Default: "true"},
	}); err != nil {
		return err
	}

	// 3. Проверка внешних ключей
	foreignKeys := []struct {
		Table     string
		Column    string
		RefTable  string
		RefColumn string
		OnDelete  string
	}{
		{"host_processes", "host_id", "hosts", "id", "CASCADE"},
		{"host_containers", "host_id", "hosts", "id", "CASCADE"},
		{"alert_rules", "host_id", "hosts", "id", "CASCADE"},
	}

	for _, fk := range foreignKeys {
		exists, err := foreignKeyExists(fk.Table, fk.Column, fk.RefTable, fk.RefColumn, fk.OnDelete)
		if err != nil {
			return fmt.Errorf("error checking foreign key: %w", err)
		}
		if !exists {
			return fmt.Errorf("missing foreign key: %s.%s -> %s.%s (%s ON DELETE)",
				fk.Table, fk.Column, fk.RefTable, fk.RefColumn, fk.OnDelete)
		}
	}

	log.Println("Database structure verification successful")
	return nil
}

// Описывает ожидаемую структуру столбца
type ColumnDefinition struct {
	Name       string
	Type       string
	NotNull    bool
	PrimaryKey bool
	Default    string
}

// Проверяет существование таблицы
func tableExists(tableName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = $1
		)`

	var exists bool
	err := DB.QueryRow(query, tableName).Scan(&exists)
	return exists, err
}

// Проверяет, является ли столбец частью первичного ключа
func isPrimaryKeyColumn(tableName, columnName string) (bool, error) {
	query := `
        SELECT COUNT(*)
        FROM information_schema.table_constraints tc
        JOIN information_schema.key_column_usage kcu
            ON tc.constraint_name = kcu.constraint_name
            AND tc.table_schema = kcu.table_schema
        WHERE tc.constraint_type = 'PRIMARY KEY'
            AND tc.table_name = $1
            AND kcu.column_name = $2
    `

	var count int
	err := DB.QueryRow(query, tableName, columnName).Scan(&count)
	return count > 0, err
}

// Проверяет структуру таблицы
func verifyTableStructure(tableName string, columns []ColumnDefinition) error {
	for _, expected := range columns {
		// Получаем информацию о столбце
		actual, err := getColumnInfo(tableName, expected.Name)
		if err != nil {
			return fmt.Errorf("error getting column info for %s.%s: %w", tableName, expected.Name, err)
		}

		// Проверяем тип данных
		if !strings.EqualFold(actual.DataType, expected.Type) {
			return fmt.Errorf("type mismatch for %s.%s: expected %s, got %s",
				tableName, expected.Name, expected.Type, actual.DataType)
		}

		// Проверяем NOT NULL
		if expected.NotNull && actual.IsNullable == "YES" {
			return fmt.Errorf("column %s.%s should be NOT NULL", tableName, expected.Name)
		}

		// Проверяем значение по умолчанию (если указано)
		if expected.Default != "" {
			// Обрабатываем случай, когда значение по умолчанию NULL
			if !actual.ColumnDefault.Valid {
				return fmt.Errorf("default value mismatch for %s.%s: expected '%s', got NULL",
					tableName, expected.Name, expected.Default)
			}

			// Нормализуем значения по умолчанию
			normalizedExpected := strings.ToLower(strings.TrimSpace(expected.Default))
			normalizedActual := strings.ToLower(strings.TrimSpace(actual.ColumnDefault.String))

			if normalizedActual != normalizedExpected {
				return fmt.Errorf("default value mismatch for %s.%s: expected '%s', got '%s'",
					tableName, expected.Name, normalizedExpected, normalizedActual)
			}
		}
		// Проверка первичного ключа
		if expected.PrimaryKey {
			isPK, err := isPrimaryKeyColumn(tableName, expected.Name)
			if err != nil {
				return fmt.Errorf("error checking PK for %s.%s: %w", tableName, expected.Name, err)
			}
			if !isPK {
				return fmt.Errorf("column %s.%s should be primary key", tableName, expected.Name)
			}
		}
	}
	return nil
}

// ColumnInfo хранит информацию о столбце
type ColumnInfo struct {
	ColumnName    string
	DataType      string
	IsNullable    string
	ColumnDefault sql.NullString
}

// getColumnInfo возвращает информацию о столбце
func getColumnInfo(tableName, columnName string) (ColumnInfo, error) {
	query := `
		SELECT 
			column_name, 
			data_type,
			is_nullable,
			column_default
		FROM information_schema.columns 
		WHERE table_name = $1 AND column_name = $2`

	var info ColumnInfo
	err := DB.QueryRow(query, tableName, columnName).Scan(
		&info.ColumnName,
		&info.DataType,
		&info.IsNullable,
		&info.ColumnDefault,
	)

	if err == sql.ErrNoRows {
		return ColumnInfo{}, fmt.Errorf("column %s not found in table %s", columnName, tableName)
	}

	return info, err
}

// foreignKeyExists проверяет существование внешнего ключа
func foreignKeyExists(table, column, refTable, refColumn, onDelete string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM information_schema.table_constraints tc
			JOIN information_schema.key_column_usage kcu
				ON tc.constraint_name = kcu.constraint_name
			JOIN information_schema.referential_constraints rc
				ON tc.constraint_name = rc.constraint_name
			WHERE 
				tc.table_name = $1 AND 
				kcu.column_name = $2 AND
				tc.constraint_type = 'FOREIGN KEY' AND
				rc.unique_constraint_name IN (
					SELECT constraint_name
					FROM information_schema.table_constraints
					WHERE table_name = $3 AND constraint_type = 'PRIMARY KEY'
				) AND
				rc.delete_rule = $4
		)`

	var exists bool
	err := DB.QueryRow(query, table, column, refTable, onDelete).Scan(&exists)
	return exists, err
}
