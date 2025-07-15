package postgres

import (
	"center/internal/models"
	"context"
	"fmt"
	"database/sql"
)

// PostgresContainerRepository реализация репозитория хостов
type PostgresContainerRepository struct {
	db *sql.DB
}

func NewPostgresContainerRepository(db *sql.DB) *PostgresContainerRepository {
	return &PostgresContainerRepository{db: db}
}

func (r *PostgresContainerRepository) GetByHostID(ctx context.Context, hostID int) ([]models.Container, error) {
	const query = `
		SELECT id, host_id, container_name 
		FROM host_containers 
		WHERE host_id = $1
	`
	
	rows, err := r.db.QueryContext(ctx, query, hostID)
	if err != nil {
		return nil, fmt.Errorf("failed to query containers: %w", err)
	}
	defer rows.Close()

	var containers []models.Container
	for rows.Next() {
		var c models.Container
		if err := rows.Scan(&c.ID, &c.HostID, &c.ContainerName); err != nil {
			return nil, fmt.Errorf("failed to scan container row: %w", err)
		}
		containers = append(containers, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return containers, nil
}

func (r *PostgresContainerRepository) GetByID(ctx context.Context, id int) (*models.Container, error) {
	const query = `
		SELECT id, host_id, container_name 
		FROM host_containers 
		WHERE id = $1
	`
	
	var container models.Container
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&container.ID,
		&container.HostID,
		&container.ContainerName,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("failed to get container: %w", err)
	default:
		return &container, nil
	}
}

func (r *PostgresContainerRepository) Create(ctx context.Context, container *models.Container) (int, error) {
	const query = `
		INSERT INTO host_containers (host_id, container_name)
		VALUES ($1, $2)
		RETURNING id
	`
	
	var id int
	err := r.db.QueryRowContext(
		ctx,
		query,
		container.HostID,
		container.ContainerName,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create container: %w", err)
	}

	return id, nil
}

func (r *PostgresContainerRepository) Update(ctx context.Context, container *models.Container) error {
	const query = `
		UPDATE host_containers 
		SET container_name = $1
		WHERE id = $2
	`
	
	result, err := r.db.ExecContext(
		ctx,
		query,
		container.ContainerName,
		container.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update container: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("container with ID %d not found", container.ID)
	}

	return nil
}

func (r *PostgresContainerRepository) Delete(ctx context.Context, id int) error {
	const query = `DELETE FROM host_containers WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete container: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("container with ID %d not found", id)
	}

	return nil
}

func (r *PostgresContainerRepository) Exists(ctx context.Context, hostID int, containerName string) (bool, error) {
	const query = `
		SELECT EXISTS(
			SELECT 1 
			FROM host_containers 
			WHERE host_id = $1 AND container_name = $2
		)
	`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, hostID, containerName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check container existence: %w", err)
	}

	return exists, nil
}


/*
// Инициализация
db, _ := sql.Open("postgres", connectionString)
containerRepo := NewPostgresContainerRepository(db)

// Создание контейнера
newContainer := &models.Container{
    HostID:        1,
    ContainerName: "my-container",
}
id, err := containerRepo.Create(context.Background(), newContainer)

// Получение контейнеров хоста
containers, err := containerRepo.GetByHostID(context.Background(), 1)

// Проверка существования
exists, err := containerRepo.Exists(context.Background(), 1, "my-container")
*/