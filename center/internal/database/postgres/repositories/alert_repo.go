package repositories

import (
	"center/internal/models"
	"context"
	"database/sql"
	"errors"
)

type PostgresAlertRepository struct {
	db *sql.DB
}

func NewPostgresAlertRepository(db *sql.DB) *PostgresAlertRepository {
	return &PostgresAlertRepository{db: db}
}

func (r *PostgresAlertRepository) GetByHostID(ctx context.Context, hostID int) ([]models.AlertRule, error) {
	const query = `
		SELECT id, host_id, metric_name, threshold_value, condition, enabled
		FROM alert_rules
		WHERE host_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, hostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []models.AlertRule
	for rows.Next() {
		var alert models.AlertRule
		if err := rows.Scan(
			&alert.ID,
			&alert.HostID,
			&alert.MetricName,
			&alert.ThresholdValue,
			&alert.Condition,
			&alert.Enabled,
		); err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return alerts, nil
}

func (r *PostgresAlertRepository) GetByID(ctx context.Context, id int) (*models.AlertRule, error) {
	const query = `
		SELECT id, host_id, metric_name, threshold_value, condition, enabled
		FROM alert_rules
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var alert models.AlertRule
	err := row.Scan(
		&alert.ID,
		&alert.HostID,
		&alert.MetricName,
		&alert.ThresholdValue,
		&alert.Condition,
		&alert.Enabled,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &alert, nil
}

func (r *PostgresAlertRepository) Create(ctx context.Context, alert *models.AlertRule) (int, error) {
	const query = `
		INSERT INTO alert_rules (host_id, metric_name, threshold_value, condition, enabled)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id int
	err := r.db.QueryRowContext(ctx, query,
		alert.HostID,
		alert.MetricName,
		alert.ThresholdValue,
		alert.Condition,
		alert.Enabled,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostgresAlertRepository) Update(ctx context.Context, alert *models.AlertRule) error {
	const query = `
		UPDATE alert_rules
		SET 
			host_id = $2,
			metric_name = $3,
			threshold_value = $4,
			condition = $5,
			enabled = $6
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		alert.ID,
		alert.HostID,
		alert.MetricName,
		alert.ThresholdValue,
		alert.Condition,
		alert.Enabled,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PostgresAlertRepository) Delete(ctx context.Context, id int) error {
	const query = `DELETE FROM alert_rules WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PostgresAlertRepository) SetEnabled(ctx context.Context, id int, enabled bool) error {
	const query = `UPDATE alert_rules SET enabled = $2 WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id, enabled)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PostgresAlertRepository) GetActive(ctx context.Context) ([]models.AlertRule, error) {
	const query = `
		SELECT id, host_id, metric_name, threshold_value, condition, enabled
		FROM alert_rules
		WHERE enabled = true
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []models.AlertRule
	for rows.Next() {
		var alert models.AlertRule
		if err := rows.Scan(
			&alert.ID,
			&alert.HostID,
			&alert.MetricName,
			&alert.ThresholdValue,
			&alert.Condition,
			&alert.Enabled,
		); err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return alerts, nil
}

/*
package services

import (
	"context"
	"monitoring-center/internal/models"
	"monitoring-center/internal/repositories"
)

type AlertService struct {
	alertRepo repositories.AlertRepository
}

func NewAlertService(alertRepo repositories.AlertRepository) *AlertService {
	return &AlertService{alertRepo: alertRepo}
}

func (s *AlertService) CreateAlertRule(ctx context.Context, rule *models.AlertRule) error {
	_, err := s.alertRepo.Create(ctx, rule)
	return err
}

func (s *AlertService) CheckThresholds(ctx context.Context, metrics models.HostMetrics) ([]models.Alert, error) {
	rules, err := s.alertRepo.GetByHostID(ctx, metrics.HostID)
	if err != nil {
		return nil, err
	}

	var triggeredAlerts []models.Alert
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		var value float64
		switch rule.MetricName {
		case "cpu_usage_percent":
			value = metrics.CPU
		case "memory_usage_percent":
			value = metrics.RAM
		case "disk_usage_percent":
			value = metrics.Disk
		default:
			continue
		}

		if rule.CheckCondition(value) {
			triggeredAlerts = append(triggeredAlerts, models.Alert{
				RuleID:      rule.ID,
				HostID:      metrics.HostID,
				MetricName:  rule.MetricName,
				MetricValue: value,
				Threshold:   rule.ThresholdValue,
				Condition:   rule.Condition,
				TriggeredAt: time.Now(),
			})
		}
	}

	return triggeredAlerts, nil
}
*/
