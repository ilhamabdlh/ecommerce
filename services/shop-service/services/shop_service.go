package services

import (
	"errors"

	"ecommerce/shop-service/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type shopService struct {
	shopRepo models.ShopRepository
}

func NewShopService(shopRepo models.ShopRepository) models.ShopService {
	return &shopService{
		shopRepo: shopRepo,
	}
}

func (s *shopService) CreateShop(shop *models.Shop) error {
	if shop.Name == "" {
		return errors.New("shop name is required")
	}
	if shop.Location == "" {
		return errors.New("shop location is required")
	}

	return s.shopRepo.Create(shop)
}

func (s *shopService) GetShop(id primitive.ObjectID) (*models.Shop, error) {
	shop, err := s.shopRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return shop, nil
}

func (s *shopService) ListShops() ([]models.Shop, error) {
	return s.shopRepo.GetAll()
}

func (s *shopService) UpdateShop(shop *models.Shop) error {
	if shop.ID.IsZero() {
		return errors.New("shop ID is required")
	}
	if shop.Name == "" {
		return errors.New("shop name is required")
	}
	if shop.Location == "" {
		return errors.New("shop location is required")
	}

	return s.shopRepo.Update(shop)
}

func (s *shopService) DeleteShop(id primitive.ObjectID) error {
	return s.shopRepo.Delete(id)
}

func (s *shopService) AddWarehouseToShop(shopID, warehouseID primitive.ObjectID) error {
	shop, err := s.shopRepo.GetByID(shopID)
	if err != nil {
		return err
	}

	// Check if warehouse already exists in shop
	for _, wID := range shop.Warehouses {
		if wID == warehouseID {
			return errors.New("warehouse already associated with this shop")
		}
	}

	return s.shopRepo.AddWarehouse(shopID, warehouseID)
}

func (s *shopService) RemoveWarehouseFromShop(shopID, warehouseID primitive.ObjectID) error {
	shop, err := s.shopRepo.GetByID(shopID)
	if err != nil {
		return err
	}

	// Check if warehouse exists in shop
	warehouseExists := false
	for _, wID := range shop.Warehouses {
		if wID == warehouseID {
			warehouseExists = true
			break
		}
	}

	if !warehouseExists {
		return errors.New("warehouse not associated with this shop")
	}

	return s.shopRepo.RemoveWarehouse(shopID, warehouseID)
}

func (s *shopService) ActivateShop(id primitive.ObjectID) error {
	return s.shopRepo.UpdateStatus(id, "active")
}

func (s *shopService) DeactivateShop(id primitive.ObjectID) error {
	return s.shopRepo.UpdateStatus(id, "inactive")
}
