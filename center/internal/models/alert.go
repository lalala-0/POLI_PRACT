package models

import "time"

// AlertRule представляет правило для генерации уведомлений
type AlertRule struct {
	ID                int     `json:"id" db:"id"`
	HostID            int     `json:"host_id" db:"host_id"`
	MetricName        string  `json:"metric_name" binding:"required" db:"metric_name"`
	ThresholdValue    float64 `json:"threshold_value" binding:"required" db:"threshold_value"`
	Condition         string  `json:"condition" binding:"required" db:"condition"` // "greater", "less", "equal"
	Enabled           bool    `json:"enabled" db:"enabled"`
	TimeWindowSeconds int     `json:"time_window_seconds" db:"time_window_seconds"` // Например, 60
	FailPercent       float64 `json:"fail_percent" db:"fail_percent"`
}

// AlertInput представляет данные для создания правила оповещения
type AlertInput struct {
	MetricName     string  `json:"metric_name" binding:"required"`
	ThresholdValue float64 `json:"threshold_value" binding:"required"`
	Condition      string  `json:"condition" binding:"required"`
	Enabled        bool    `json:"enabled"`
}

// Alert представляет сгенерированное оповещение
type Alert struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	HostID    int       `json:"host_id" bson:"host_id"`
	Hostname  string    `json:"hostname" bson:"hostname"`
	RuleID    int       `json:"rule_id" bson:"rule_id"`
	Message   string    `json:"message" bson:"message"`
	Value     float64   `json:"value" bson:"value"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Resolved  bool      `json:"resolved" bson:"resolved"`
}
