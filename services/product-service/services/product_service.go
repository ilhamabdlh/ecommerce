package services

import (
	"errors"

	"ecommerce/product-service/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type productService struct {
	productRepo models.ProductRepository
}

func NewProductService(productRepo models.ProductRepository) models.ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

func (s *productService) CreateProduct(product *models.Product) error {
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}
	if product.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}

	return s.productRepo.Create(product)
}

func (s *productService) GetProduct(id primitive.ObjectID) (*models.Product, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) ListProducts() ([]models.Product, error) {
	return s.productRepo.GetAll()
}

func (s *productService) UpdateProduct(product *models.Product) error {
	if product.ID.IsZero() {
		return errors.New("product ID is required")
	}
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}
	if product.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}

	return s.productRepo.Update(product)
}

func (s *productService) DeleteProduct(id primitive.ObjectID) error {
	return s.productRepo.Delete(id)
}

func (s *productService) UpdateProductStock(id primitive.ObjectID, quantity int) error {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return err
	}

	if product.Stock+quantity < 0 {
		return errors.New("insufficient stock")
	}

	return s.productRepo.UpdateStock(id, quantity)
}
