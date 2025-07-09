package db

import (
	"database/sql"
	"fmt"
	"log"

	// Импортируем все возможные драйверы
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"
    _ "github.com/mattn/go-sqlite3"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var DB *sql.DB

type PostgresConfig struct {
    Driver   string
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

func InitPostgres(cfg PostgresConfig) error {
	connStr := generateConnectionString(cnf)
	var err error
	DB, err = sql.Open(cfg.driver, connStr)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	log.Println("Connected to PostgreSQL database")
	return nil
}


func generateConnectionString(cfg PostgresConfig) string {
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

func configureConnectionPool() {
	//максимальное количество одновременно открытых соединений с БД
    DB.SetMaxOpenConns(25)
	//количество неактивных соединений, которые сохраняются в пуле
    DB.SetMaxIdleConns(5)
	//максимальное время жизни соединения
    DB.SetConnMaxLifetime(5 * time.Minute)
}


