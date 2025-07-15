package repositories

import (
	"center/internal/models"
	"database/sql"
	"go.mongodb.org/mongo-driver/mongo"

	"context"
	//"database/sql"
	"time"
)

// HostRepository интерфейс для работы с хостами в БД
type HostRepository interface {
	NewHostRepository(db *sql.DB) *HostRepository
	GetAll(ctx context.Context) ([]models.Host, error)
	GetHostCount() (int, error)
	GetByID(ctx context.Context, id int) (*models.Host, error)
	Create(ctx context.Context, host *models.Host) (int, error)
	Update(ctx context.Context, host *models.Host) error
	Delete(ctx context.Context, id int) error
	UpdateStatus(ctx context.Context, id int, status string) error
	GetMaster(ctx context.Context) (*models.Host, error)
	SetMaster(ctx context.Context, id int) error
	//Ping(ctx context.Context) error
}

// ProcessRepository интерфейс для работы с процессами в БД
type ProcessRepository interface {
	NewProcessRepository(db *sql.DB) *ProcessRepository
	GetByHostID(ctx context.Context, hostID int) ([]models.Process, error)
	GetByID(ctx context.Context, id int) (*models.Process, error)
	Create(ctx context.Context, process *models.Process) (int, error)
	Update(ctx context.Context, process *models.Process) error
	Delete(ctx context.Context, id int) error
	Exists(ctx context.Context, hostID int, processName string) (bool, error)
}

// ContainerRepository интерфейс для работы с контейнерами в БД
type ContainerRepository interface {
	NewContainerRepository(db *sql.DB) *ContainerRepository
	GetByHostID(ctx context.Context, hostID int) ([]models.Container, error)
	GetByID(ctx context.Context, id int) (*models.Container, error)
	Create(ctx context.Context, container *models.Container) (int, error)
	Update(ctx context.Context, container *models.Container) error
	Delete(ctx context.Context, id int) error
	Exists(ctx context.Context, hostID int, containerName string) (bool, error)
}

// AlertRepository интерфейс для работы с правилами оповещений в БД
type AlertRepository interface {
	NewAlertRepository(db *sql.DB) *AlertRepository
	GetByHostID(ctx context.Context, hostID int) ([]models.AlertRule, error)
	GetByID(ctx context.Context, id int) (*models.AlertRule, error)
	Create(ctx context.Context, alert *models.AlertRule) (int, error)
	Update(ctx context.Context, alert *models.AlertRule) error
	Delete(ctx context.Context, id int) error
	SetEnabled(ctx context.Context, id int, enabled bool) error
	GetActive(ctx context.Context) ([]models.AlertRule, error)
}

// MetricRepository интерфейс для работы с метриками в MongoDB
type MetricRepository interface {
	NewMetricRepository(db *mongo.Database) *MetricRepository
	SaveSystemMetrics(ctx context.Context, metrics *models.SystemMetrics) error
	SaveProcessMetrics(ctx context.Context, metrics *models.ProcessMetrics) error
	SaveContainerMetrics(ctx context.Context, metrics *models.ContainerMetrics) error
	SaveNetworkMetrics(ctx context.Context, metrics *models.NetworkMetrics) error
	GetLastSystemMetrics(ctx context.Context, hostID int) (*models.SystemMetrics, error)
	GetSystemMetricsInRange(ctx context.Context, hostID int, from, to time.Time) ([]models.SystemMetrics, error)
	GetProcessMetricsInRange(ctx context.Context, hostID int, from, to time.Time) ([]models.ProcessMetrics, error)
	GetContainerMetricsInRange(ctx context.Context, hostID int, from, to time.Time) ([]models.ContainerMetrics, error)
	GetNetworkMetricsInRange(ctx context.Context, hostID int, from, to time.Time) ([]models.NetworkMetrics, error)
	SetupTTLIndex(ctx context.Context, collectionName string, ttlSeconds int32) error
	CleanupOldMetrics(ctx context.Context, collectionName string, threshold time.Time) error
	Ping(ctx context.Context) error
}
