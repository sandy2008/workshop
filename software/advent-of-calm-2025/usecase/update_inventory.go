package usecase

import (
	"context"
	"github.com/sokoide/advent-of-calm-2025/cleanarch/domain/service"
)

type UpdateInventoryInput struct {
	ProductID string
	Quantity  int
}

type UpdateInventoryUsecase struct {
	inventoryService *service.InventoryDomainService
}

func NewUpdateInventoryUsecase(svc *service.InventoryDomainService) *UpdateInventoryUsecase {
	return &UpdateInventoryUsecase{inventoryService: svc}
}

func (u *UpdateInventoryUsecase) Execute(ctx context.Context, input UpdateInventoryInput) error {
	return u.inventoryService.UpdateStock(ctx, input.ProductID, input.Quantity)
}
