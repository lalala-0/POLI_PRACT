package postgres

import (
	"POLI_PRACT/models"
	"context"
	"database/sql"
)

// PostgresHostRepository реализация репозитория хостов
type PostgresContainerRepository struct {
	db *sql.DB
}

func (p PostgresContainerRepository) GetAll(ctx context.Context) ([]models.Host, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresContainerRepository) GetByID(ctx context.Context, id int) (*models.Host, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresContainerRepository) Create(ctx context.Context, host *models.Host) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresContainerRepository) Update(ctx context.Context, host *models.Host) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresContainerRepository) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresContainerRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresContainerRepository) GetMaster(ctx context.Context) (*models.Host, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresContainerRepository) SetMaster(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func NewContainerRepository(db *sql.DB) HostRepository {
	return &PostgresContainerRepository{db: db}
}
