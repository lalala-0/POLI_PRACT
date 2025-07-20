package services

import (
	"center/internal/config"
	"center/internal/database/mongodb/repositories"
	pg_repo "center/internal/database/postgres/repositories"
	"center/internal/models"
	"context"
	"errors"
	"log"
	"time"
)

// HostService реализует бизнес-логику работы с хостами
type HostService struct {
	HostRepo      pg_repo.PostgresHostRepository
	ProcessRepo   pg_repo.PostgresProcessRepository
	ContainerRepo pg_repo.PostgresContainerRepository
	AlertRepo     pg_repo.PostgresAlertRepository
	MetricRepo    repositories.MongoMetricRepository
}

func NewHostService(
	hostRepo pg_repo.PostgresHostRepository,
	processRepo pg_repo.PostgresProcessRepository,
	containerRepo pg_repo.PostgresContainerRepository,
	alertRepo pg_repo.PostgresAlertRepository,
	metricRepo repositories.MongoMetricRepository,
) *HostService {
	return &HostService{
		HostRepo:      hostRepo,
		ProcessRepo:   processRepo,
		ContainerRepo: containerRepo,
		AlertRepo:     alertRepo,
		MetricRepo:    metricRepo,
	}
}

// Host Operations
func (s *HostService) CreateHost(ctx context.Context, hostInput models.HostInput) (int, error) {
	host := models.Host{
		Hostname:  hostInput.Hostname,
		IPAddress: hostInput.IPAddress,
		AgentPort: hostInput.AgentPort,
		Priority:  hostInput.Priority,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return s.HostRepo.Create(ctx, &host)
}

func (s *HostService) GetHost(ctx context.Context, id int) (*models.Host, error) {
	return s.HostRepo.GetByID(ctx, id)
}

func (s *HostService) UpdateHost(ctx context.Context, id int, hostInput models.HostInput) error {
	host, err := s.HostRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	host.Hostname = hostInput.Hostname
	host.IPAddress = hostInput.IPAddress
	host.AgentPort = hostInput.AgentPort
	host.Priority = hostInput.Priority
	host.UpdatedAt = time.Now()

	err = s.HostRepo.Update(ctx, host)
	if err != nil {
		return err
	}

	// Если изменился приоритет, перевыбираем мастера
	if hostInput.Priority != 0 {
		if err := s.electMasterHost(ctx); err != nil {
			log.Printf("Failed to elect master host: %v", err)
		}
	}

	return nil
}

func (s *HostService) GetAllHosts(ctx context.Context) ([]models.Host, error) {
	return s.HostRepo.GetAll(ctx)
}

func (s *HostService) UpdateHostStatus(ctx context.Context, id int, status string) error {
	return s.HostRepo.UpdateStatus(ctx, id, status)
}

func (s *HostService) DeleteHost(ctx context.Context, id int) error {
	if err := s.HostRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Перевыбираем мастера после удаления
	return s.electMasterHost(ctx)
}

func (s *HostService) SetMasterHost(ctx context.Context, id int) error {
	// Установка нового мастера
	return s.HostRepo.SetMaster(ctx, id)
}

// Автоматически выбирает мастер-хост на основе приоритета
func (s *HostService) electMasterHost(ctx context.Context) error {
	hosts, err := s.HostRepo.GetAll(ctx)
	if err != nil {
		return err
	}
	//log.Println("------------------------1------------------")
	if len(hosts) == 0 {
		return nil
	}

	//log.Println("------------------------2------------------")
	// Ищем хост с наивысшим приоритетом
	var masterHost *models.Host
	for _, host := range hosts {
		if masterHost == nil {
			masterHost = &host
			continue
		}

		// Сравниваем приоритеты
		if host.Priority > masterHost.Priority {
			masterHost = &host
		} else if host.Priority == masterHost.Priority {
			// При равном приоритете выбираем самый старый хост
			if host.CreatedAt.Before(masterHost.CreatedAt) {
				masterHost = &host
			}
		}
	}

	//log.Println("------------------------------------------", masterHost.ID)
	// Устанавливаем найденный хост как мастер
	return s.SetMasterHost(ctx, masterHost.ID)
}

// Process Operations
func (s *HostService) AddProcess(ctx context.Context, hostID int, processName string) (int, error) {
	exists, err := s.ProcessRepo.Exists(ctx, hostID, processName)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("process already monitored")
	}

	process := &models.Process{
		HostID:      hostID,
		ProcessName: processName,
	}
	return s.ProcessRepo.Create(ctx, process)
}

// Container Operations
func (s *HostService) AddContainer(ctx context.Context, hostID int, containerName string) (int, error) {
	exists, err := s.ContainerRepo.Exists(ctx, hostID, containerName)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("container already monitored")
	}

	container := &models.Container{
		HostID:        hostID,
		ContainerName: containerName,
	}
	return s.ContainerRepo.Create(ctx, container)
}

// Alert Operations
func (s *HostService) CreateAlertRule(ctx context.Context, hostID int, alertInput models.AlertInput) (int, error) {
	rule := &models.AlertRule{
		HostID:         hostID,
		MetricName:     alertInput.MetricName,
		ThresholdValue: alertInput.ThresholdValue,
		Condition:      alertInput.Condition,
		Enabled:        alertInput.Enabled,
	}
	return s.AlertRepo.Create(ctx, rule)

}
func (s *HostService) GetAlertsByHostID(ctx context.Context, hostID int) ([]models.AlertRule, error) {
	return s.AlertRepo.GetByHostID(ctx, hostID)
}

// Metrics Operations
func (s *HostService) SaveSystemMetrics(ctx context.Context, metrics *models.SystemMetrics) error {
	return s.MetricRepo.SaveSystemMetrics(ctx, metrics)
}

func (s *HostService) SaveProcessMetrics(ctx context.Context, metrics *models.ProcessMetrics) error {
	return s.MetricRepo.SaveProcessMetrics(ctx, metrics)
}

func (s *HostService) SaveContainerMetrics(ctx context.Context, metrics *models.ContainerMetrics) error {
	return s.MetricRepo.SaveContainerMetrics(ctx, metrics)
}

func (s *HostService) SaveNetworkMetrics(ctx context.Context, metrics *models.NetworkMetrics) error {
	return s.MetricRepo.SaveNetworkMetrics(ctx, metrics)
}

func (s *HostService) LoadInitialData(ctx context.Context, cfg *config.AppConfig) error {
	//Проверка, есть ли уже данные
	count, err := s.HostRepo.GetHostCount()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Данные уже есть, пропускаем
	}

	// Загрузка данных из конфига
	for _, hostCfg := range cfg.InitialData.Hosts {
		// Создание хоста через сервис
		hostID, err := s.CreateHost(ctx, models.HostInput{
			Hostname:  hostCfg.Hostname,
			IPAddress: hostCfg.IPAddress,
			AgentPort: hostCfg.AgentPort,
			Priority:  hostCfg.Priority,
		})
		if err != nil {
			log.Printf("Failed to create host %s: %v", hostCfg.Hostname, err)
			continue
		}

		// Добавление процессов
		for _, process := range hostCfg.Processes {
			if _, err := s.AddProcess(ctx, hostID, process); err != nil {
				log.Printf("Failed to add process %s to host %s: %v", process, hostCfg.Hostname, err)
			}
		}

		// Добавление контейнеров
		for _, container := range hostCfg.Containers {
			if _, err := s.AddContainer(ctx, hostID, container); err != nil {
				log.Printf("Failed to add container %s to host %s: %v", container, hostCfg.Hostname, err)
			}
		}

		// Добавление правил оповещений
		for _, alert := range hostCfg.Alerts {
			if _, err := s.CreateAlertRule(ctx, hostID, models.AlertInput{
				MetricName:     alert.MetricName,
				ThresholdValue: alert.ThresholdValue,
				Condition:      alert.Condition,
				Enabled:        alert.Enabled,
			}); err != nil {
				log.Printf("Failed to add alert for %s to host %s: %v", alert.MetricName, hostCfg.Hostname, err)
			}
		}

	}

	log.Println("Initial data loaded from config using services")
	// После загрузки всех хостов выбираем мастера
	return s.electMasterHost(ctx)
}
