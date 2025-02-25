package e2e

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"assessment/adapter/registry"
	"assessment/adapter/sorter"
	"assessment/domain/model"
	"assessment/infrastructure/config"
	"assessment/infrastructure/persistence"
	"assessment/usecase"
)

// TestMainE2E tests the main application flow end-to-end
func TestMainE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	// Capture stdout to verify output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Initialize repository with sample data
	repo := initializeRepository()
	if repo == nil {
		t.Fatal("Failed to initialize repository")
	}

	// Load configuration
	cfg := loadConfiguration()
	if cfg == nil {
		t.Fatal("Failed to load configuration")
	}

	// Initialize sorter registry and use case
	sorterRegistry := registry.NewSorterRegistry()
	sorterUseCase := usecase.NewProductSorterUseCase(sorterRegistry)
	sorterUseCase.SetConfig(cfg)
	sorter.InitializeDefaultSorters(sorterRegistry, cfg)

	// Display available sorters
	fmt.Println("Available sorters:")
	for _, name := range sorterUseCase.GetAvailableSorters() {
		fmt.Printf("- %s\n", name)
	}
	fmt.Println()

	// Retrieve products from repository
	products, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Error retrieving products: %v", err)
	}

	// Display products sorted by price
	displaySortedProducts(sorterUseCase, products, "Price (ascending)")

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("Failed to copy from pipe: %v", err)
	}
	output := buf.String()

	// Verify output contains expected strings
	expectedStrings := []string{
		"Available sorters:",
		"Products sorted by Price (ascending):",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Output missing expected string: %s", expected)
		}
	}
}

// initializeRepository creates and populates a repository with sample data
func initializeRepository() *persistence.InMemoryProductRepository {
	// Sample product data for demonstration
	sampleProducts := []map[string]interface{}{
		{
			"id":          1,
			"name":        "Alabaster Table",
			"price":       12.99,
			"created":     "2019-01-04",
			"sales_count": 32,
			"views_count": 730,
		},
		{
			"id":          2,
			"name":        "Zebra Table",
			"price":       44.49,
			"created":     "2012-01-04",
			"sales_count": 301,
			"views_count": 3279,
		},
		{
			"id":          3,
			"name":        "Coffee Table",
			"price":       10.00,
			"created":     "2014-05-28",
			"sales_count": 1048,
			"views_count": 20123,
		},
	}

	products := make(model.ProductList, 0, len(sampleProducts))

	for _, data := range sampleProducts {
		createdStr := data["created"].(string)
		created, err := model.ParseTime(createdStr)
		if err != nil {
			fmt.Printf("Error parsing date %s: %v\n", createdStr, err)
			continue
		}

		product := &model.Product{
			ID:         data["id"].(int),
			Name:       data["name"].(string),
			Price:      data["price"].(float64),
			Created:    created,
			SalesCount: data["sales_count"].(int),
			ViewsCount: data["views_count"].(int),
		}
		products = append(products, product)
	}

	repo := persistence.NewInMemoryProductRepository()
	if err := repo.Save(products); err != nil {
		fmt.Printf("Error saving products: %v\n", err)
		return nil
	}

	return repo
}

// loadConfiguration loads configuration from file or uses defaults
func loadConfiguration() *config.Config {
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

	return cfg
}

// displaySortedProducts sorts and displays products using the specified sorter
func displaySortedProducts(sorterUseCase *usecase.ProductSorterUseCase, products model.ProductList, sorterName string) {
	sortedProducts, err := sorterUseCase.SortProducts(products, sorterName)
	if err != nil {
		fmt.Printf("Error sorting products by %s: %v\n", sorterName, err)
		return
	}

	fmt.Printf("Products sorted by %s:\n", sorterName)
	for _, p := range sortedProducts {
		fmt.Println(p.String())
	}
}
