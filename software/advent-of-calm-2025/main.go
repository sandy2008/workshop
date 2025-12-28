package main

import (
	"context"

	"github.com/sokoide/advent-of-calm-2025/cleanarch/domain/service"
	"github.com/sokoide/advent-of-calm-2025/cleanarch/infra/client"
	"github.com/sokoide/advent-of-calm-2025/cleanarch/infra/messaging"
	"github.com/sokoide/advent-of-calm-2025/cleanarch/infra/repository"
	"github.com/sokoide/advent-of-calm-2025/cleanarch/infra/util"
	"github.com/sokoide/advent-of-calm-2025/cleanarch/usecase"
)

func main() {
	// 1. Setup Infrastructure
	orderRepo := &repository.PostgresOrderRepository{}
	inventoryClient := &client.RestInventoryClient{}
	paymentPub := &messaging.RabbitMQPaymentPublisher{}
	idGen := &util.UUIDGenerator{}

	// 2. Setup Domain Service
	orderDomainSvc := service.NewOrderDomainService(inventoryClient)

	// 3. Setup Usecase
	createOrderUsecase := usecase.NewCreateOrderUsecase(orderRepo, orderDomainSvc, paymentPub, idGen)

	// 4. Run Usecase
	ctx := context.Background()
	input := usecase.CreateOrderInput{
		CustomerID: "customer-123",
		ProductID:  "product-456",
		Quantity:   1,
		Amount:     99.99,
	}

	err := createOrderUsecase.Execute(ctx, input)
	if err != nil {
		panic(err)
	}
}
