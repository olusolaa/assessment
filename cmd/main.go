package main

import (
	"fmt"
	"os"

	"assessment/adapter/registry"
	"assessment/adapter/sorter"
	"assessment/domain/model"
	"assessment/infrastructure/config"
	"assessment/infrastructure/persistence"
	"assessment/usecase"
)

func main() {
	// Initialize repository
	repo := persistence.NewInMemoryProductRepository()
	if repo == nil {
		fmt.Println("Failed to initialize repository")
		os.Exit(1)
	}

	// Load sample data
	if err := loadSampleData(repo); err != nil {
		fmt.Printf("Error loading sample data: %v\n", err)
		os.Exit(1)
	}

	// Load configuration
	cfg := config.NewConfig()
	configFile := "infrastructure/config/sample_config.json"

	if _, err := os.Stat(configFile); err == nil {
		if err := cfg.LoadFromFile(configFile); err != nil {
			fmt.Printf("Warning: Failed to load config file: %v\n", err)
		} else {
			fmt.Println("Configuration loaded from", configFile)
		}
	} else {
		fmt.Println("Using default configuration")
	}

	// Initialize sorter registry and use case
	sorterRegistry := registry.NewSorterRegistry()
	sorterUseCase := usecase.NewProductSorterUseCase(sorterRegistry)
	sorterUseCase.SetConfig(cfg)
	sorter.InitializeDefaultSorters(sorterRegistry, cfg)

	// Run the application
	runApp(repo, sorterUseCase)
}

// runApp runs the main application logic
func runApp(repo *persistence.InMemoryProductRepository, sorterUseCase *usecase.ProductSorterUseCase) {
	// Display available sorters
	fmt.Println("Available sorters:")
	for _, name := range sorterUseCase.GetAvailableSorters() {
		fmt.Printf("- %s\n", name)
	}
	fmt.Println()

	// Retrieve products from repository
	products, err := repo.GetAll()
	if err != nil {
		fmt.Printf("Error retrieving products: %v\n", err)
		return
	}

	// Display products sorted by price
	sortedProducts, err := sorterUseCase.SortProducts(products, "Price (ascending)")
	if err != nil {
		fmt.Printf("Error sorting products by price: %v\n", err)
		return
	}

	fmt.Println("Products sorted by Price (ascending):")
	for _, p := range sortedProducts {
		fmt.Println(p.String())
	}
}

// loadSampleData loads sample product data into the repository
func loadSampleData(repo *persistence.InMemoryProductRepository) error {
	// Sample product data
	sampleProducts := []struct {
		ID         int
		Name       string
		Price      float64
		Created    string
		SalesCount int
		ViewsCount int
	}{
		{1, "Alabaster Table", 12.99, "2019-01-04", 32, 730},
		{2, "Zebra Table", 44.49, "2012-01-04", 301, 3279},
		{3, "Coffee Table", 10.00, "2014-05-28", 1048, 20123},
	}

	products := make(model.ProductList, 0, len(sampleProducts))

	for _, data := range sampleProducts {
		created, err := model.ParseTime(data.Created)
		if err != nil {
			return fmt.Errorf("error parsing date %s: %w", data.Created, err)
		}

		product := &model.Product{
			ID:         data.ID,
			Name:       data.Name,
			Price:      data.Price,
			Created:    created,
			SalesCount: data.SalesCount,
			ViewsCount: data.ViewsCount,
		}
		products = append(products, product)
	}

	return repo.Save(products)
}
