package service

import (
	"context"
	"github.com/sokoide/advent-of-calm-2025/cleanarch/domain/repository"
)

type InventoryDomainService struct {
	repo repository.InventoryRepository
}

func NewInventoryDomainService(repo repository.InventoryRepository) *InventoryDomainService {
	return &InventoryDomainService{repo: repo}
}

func (s *InventoryDomainService) GetStock(ctx context.Context, productID string) (int, error) {
	return s.repo.GetStock(ctx, productID)
}

func (s *InventoryDomainService) UpdateStock(ctx context.Context, productID string, quantity int) error {
	// 業務ルール（例：負の在庫を許可しない等）をここに実装可能
	return s.repo.UpdateStock(ctx, productID, quantity)
}
