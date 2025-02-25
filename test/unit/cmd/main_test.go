package cmd_test

import (
	"testing"
	"time"
	
	"assessment/adapter/registry"
	"assessment/adapter/sorter"
	"assessment/domain/model"
	"assessment/infrastructure/config"
	"assessment/infrastructure/persistence"
	"assessment/usecase"
)


func TestMainFlow(t *testing.T) {
	
	
	
	products := createTestProducts()
	
	
	repo := persistence.NewInMemoryProductRepository()
	
	
	err := repo.Save(products)
	if err != nil {
		t.Fatalf("Failed to save products: %v", err)
	}
	
	
	reg := registry.NewSorterRegistry()
	
	
	sorterUseCase := usecase.NewProductSorterUseCase(reg)
	
	
	cfg := config.NewConfig()
	cfg.DisabledSorters = []string{"Name (descending)"}
	cfg.DefaultPageSize = 2
	
	
	sorterUseCase.SetConfig(cfg)
	
	
	sorter.InitializeDefaultSorters(reg, cfg)
	
	
	repoProducts, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get products: %v", err)
	}
	
	
	priceSortedProducts, err := sorterUseCase.SortProducts(repoProducts, "Price (ascending)")
	if err != nil {
		t.Fatalf("Failed to sort by price: %v", err)
	}
	
	
	if len(priceSortedProducts) != 3 {
		t.Errorf("Sorted products count mismatch: got %d, want %d", len(priceSortedProducts), 3)
	}
	
	
	paginationOptions := usecase.PaginationOptions{
		Page:     1,
		PageSize: 2,
	}
	
	
	paginatedResult, err := sorterUseCase.SortAndPaginateProducts(repoProducts, "Price (ascending)", paginationOptions)
	if err != nil {
		t.Fatalf("Failed to paginate products: %v", err)
	}
	
	
	if len(paginatedResult.Items) != 2 {
		t.Errorf("Paginated items count mismatch: got %d, want %d", len(paginatedResult.Items), 2)
	}
	
	
	printProducts(t, paginatedResult.Items)
}


func createTestProducts() model.ProductList {
	
	date1, _ := time.Parse("2006-01-02", "2019-01-04")
	date2, _ := time.Parse("2006-01-02", "2012-01-04")
	date3, _ := time.Parse("2006-01-02", "2014-05-28")
	
	
	return model.ProductList{
		{
			ID:         1,
			Name:       "Alabaster Table",
			Price:      12.99,
			Created:    date1,
			SalesCount: 32,
			ViewsCount: 730,
		},
		{
			ID:         2,
			Name:       "Zebra Table",
			Price:      44.49,
			Created:    date2,
			SalesCount: 301,
			ViewsCount: 3279,
		},
		{
			ID:         3,
			Name:       "Coffee Table",
			Price:      10.00,
			Created:    date3,
			SalesCount: 1048,
			ViewsCount: 20123,
		},
	}
}


func printProducts(t *testing.T, products model.ProductList) {
	for _, p := range products {
		t.Logf("%s", p.String())
	}
} 