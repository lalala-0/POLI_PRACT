package transport

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// @title Agent API
// @version 1.0
// @description API сервиса мониторинга
// @host localhost:8080
// @BasePath /

// healthCheck проверяет работоспособность агента
// @Summary Проверка работоспособности
// @Description Проверяет работоспособность агента и возвращает статус
// @Tags system
// @Produce json
// @Success 200 {object} object{status=string,message=string} "Агент работает"
// @Router /api/health [get]
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Агент работает",
	})
}

// getMetrics возвращает все собранные метрики
// @Summary Получение всех метрик
// @Description Возвращает все собранные метрики агента
// @Tags metrics
// @Produce json
// @Success 200 {object} models.AgentMetrics "Все метрики"
// @Router /api/metrics [get]
func (s *Server) getMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, s.lastMetrics)
}

// getSystemMetrics возвращает только системные метрики
// @Summary Получение системных метрик
// @Description Возвращает только системные метрики агента
// @Tags metrics
// @Produce json
// @Success 200 {object} object{host_id=string,timestamp=string,system=object} "Системные метрики"
// @Router /api/metrics/system [get]
func (s *Server) getSystemMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"host_id":   s.lastMetrics.HostID,
		"timestamp": s.lastMetrics.Timestamp,
		"system":    s.lastMetrics.System,
	})
}

// getProcessMetrics возвращает только метрики процессов
// @Summary Получение метрик процессов
// @Description Возвращает метрики отслеживаемых процессов
// @Tags metrics
// @Produce json
// @Success 200 {object} object{host_id=string,timestamp=string,processes=array} "Метрики процессов"
// @Failure 400 {object} object{status=string,message=string} "Список отслеживаемых процессов не настроен"
// @Router /api/metrics/processes [get]
func (s *Server) getProcessMetrics(c *gin.Context) {
	if !s.metricsService.IsProcessConfigSet() {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Список отслеживаемых процессов не настроен",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"host_id":   s.lastMetrics.HostID,
		"timestamp": s.lastMetrics.Timestamp,
		"processes": s.lastMetrics.Processes,
	})
}

// getNetworkMetrics возвращает только сетевые метрики
// @Summary Получение сетевых метрик
// @Description Возвращает информацию о сетевых соединениях и открытых портах
// @Tags metrics
// @Produce json
// @Success 200 {object} object{host_id=string,timestamp=string,ports=array} "Сетевые метрики"
// @Router /api/metrics/network [get]
func (s *Server) getNetworkMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"host_id":   s.lastMetrics.HostID,
		"timestamp": s.lastMetrics.Timestamp,
		"ports":     s.lastMetrics.Ports,
	})
}

// getContainerMetrics возвращает только метрики контейнеров
// @Summary Получение метрик контейнеров
// @Description Возвращает метрики отслеживаемых Docker контейнеров
// @Tags metrics
// @Produce json
// @Success 200 {object} object{host_id=string,timestamp=string,containers=array} "Метрики контейнеров"
// @Failure 400 {object} object{status=string,message=string} "Список отслеживаемых контейнеров не настроен"
// @Router /api/metrics/containers [get]
func (s *Server) getContainerMetrics(c *gin.Context) {
	if !s.metricsService.IsContainerConfigSet() {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Список отслеживаемых контейнеров не настроен",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"host_id":    s.lastMetrics.HostID,
		"timestamp":  s.lastMetrics.Timestamp,
		"containers": s.lastMetrics.Containers,
	})
}

// updateProcessConfig обновляет список отслеживаемых процессов
// @Summary Обновление списка отслеживаемых процессов
// @Description Устанавливает список процессов, метрики которых будут собираться
// @Tags configuration
// @Accept json
// @Produce json
// @Param request body object{processes=array} true "Массив имён процессов для отслеживания"
// @Success 200 {object} object{status=string,message=string} "Конфигурация успешно обновлена"
// @Failure 400 {object} object{status=string,message=string} "Некорректный формат данных или пустой список"
// @Failure 500 {object} object{status=string,message=string} "Внутренняя ошибка сервера"
// @Router /api/config/processes [post]
func (s *Server) updateProcessConfig(c *gin.Context) {
	var config struct {
		Processes []string `json:"processes"`
	}

	if err := c.BindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	if len(config.Processes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Список процессов не может быть пустым",
		})
		return
	}

	// Обновляем конфигурацию через сервис
	if err := s.metricsService.UpdateProcessConfig(config.Processes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось обновить конфигурацию",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Конфигурация процессов обновлена",
	})
}

// updateContainerConfig обновляет список отслеживаемых контейнеров
// @Summary Обновление списка отслеживаемых контейнеров
// @Description Устанавливает список Docker контейнеров, метрики которых будут собираться
// @Tags configuration
// @Accept json
// @Produce json
// @Param request body object{containers=array} true "Массив имён контейнеров для отслеживания"
// @Success 200 {object} object{status=string,message=string} "Конфигурация успешно обновлена"
// @Failure 400 {object} object{status=string,message=string} "Некорректный формат данных или пустой список"
// @Failure 500 {object} object{status=string,message=string} "Внутренняя ошибка сервера"
// @Router /api/config/containers [post]
func (s *Server) updateContainerConfig(c *gin.Context) {
	var config struct {
		Containers []string `json:"containers"`
	}

	if err := c.BindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	if len(config.Containers) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Список контейнеров не может быть пустым",
		})
		return
	}

	if err := s.metricsService.UpdateContainerConfig(config.Containers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось обновить конфигурацию",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Конфигурация контейнеров обновлена",
	})
}

// updateCollectionInterval обновляет интервал сбора метрик
// @Summary Обновление интервала сбора метрик
// @Description Изменяет временной интервал между сборами метрик агента
// @Tags configuration
// @Accept json
// @Produce json
// @Param request body object{interval_seconds=integer} true "Интервал сбора в секундах (положительное число)"
// @Success 200 {object} object{status=string,message=string,interval_seconds=integer} "Интервал успешно обновлен"
// @Failure 400 {object} object{status=string,message=string} "Некорректные входные данные"
// @Failure 500 {object} object{status=string,message=string} "Внутренняя ошибка сервера"
// @Router /api/config/collection-interval [post]
// updateCollectionInterval обновляет интервал сбора метрик
// updateCollectionInterval обновляет интервал сбора метрик
func (s *Server) updateCollectionInterval(c *gin.Context) {
	var config struct {
		Interval int64 `json:"interval_seconds"` // Интервал в секундах
	}

	if err := c.BindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	if config.Interval <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Интервал сбора метрик должен быть положительным числом",
		})
		return
	}

	// Преобразуем секунды в Duration
	interval := time.Duration(config.Interval) * time.Second

	// Предполагается, что метод UpdateCollectionInterval существует в сервисе метрик
	if err := s.metricsService.UpdateCollectionInterval(interval); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось обновить интервал сбора метрик",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":           "success",
		"message":          "Интервал сбора метрик обновлен",
		"interval_seconds": config.Interval,
	})
}
