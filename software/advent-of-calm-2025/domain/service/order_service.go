package service

import (
	"context"
	"github.com/sokoide/advent-of-calm-2025/cleanarch/domain/entity"
	"github.com/sokoide/advent-of-calm-2025/cleanarch/domain/repository"
)

type OrderDomainService struct {
	inventoryClient repository.InventoryClient
}

func NewOrderDomainService(ic repository.InventoryClient) *OrderDomainService {
	return &OrderDomainService{inventoryClient: ic}
}

func (s *OrderDomainService) ValidateAndReserveStock(ctx context.Context, productID string, quantity int) error {
	if productID == "" {
		return entity.ErrInvalidProductID
	}
	if quantity <= 0 {
		return entity.ErrInvalidQuantity
	}

	available, err := s.inventoryClient.CheckAndReserve(ctx, productID, quantity)
	if err != nil {
		return err
	}
	if !available {
		return entity.ErrInsufficientStock
	}
	return nil
}
