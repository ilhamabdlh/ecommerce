package domain

import (
	"context"
	"ecommerce/warehouse-service/models"
)

type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *models.Warehouse) error
	GetByID(ctx context.Context, id string) (*models.Warehouse, error)
	UpdateStock(ctx context.Context, warehouseID string, productID string, quantity int) error
	TransferStock(ctx context.Context, transfer *models.StockTransfer) error
	UpdateStatus(ctx context.Context, warehouseID string, status string) error
	GetAll(ctx context.Context) ([]*models.Warehouse, error)
}

type WarehouseUsecase interface {
	CreateWarehouse(ctx context.Context, warehouse *models.Warehouse) error
	GetWarehouse(ctx context.Context, id string) (*models.Warehouse, error)
	UpdateStock(ctx context.Context, warehouseID string, productID string, quantity int) error
	TransferStock(ctx context.Context, transfer *models.StockTransfer) error
	ActivateWarehouse(ctx context.Context, warehouseID string) error
	DeactivateWarehouse(ctx context.Context, warehouseID string) error
	GetAllWarehouses(ctx context.Context) ([]*models.Warehouse, error)
}
