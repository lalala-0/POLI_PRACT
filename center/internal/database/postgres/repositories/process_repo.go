package repositories

import (
	"center/internal/models"
	"context"
	"database/sql"
	"errors"
)

type PostgresProcessRepository struct {
	db *sql.DB
}

func NewPostgresProcessRepository(db *sql.DB) *PostgresProcessRepository {
	return &PostgresProcessRepository{db: db}
}

// GetByHostID возвращает все процессы для указанного хоста
func (p *PostgresProcessRepository) GetByHostID(ctx context.Context, hostID int) ([]models.Process, error) {
	const query = `SELECT id, host_id, process_name FROM host_processes WHERE host_id = $1`

	rows, err := p.db.QueryContext(ctx, query, hostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var processes []models.Process
	for rows.Next() {
		var proc models.Process
		if err := rows.Scan(&proc.ID, &proc.HostID, &proc.ProcessName); err != nil {
			return nil, err
		}
		processes = append(processes, proc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return processes, nil
}

// GetByID возвращает процесс по его идентификатору
func (p *PostgresProcessRepository) GetByID(ctx context.Context, id int) (*models.Process, error) {
	const query = `SELECT id, host_id, process_name FROM host_processes WHERE id = $1`

	row := p.db.QueryRowContext(ctx, query, id)

	var proc models.Process
	err := row.Scan(&proc.ID, &proc.HostID, &proc.ProcessName)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return &proc, nil
	}
}

// Create создает новый процесс и возвращает его ID
func (p *PostgresProcessRepository) Create(ctx context.Context, process *models.Process) (int, error) {
	const query = `INSERT INTO host_processes (host_id, process_name) 
				   VALUES ($1, $2) 
				   RETURNING id`

	var id int
	err := p.db.QueryRowContext(
		ctx,
		query,
		process.HostID,
		process.ProcessName,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update обновляет данные процесса
func (p *PostgresProcessRepository) Update(ctx context.Context, process *models.Process) error {
	const query = `
		UPDATE host_processes 
		SET host_id = $1, 
			process_name = $2 
		WHERE id = $3
	`

	result, err := p.db.ExecContext(
		ctx,
		query,
		process.HostID,
		process.ProcessName,
		process.ID,
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

// Delete удаляет процесс по идентификатору
func (p *PostgresProcessRepository) Delete(ctx context.Context, id int) error {
	const query = `DELETE FROM host_processes WHERE id = $1`

	result, err := p.db.ExecContext(ctx, query, id)
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

// Exists проверяет существование процесса для хоста
func (p *PostgresProcessRepository) Exists(ctx context.Context, hostID int, processName string) (bool, error) {
	const query = `
		SELECT EXISTS(
			SELECT 1 
			FROM host_processes 
			WHERE host_id = $1 AND process_name = $2
		)
	`

	var exists bool
	err := p.db.QueryRowContext(ctx, query, hostID, processName).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

/*
func (s *HostService) AddProcessToHost(ctx context.Context, hostID int, name string) error {
    // Проверка существования процесса
    exists, err := s.processRepo.Exists(ctx, hostID, name)
    if err != nil {
        return err
    }
    if exists {
        return errors.New("process already monitored")
    }

    // Создание нового процесса
    _, err = s.processRepo.Create(ctx, &models.Process{
        HostID:      hostID,
        ProcessName: name,
    })
    return err
}

*/
