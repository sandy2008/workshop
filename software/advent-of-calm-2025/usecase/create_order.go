package usecase

import (
	"context"
	"time"

	"github.com/sokoide/advent-of-calm-2025/cleanarch/domain/entity"
	"github.com/sokoide/advent-of-calm-2025/cleanarch/domain/repository"
	"github.com/sokoide/advent-of-calm-2025/cleanarch/domain/service"
)

type CreateOrderInput struct {
	CustomerID string
	ProductID  string
	Quantity   int
	Amount     float64
}

type CreateOrderUsecase struct {
	orderRepo    repository.OrderRepository
	orderService *service.OrderDomainService
	paymentPub   repository.PaymentPublisher
	idGen        repository.IDGenerator
}

func NewCreateOrderUsecase(
	repo repository.OrderRepository,
	svc *service.OrderDomainService,
	pub repository.PaymentPublisher,
	idGen repository.IDGenerator,
) *CreateOrderUsecase {
	return &CreateOrderUsecase{
		orderRepo:    repo,
		orderService: svc,
		paymentPub:   pub,
		idGen:        idGen,
	}
}

func (u *CreateOrderUsecase) Execute(ctx context.Context, input CreateOrderInput) error {
	// 1. Domain Service を使用して在庫を確保
	if err := u.orderService.ValidateAndReserveStock(ctx, input.ProductID, input.Quantity); err != nil {
		return err
	}

	// 2. エンティティの作成
	order := &entity.Order{
		ID:         u.idGen.GenerateID(),
		CustomerID: input.CustomerID,
		Amount:     input.Amount,
		Status:     entity.OrderStatusPending,
		CreatedAt:  time.Now(),
	}

	// 3. 永続化
	if err := u.orderRepo.Save(ctx, order); err != nil {
		return err
	}

	// 4. 非同期処理の開始（支払い処理をキューへ）
	return u.paymentPub.PublishPaymentTask(ctx, order)
}
