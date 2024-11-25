package repository

import (
	"context"
	"ecommerce/user-service/models"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByPhone(ctx context.Context, phone string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, id string, user *models.User) error
}
