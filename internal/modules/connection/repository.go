package connection

import (
	"context"

	"github.com/khanalsaroj/typegen-server/internal/domain"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *domain.DatabaseConnection) error
	FindByID(ctx context.Context, id uint) (*domain.DatabaseConnection, error)
	FindAll(ctx context.Context, offset, limit int) ([]*domain.DatabaseConnection, int64, error)
	Update(ctx context.Context, DatabaseConnection *domain.DatabaseConnection) error
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, connection *domain.DatabaseConnection) error {
	return r.db.WithContext(ctx).Create(connection).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (*domain.DatabaseConnection, error) {
	var connection domain.DatabaseConnection
	if err := r.db.WithContext(ctx).First(&connection, id).Error; err != nil {
		return nil, err
	}
	return &connection, nil
}

func (r *repository) FindAll(ctx context.Context, offset, limit int) ([]*domain.DatabaseConnection, int64, error) {
	var users []*domain.DatabaseConnection
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.DatabaseConnection{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *repository) Update(ctx context.Context, connection *domain.DatabaseConnection) error {
	return r.db.WithContext(ctx).Save(connection).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.DatabaseConnection{}, id).Error
}
