package usecase

import (
	"context"
	"github.com/sokoide/advent-of-calm-2025/cleanarch/domain/service"
)

type CheckInventoryInput struct {
	ProductID string
}

type CheckInventoryOutput struct {
	ProductID string
	Quantity  int
}

type CheckInventoryUsecase struct {
	inventoryService *service.InventoryDomainService
}

func NewCheckInventoryUsecase(svc *service.InventoryDomainService) *CheckInventoryUsecase {
	return &CheckInventoryUsecase{inventoryService: svc}
}

func (u *CheckInventoryUsecase) Execute(ctx context.Context, input CheckInventoryInput) (*CheckInventoryOutput, error) {
	stock, err := u.inventoryService.GetStock(ctx, input.ProductID)
	if err != nil {
		return nil, err
	}

	return &CheckInventoryOutput{
		ProductID: input.ProductID,
		Quantity:  stock,
	}, nil
}
