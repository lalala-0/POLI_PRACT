package postgres

import (
	"POLI_PRACT/models"
	"context"
	"database/sql"
)

type PostgresProcessRepository struct {
	db *sql.DB
}

func (p PostgresProcessRepository) GetByHostID(ctx context.Context, hostID int) ([]models.Process, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresProcessRepository) GetByID(ctx context.Context, id int) (*models.Process, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresProcessRepository) Create(ctx context.Context, process *models.Process) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresProcessRepository) Update(ctx context.Context, process *models.Process) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresProcessRepository) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresProcessRepository) Exists(ctx context.Context, hostID int, processName string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewProcessRepository(db *sql.DB) ProcessRepository {
	return &PostgresProcessRepository{db: db}
}
