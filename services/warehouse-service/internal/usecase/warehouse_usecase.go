package usecase

import (
	"context"
	"ecommerce/warehouse-service/models"
	"ecommerce/warehouse-service/repository"
	"errors"
	"time"
)

type warehouseUsecase struct {
	warehouseRepo repository.WarehouseRepository
}

func NewWarehouseUsecase(repo repository.WarehouseRepository) *warehouseUsecase {
	return &warehouseUsecase{
		warehouseRepo: repo,
	}
}

func (u *warehouseUsecase) CreateWarehouse(ctx context.Context, warehouse *models.Warehouse) error {
	if warehouse.Name == "" {
		return errors.New("warehouse name is required")
	}
	if warehouse.Location == "" {
		return errors.New("warehouse location is required")
	}

	warehouse.Status = "active"
	warehouse.Stock = make(map[string]int)
	return u.warehouseRepo.Create(ctx, warehouse)
}

func (u *warehouseUsecase) GetWarehouse(ctx context.Context, id string) (*models.Warehouse, error) {
	return u.warehouseRepo.GetByID(ctx, id)
}

func (u *warehouseUsecase) UpdateStock(ctx context.Context, warehouseID string, productID string, quantity int) error {
	warehouse, err := u.warehouseRepo.GetByID(ctx, warehouseID)
	if err != nil {
		return err
	}

	if warehouse.Status != "active" {
		return errors.New("warehouse is not active")
	}

	currentStock := warehouse.Stock[productID]
	if currentStock+quantity < 0 {
		return errors.New("insufficient stock")
	}

	return u.warehouseRepo.UpdateStock(ctx, warehouseID, productID, quantity)
}

func (u *warehouseUsecase) TransferStock(ctx context.Context, transfer *models.StockTransfer) error {
	if transfer.Quantity <= 0 {
		return errors.New("transfer quantity must be positive")
	}

	// Validate source warehouse
	sourceWarehouse, err := u.warehouseRepo.GetByID(ctx, transfer.FromWarehouse)
	if err != nil {
		return err
	}
	if sourceWarehouse.Status != "active" {
		return errors.New("source warehouse is not active")
	}
	if sourceWarehouse.Stock[transfer.ProductID] < transfer.Quantity {
		return errors.New("insufficient stock in source warehouse")
	}

	// Validate destination warehouse
	destWarehouse, err := u.warehouseRepo.GetByID(ctx, transfer.ToWarehouse)
	if err != nil {
		return err
	}
	if destWarehouse.Status != "active" {
		return errors.New("destination warehouse is not active")
	}

	transfer.Status = "pending"
	transfer.TransferDate = time.Now()

	err = u.warehouseRepo.TransferStock(ctx, transfer)
	if err != nil {
		transfer.Status = "failed"
		return err
	}

	transfer.Status = "completed"
	return nil
}

func (u *warehouseUsecase) ActivateWarehouse(ctx context.Context, warehouseID string) error {
	return u.warehouseRepo.UpdateStatus(ctx, warehouseID, "active")
}

func (u *warehouseUsecase) DeactivateWarehouse(ctx context.Context, warehouseID string) error {
	return u.warehouseRepo.UpdateStatus(ctx, warehouseID, "inactive")
}

func (u *warehouseUsecase) GetAllWarehouses(ctx context.Context) ([]*models.Warehouse, error) {
	return u.warehouseRepo.GetAll(ctx)
}
