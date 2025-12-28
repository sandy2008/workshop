package service

import (
	"context"
	"github.com/sokoide/advent-of-calm-2025/cleanarch/domain/entity"
	"github.com/sokoide/advent-of-calm-2025/cleanarch/domain/repository"
)

type InventoryDomainService struct {
	repo repository.InventoryRepository
}

func NewInventoryDomainService(repo repository.InventoryRepository) *InventoryDomainService {
	return &InventoryDomainService{repo: repo}
}

func (s *InventoryDomainService) GetStock(ctx context.Context, productID string) (int, error) {
	if productID == "" {
		return 0, entity.ErrInvalidProductID
	}
	return s.repo.GetStock(ctx, productID)
}

func (s *InventoryDomainService) UpdateStock(ctx context.Context, productID string, quantity int) error {
	if productID == "" {
		return entity.ErrInvalidProductID
	}
	if quantity < 0 {
		return entity.ErrInvalidQuantity
	}
	return s.repo.UpdateStock(ctx, productID, quantity)
}
