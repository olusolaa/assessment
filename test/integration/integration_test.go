package integration

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


func TestIntegration(t *testing.T) {
	
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
	
	
	availableSorters := sorterUseCase.GetAvailableSorters()
	if len(availableSorters) != 7 {
		t.Errorf("Available sorters count mismatch: got %d, want %d", len(availableSorters), 7)
	}
	
	
	for _, name := range availableSorters {
		if name == "Name (descending)" {
			t.Error("Disabled sorter is available")
		}
	}
	
	
	priceSortedProducts, err := sorterUseCase.SortProducts(repoProducts, "Price (ascending)")
	if err != nil {
		t.Fatalf("Failed to sort by price: %v", err)
	}
	
	
	if len(priceSortedProducts) != 3 {
		t.Errorf("Sorted products count mismatch: got %d, want %d", len(priceSortedProducts), 3)
	}
	
	if priceSortedProducts[0].Price > priceSortedProducts[1].Price || priceSortedProducts[1].Price > priceSortedProducts[2].Price {
		t.Error("Products not sorted correctly by price (ascending)")
	}
	
	
	nameSortedProducts, err := sorterUseCase.SortProducts(repoProducts, "Name (ascending)")
	if err != nil {
		t.Fatalf("Failed to sort by name: %v", err)
	}
	
	
	if nameSortedProducts[0].Name > nameSortedProducts[1].Name || nameSortedProducts[1].Name > nameSortedProducts[2].Name {
		t.Error("Products not sorted correctly by name (ascending)")
	}
	
	
	dateSortedProducts, err := sorterUseCase.SortProducts(repoProducts, "Creation Date (ascending)")
	if err != nil {
		t.Fatalf("Failed to sort by date: %v", err)
	}
	
	
	if dateSortedProducts[0].Created.After(dateSortedProducts[1].Created) || dateSortedProducts[1].Created.After(dateSortedProducts[2].Created) {
		t.Error("Products not sorted correctly by date (ascending)")
	}
	
	
	spvSortedProducts, err := sorterUseCase.SortProducts(repoProducts, "Sales per View (descending)")
	if err != nil {
		t.Fatalf("Failed to sort by sales per view: %v", err)
	}
	
	
	spv0 := float64(spvSortedProducts[0].SalesCount) / float64(spvSortedProducts[0].ViewsCount)
	spv1 := float64(spvSortedProducts[1].SalesCount) / float64(spvSortedProducts[1].ViewsCount)
	spv2 := float64(spvSortedProducts[2].SalesCount) / float64(spvSortedProducts[2].ViewsCount)
	
	
	if spv0 < spv1 || spv1 < spv2 {
		t.Error("Products not sorted correctly by sales per view (descending)")
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
	
	if paginatedResult.Page != 1 {
		t.Errorf("Page mismatch: got %d, want %d", paginatedResult.Page, 1)
	}
	
	if paginatedResult.PageSize != 2 {
		t.Errorf("PageSize mismatch: got %d, want %d", paginatedResult.PageSize, 2)
	}
	
	if paginatedResult.TotalItems != 3 {
		t.Errorf("TotalItems mismatch: got %d, want %d", paginatedResult.TotalItems, 3)
	}
	
	if paginatedResult.TotalPages != 2 {
		t.Errorf("TotalPages mismatch: got %d, want %d", paginatedResult.TotalPages, 2)
	}
	
	if !paginatedResult.HasNext {
		t.Error("HasNext should be true")
	}
	
	if paginatedResult.HasPrev {
		t.Error("HasPrev should be false")
	}
	
	
	paginationOptions.Page = 2
	
	
	paginatedResult, err = sorterUseCase.SortAndPaginateProducts(repoProducts, "Price (ascending)", paginationOptions)
	if err != nil {
		t.Fatalf("Failed to paginate products: %v", err)
	}
	
	
	if len(paginatedResult.Items) != 1 {
		t.Errorf("Paginated items count mismatch: got %d, want %d", len(paginatedResult.Items), 1)
	}
	
	if paginatedResult.Page != 2 {
		t.Errorf("Page mismatch: got %d, want %d", paginatedResult.Page, 2)
	}
	
	if !paginatedResult.HasPrev {
		t.Error("HasPrev should be true")
	}
	
	if paginatedResult.HasNext {
		t.Error("HasNext should be false")
	}
	
	
	_, err = sorterUseCase.SortProducts(repoProducts, "Name (descending)")
	if err == nil {
		t.Error("Sorting with disabled sorter should fail")
	}
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