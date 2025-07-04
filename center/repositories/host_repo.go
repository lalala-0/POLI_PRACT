package postgres

import (
	"POLI_PRACT/models"
	"context"
	"database/sql"
	"time"
)

// PostgresHostRepository реализация репозитория хостов
type PostgresHostRepository struct {
	db *sql.DB
}

func (r *PostgresHostRepository) GetByID(ctx context.Context, id int) (*models.Host, error) {
	query := `SELECT id, hostname, ip_address, priority, is_master, status, last_check, created_at, updated_at FROM hosts WHERE id = $1`

	var host models.Host
	err := r.db.QueryRowContext(ctx, query, id).Scan(&host.ID, &host.Hostname, &host.IPAddress,
		&host.Priority, &host.IsMaster, &host.Status, &host.LastCheck, &host.CreatedAt, &host.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &host, nil
}

func (r *PostgresHostRepository) Create(ctx context.Context, host *models.Host) (int, error) {
	query := `INSERT INTO hosts (hostname, ip_address, priority, is_master, status, last_check)
              VALUES ($1, $2, $3, $4, $5, $6)
              RETURNING id, created_at, updated_at`

	// Установка значений по умолчанию, если они не заданы
	if host.Status == "" {
		host.Status = "unknown"
	}

	var id int
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(
		ctx,
		query,
		host.Hostname,
		host.IPAddress,
		host.Priority,
		host.IsMaster,
		host.Status,
		host.LastCheck,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return 0, err
	}

	// Обновляем структуру host с полученными значениями
	host.ID = id
	host.CreatedAt = createdAt
	host.UpdatedAt = updatedAt

	return id, nil
}

func (r *PostgresHostRepository) Update(ctx context.Context, host *models.Host) error {
	query := `UPDATE hosts
			  SET hostname = $1,
				  ip_address = $2,
				  priority = $3,
				  is_master = $4,
				  status = $5,
				  last_check = $6,
				  updated_at = NOW()
			  WHERE id = $7`

	result, err := r.db.ExecContext(
		ctx,
		query,
		host.Hostname,
		host.IPAddress,
		host.Priority,
		host.IsMaster,
		host.Status,
		host.LastCheck,
		host.ID,
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

	// Обновляем updated_at в структуре host
	var updatedAt time.Time
	err = r.db.QueryRowContext(ctx, `SELECT updated_at FROM hosts WHERE id = $1`, host.ID).Scan(&updatedAt)
	if err != nil {
		return err
	}
	host.UpdatedAt = updatedAt

	return nil
}

func (r *PostgresHostRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM hosts WHERE id = $1`

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

func (r *PostgresHostRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE hosts SET status = $1, last_check = NOW(), updated_at = NOW() WHERE id = $2`

	result, err := r.db.ExecContext(
		ctx,
		query,
		status,
		id,
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

// NewHostRepository создает новый репозиторий хостов
func NewHostRepository(db *sql.DB) HostRepository {
	return &PostgresHostRepository{db: db}
}

// GetAll возвращает все хосты
func (r *PostgresHostRepository) GetAll(ctx context.Context) ([]models.Host, error) {
	query := `SELECT id, hostname, ip_address, priority, is_master, status, last_check, 
                  created_at, updated_at FROM hosts ORDER BY priority DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hosts []models.Host
	for rows.Next() {
		var h models.Host
		if err := rows.Scan(&h.ID, &h.Hostname, &h.IPAddress, &h.Priority, &h.IsMaster,
			&h.Status, &h.LastCheck, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, err
		}
		hosts = append(hosts, h)
	}

	return hosts, nil
}

// GetMaster возвращает текущий мастер-хост
func (r *PostgresHostRepository) GetMaster(ctx context.Context) (*models.Host, error) {
	query := `SELECT id, hostname, ip_address, priority, is_master, status, last_check, 
                  created_at, updated_at FROM hosts WHERE is_master = true LIMIT 1`

	var host models.Host
	err := r.db.QueryRowContext(ctx, query).Scan(&host.ID, &host.Hostname, &host.IPAddress,
		&host.Priority, &host.IsMaster, &host.Status, &host.LastCheck, &host.CreatedAt, &host.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &host, nil
}

// SetMaster устанавливает новый мастер-хост
func (r *PostgresHostRepository) SetMaster(ctx context.Context, id int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Сначала сбрасываем флаг is_master у всех хостов
	_, err = tx.ExecContext(ctx, `UPDATE hosts SET is_master = false`)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Устанавливаем флаг is_master для выбранного хоста
	_, err = tx.ExecContext(ctx, `UPDATE hosts SET is_master = true WHERE id = $1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
