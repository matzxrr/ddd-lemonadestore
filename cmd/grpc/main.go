package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    
    // Application layer imports
    customerCmds "github.com/matzxrr/ddd-lemonadestore/internal/application/customer/commands"
    customerQueries "github.com/matzxrr/ddd-lemonadestore/internal/application/customer/queries"
    orderCmds "github.com/matzxrr/ddd-lemonadestore/internal/application/order/commands"
    orderHandlers "github.com/matzxrr/ddd-lemonadestore/internal/application/order/event_handlers"
    orderQueries "github.com/matzxrr/ddd-lemonadestore/internal/application/order/queries"
    storeCmds "github.com/matzxrr/ddd-lemonadestore/internal/application/store/commands"
    storeQueries "github.com/matzxrr/ddd-lemonadestore/internal/application/store/queries"
    
    // Domain imports
    "github.com/matzxrr/ddd-lemonadestore/internal/domain/shared"
    "github.com/matzxrr/ddd-lemonadestore/internal/domain/store"
    
    // Infrastructure imports
    "github.com/matzxrr/ddd-lemonadestore/internal/infrastructure/events"
    "github.com/matzxrr/ddd-lemonadestore/internal/infrastructure/persistence/memory"
    
    // Interface imports
    grpcServer "github.com/matzxrr/ddd-lemonadestore/internal/interfaces/grpc"
    "github.com/matzxrr/ddd-lemonadestore/internal/interfaces/grpc/services"
)

func main() {
    // Initialize infrastructure
    // WHY: Creates all the technical implementations needed by the application
    
    // 1. Create repositories (in-memory for demo)
    storeRepo := memory.NewInMemoryStoreRepository()
    orderRepo := memory.NewInMemoryOrderRepository()
    customerRepo := memory.NewInMemoryCustomerRepository()
    
    // 2. Create unit of work
    uow := memory.NewInMemoryUnitOfWork(storeRepo, orderRepo, customerRepo)
    
    // 3. Create event bus
    eventBus := events.NewInMemoryEventBus()
    
    // Initialize application layer
    // WHAT: Create all command and query handlers
    
    // Store handlers
    addInventoryHandler := storeCmds.NewAddInventoryHandler(storeRepo, eventBus)
    updatePriceHandler := storeCmds.NewUpdatePriceHandler(storeRepo, eventBus)
    getProductHandler := storeQueries.NewGetProductHandler(storeRepo)
    getInventoryHandler := storeQueries.NewGetInventoryHandler(storeRepo)
    
    // Order handlers
    createOrderHandler := orderCmds.NewCreateOrderHandler(uow, eventBus)
    cancelOrderHandler := orderCmds.NewCancelOrderHandler(uow, eventBus)
    getOrderHandler := orderQueries.NewGetOrderHandler(orderRepo)
    listOrdersHandler := orderQueries.NewListOrdersHandler(orderRepo)
    
    // Customer handlers
    registerCustomerHandler := customerCmds.NewRegisterCustomerHandler(customerRepo, eventBus)
    updateCustomerHandler := customerCmds.NewUpdateCustomerHandler(customerRepo, eventBus)
    getCustomerHandler := customerQueries.NewGetCustomerHandler(customerRepo)
    
    // Register event handlers
    // WHY: Implements eventual consistency between aggregates
    orderPlacedHandler := orderHandlers.NewOrderPlacedHandler(customerRepo)
    eventBus.Subscribe("order.confirmed", orderPlacedHandler.Handle)
    
    // Initialize sample data
    initializeSampleData(storeRepo)
    
    // Initialize presentation layer
    // WHERE: Create gRPC services that expose application functionality
    
    storeService := services.NewStoreService(
        addInventoryHandler,
        updatePriceHandler,
        getProductHandler,
        getInventoryHandler,
    )
    
    orderService := services.NewOrderService(
        createOrderHandler,
        cancelOrderHandler,
        getOrderHandler,
        listOrdersHandler,
    )
    
    customerService := services.NewCustomerService(
        registerCustomerHandler,
        updateCustomerHandler,
        getCustomerHandler,
    )
    
    // Create and start gRPC server
    server := grpcServer.NewServer(storeService, orderService, customerService)
    
    // Handle graceful shutdown
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
        <-sigChan
        
        server.Stop()
    }()
    
    // Start server
    if err := server.Start(":50051"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

// initializeSampleData creates initial store and products
// WHY: Provides data for testing the application
func initializeSampleData(storeRepo store.StoreRepository) {
    // Create address
    address, _ := shared.NewAddress(
        "123 Main St",
        "Lemonade City",
        "CA",
        "12345",
        "USA",
    )
    
    // Create store
    mainStore, _ := store.NewStore("Main Street Lemonade Stand", address)
    
    // Add products
    // Classic Lemonade
    classicPrice, _ := shared.NewMoney(299, "USD") // $2.99
    mainStore.AddProduct(
        "Classic Lemonade",
        "Our traditional lemonade made with fresh lemons",
        classicPrice,
    )
    
    // Strawberry Lemonade
    strawberryPrice, _ := shared.NewMoney(349, "USD") // $3.49
    product2, _ := mainStore.AddProduct(
        "Strawberry Lemonade",
        "Sweet strawberry mixed with tart lemonade",
        strawberryPrice,
    )
    
    // Pink Lemonade
    pinkPrice, _ := shared.NewMoney(329, "USD") // $3.29
    product3, _ := mainStore.AddProduct(
        "Pink Lemonade",
        "A fun twist on classic lemonade",
        pinkPrice,
    )
    
    // Add initial inventory
    products := mainStore.Products
    for productID := range products {
        mainStore.AddInventory(productID, 100) // 100 units each
    }
    
    // Save store
    storeRepo.Save(mainStore)
    
    log.Printf("Initialized store with ID: %s", mainStore.ID())
}
