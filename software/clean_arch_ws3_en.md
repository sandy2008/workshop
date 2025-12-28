# Clean Architecture Workshop: advent-calm-2025

In this workshop, you will learn how to build a robust and testable application using **Go** based on **Clean Architecture**.
The complete project files are in the [advent-of-calm-2025](./advent-of-calm-2025/) directory.

## 1. What is Clean Architecture?

Clean Architecture is a design philosophy that separates concerns, keeping business logic independent of frameworks, databases, and external tools.

### The 3 Layers

This workshop adopts a simple and practical **3-layer structure**.

1. **Domain Layer** - `domain/`
    * **Role**: Core business rules and data structures.
    * **Characteristics**: **Depends on nothing**. Written purely in Go.
    * **Components**: Entities, Repository Interfaces, Domain Services.

2. **Usecase Layer** - `usecase/`
    * **Role**: Application-specific business rules (what the user wants to do).
    * **Characteristics**: Depends only on the Domain Layer. Knows nothing about the DB or HTTP details.
    * **Components**: Interactors (Usecases), Input/Output Data Structures (DTOs).

3. **Infrastructure Layer** - `infra/`
    * **Role**: Detailed technical implementations (DB connections, external API calls, web frameworks).
    * **Characteristics**: **Implements the interfaces defined in the Domain Layer** (dependencies point inward).
    * **Components**: Repository implementations, Web handlers, External clients.

### The Dependency Rule

**"Dependencies always point inwards (towards the Domain)."**
Source code dependencies must always point from lower-level details (concrete implementations) to higher-level abstractions.

```mermaid
graph TD
    %% External / Frameworks & Drivers
    Customer[Customer]
    Admin[Admin]

    subgraph Gateway [Web / API / Gateway]
        OrderAPI[Order API Endpoint]
        InvAPI[Inventory API Endpoint]
    end

    subgraph UsecaseLayer [Usecase]
        OrderUC[CreateOrder UC]
        InvUC[Check/Update UC]
    end

    subgraph DomainLayer [Domain]
        Entities[Order/Inventory Entities]
        Ports["Ports (Interfaces)"]
        OrderDS[Order Domain Svc]
        InvDS[Inventory Domain Svc]
    end

    subgraph InfraLayer [Infra / Adapters]
        OrderRepoImpl[Order Repository Impl]
        InvRepoImpl[Inventory Repository Impl]
        InvClientImpl[Inventory REST Client]
        OrderDB[(Order DB)]
        InvDB[(Inventory DB)]
    end

    %% External Access
    Customer --> OrderAPI
    Admin --> InvAPI

    %% API to Usecase
    OrderAPI --> OrderUC
    InvAPI --> InvUC

    %% Usecase to Domain Dependency
    OrderUC --> OrderDS
    OrderUC --> Ports
    InvUC --> InvDS
    InvUC --> Ports

    %% Domain Service to Interface Dependency
    OrderDS --> Ports
    InvDS --> Ports

    %% Dependency Inversion (DIP)
    OrderRepoImpl -- "implements" --> Ports
    InvRepoImpl -- "implements" --> Ports
    InvClientImpl -- "implements" --> Ports

    %% Implementation to External Resources
    OrderRepoImpl --> OrderDB
    InvRepoImpl --> InvDB

    %% Service Integration & Admin Flow Unification
    %% Both Order Service and Admin use the same Inventory API.
    InvClientImpl --> InvAPI
    InvAPI --> InvUC

    style DomainLayer fill:#f9f,stroke:#333,stroke-width:2px
    style UsecaseLayer fill:#bbf,stroke:#333,stroke-width:2px
    style InfraLayer fill:#bfb,stroke:#333,stroke-width:2px
    style Gateway fill:#fff,stroke:#333,stroke-dasharray: 5 5
```

> **Note: Unifying External Interfaces**
> `Customer` (the person ordering) and `Admin` (inventory manager) interact with the system via the appropriate API endpoints in the `Gateway` layer. Furthermore, the `Inventory REST Client` within the `Order Service` uses the same `Inventory API` as the `Admin`, centralizing all inventory-related logic within the `Inventory Usecase`.

> **What are "Ports"?**
> Ports are the "contracts (interfaces) that the inner rules demand from the outside." Details about the DB or external APIs are hidden behind Ports. Usecases and Domain Services depend on Ports to define behavior only. The outside layer (Infra) implements these Ports, keeping the dependency direction pointing inward.

---

## Workshop: Building an Order System

We will implement a fictional "Order Creation System," starting from the inside and working our way out.

### Step 1: Designing the Domain Layer (`domain/`)

The Domain Layer is the **heart** of the application and consists of the following three elements. These do not depend on any external (DB or Web) concerns.

1. **Entity**: Business data and rules (e.g., `Order`, `Inventory`).
2. **Interface**: Contracts for data persistence or external integration (e.g., `OrderRepository`, `InventoryClient`).
3. **Domain Service**: Logic that spans multiple entities (e.g., `OrderDomainService`).

**Domain Service Rule Examples**

* `OrderDomainService`: Error if `ProductID` is empty, or if quantity is 0 or less.
* `InventoryDomainService`: Error if `ProductID` is empty, or if stock quantity is negative.

First, we define the core business object, the "Order," and the "Interfaces" used to interact with the outside world.

**1. Define Entities (`domain/entity/models.go`)**
Define the state and structure of an Order.

```go
type Order struct {
    ID         string
    CustomerID string
    Amount     float64
    Status     OrderStatus
    CreatedAt  time.Time
}
```

**2. Define Interfaces (Ports) (`domain/repository/interfaces.go`)**
**Abstract** how data is saved or how external services are accessed. The implementation of these interfaces will be done in Step 3.

```go
// Dependency Inversion Principle (DIP): High-level modules own the abstractions.
type OrderRepository interface {
    Save(ctx context.Context, order *entity.Order) error
    FindByID(ctx context.Context, id string) (*entity.Order, error)
}

type InventoryClient interface {
    CheckAndReserve(ctx context.Context, productID string, quantity int) (bool, error)
}

type PaymentPublisher interface {
    PublishPaymentTask(ctx context.Context, order *entity.Order) error
}
```

> **Note: Handling context.Context**
> In Go, it is common to pass `context.Context` to Ports. However, some designs prefer to **keep it within the Usecase layer** to prioritize purity. Understand this as a trade-off depending on the application's needs.

### Step 2: Implementing the Usecase Layer (`usecase/`)

Combine the Domain Layer components to implement the application feature: "Create an Order."

**Implementation (`usecase/create_order.go`)**

```go
type CreateOrderUsecase struct {
    orderRepo repository.OrderRepository // Depends on abstraction
    // ...
}

func (u *CreateOrderUsecase) Execute(ctx context.Context, input CreateOrderInput) error {
    // 1. Check stock (using Domain Service)
    // 2. Create Order entity
    // 3. Save to database (using Repository)
    // 4. Publish event
}
```

The key point here is that `CreateOrderUsecase` does not know about the concrete database (e.g., Postgres). It only knows the Ports (e.g., `OrderRepository`).

> **Note: Transactional Consistency (DB Save vs. MQ Publish)**
> In this example, "DB Save -> MQ Publish" are executed sequentially. In real-world systems, you should consider transaction boundaries and compensation (e.g., the Outbox pattern) to prevent double-sends or missing messages.

### Step 3: Implementing the Infrastructure Layer (`infra/`)

This is where concrete technologies like "PostgreSQL" or "REST API" appear. **We implement the Domain Layer interfaces defined in Step 1**.

* `PostgresOrderRepository` implements `domain.OrderRepository`.
* `RestInventoryClient` implements `domain.InventoryClient`.
* `RabbitMQPaymentPublisher` implements `domain.PaymentPublisher`.

**Repository Implementation (`infra/repository/postgres_order_repository.go`)**

```go
type PostgresOrderRepository struct {
    // DB connection instance, etc.
}

// Satisfies the domain/repository.OrderRepository interface
func (r *PostgresOrderRepository) Save(ctx context.Context, order *entity.Order) error {
    fmt.Printf("Saving order %s to Postgres\n", order.ID)
    // Actual SQL execution logic...
    return nil
}
```

### Step 4: Assembling the Application (`main.go`)

Finally, we wire up all the parts in `main.go` using **Dependency Injection**.

```go
func main() {
    // 1. Create Infrastructure objects
    orderRepo := &repository.PostgresOrderRepository{}
    inventoryClient := &client.RestInventoryClient{}
    paymentPub := &messaging.RabbitMQPaymentPublisher{}
    idGen := &util.UUIDGenerator{} // ID Generator implementation

    // 2. Create Domain Service
    orderDomainSvc := service.NewOrderDomainService(inventoryClient)
    inventoryRepo := &repository.PostgresInventoryRepository{}
    inventoryDomainSvc := service.NewInventoryDomainService(inventoryRepo)

    // 3. Inject into Usecase
    createOrderUsecase := usecase.NewCreateOrderUsecase(orderRepo, orderDomainSvc, paymentPub, idGen)
    checkInventoryUsecase := usecase.NewCheckInventoryUsecase(inventoryDomainSvc)
    updateInventoryUsecase := usecase.NewUpdateInventoryUsecase(inventoryDomainSvc)

    // 4. Run
    createOrderUsecase.Execute(ctx, input)
    checkInventoryUsecase.Execute(ctx, checkInput)
}
```

---

## Design Analysis & Quality (Clean Architecture Analysis)

This project maintains high design quality based on the following principles:

1. **Loose Coupling**: Order and Inventory concerns are separated at the domain level, making it easy to split into microservices in the future.
2. **Pure Business Logic**: The `domain` package has zero dependencies on external libraries, containing only business rules.
3. **Extensibility**: Adding new notification methods (Email, Slack, etc.) only requires adding an interface to `domain/repository` and implementing it in `infra`.

---

## How to Run

Execute the following commands in the project root directory to resolve dependencies and run the application.

```bash
# Tidy up dependencies
go mod tidy

# Run the application
go run main.go
```

If successful, the infrastructure implementation will be called, and logs (like simulated save operations) will be output.

## Summary

* **Robust against change**: Even if you switch the DB to MySQL, the code in `domain` and `usecase` does not change by a single line.
* **Easy to Test**: You only need to mock the `repository` to test the `usecase`. No database is required.
* **Separation of Concerns**: Business logic and technical details are clearly separated.
