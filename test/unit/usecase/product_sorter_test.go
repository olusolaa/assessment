package usecase_test

import (
	"testing"
	"time"
	
	"assessment/adapter/registry"
	"assessment/domain/model"
	"assessment/infrastructure/config"
	"assessment/usecase"
)


type MockSorter struct {
	name string
}

func NewMockSorter(name string) *MockSorter {
	return &MockSorter{name: name}
}

func (s *MockSorter) Sort(products model.ProductList) model.ProductList {
	
	return products.Clone()
}

func (s *MockSorter) Name() string {
	return s.name
}


func createTestProducts() model.ProductList {
	
	date, _ := time.Parse("2006-01-02", "2020-01-01")
	
	
	return model.ProductList{
		{
			ID:         1,
			Name:       "Product 1",
			Price:      10.0,
			Created:    date,
			SalesCount: 100,
			ViewsCount: 1000,
		},
		{
			ID:         2,
			Name:       "Product 2",
			Price:      20.0,
			Created:    date,
			SalesCount: 200,
			ViewsCount: 2000,
		},
		{
			ID:         3,
			Name:       "Product 3",
			Price:      30.0,
			Created:    date,
			SalesCount: 300,
			ViewsCount: 3000,
		},
	}
}

func TestProductSorterUseCaseSortProducts(t *testing.T) {
	
	reg := registry.NewSorterRegistry()
	
	
	mockSorter := NewMockSorter("MockSorter")
	
	
	reg.RegisterSorter(mockSorter)
	
	
	sorterUseCase := usecase.NewProductSorterUseCase(reg)
	
	
	products := createTestProducts()
	
	
	sortedProducts, err := sorterUseCase.SortProducts(products, "MockSorter")
	if err != nil {
		t.Fatalf("SortProducts failed: %v", err)
	}
	
	
	if len(sortedProducts) != 3 {
		t.Errorf("Product count mismatch: got %d, want %d", len(sortedProducts), 3)
	}
	
	
	_, err = sorterUseCase.SortProducts(products, "NonExistentSorter")
	if err == nil {
		t.Error("SortProducts did not return error for non-existent sorter")
	}
}

func TestProductSorterUseCaseDisabledSorter(t *testing.T) {
	
	reg := registry.NewSorterRegistry()
	
	
	mockSorter1 := NewMockSorter("MockSorter1")
	mockSorter2 := NewMockSorter("MockSorter2")
	
	
	reg.RegisterSorter(mockSorter1)
	reg.RegisterSorter(mockSorter2)
	
	
	sorterUseCase := usecase.NewProductSorterUseCase(reg)
	
	
	cfg := config.NewConfig()
	cfg.DisabledSorters = []string{"MockSorter1"}
	
	
	sorterUseCase.SetConfig(cfg)
	
	
	products := createTestProducts()
	
	
	_, err := sorterUseCase.SortProducts(products, "MockSorter1")
	if err == nil {
		t.Error("SortProducts did not return error for disabled sorter")
	}
	
	
	_, err = sorterUseCase.SortProducts(products, "MockSorter2")
	if err != nil {
		t.Errorf("SortProducts failed for enabled sorter: %v", err)
	}
}

func TestProductSorterUseCaseGetAvailableSorters(t *testing.T) {
	
	reg := registry.NewSorterRegistry()
	
	
	mockSorter1 := NewMockSorter("MockSorter1")
	mockSorter2 := NewMockSorter("MockSorter2")
	mockSorter3 := NewMockSorter("MockSorter3")
	
	
	reg.RegisterSorter(mockSorter1)
	reg.RegisterSorter(mockSorter2)
	reg.RegisterSorter(mockSorter3)
	
	
	sorterUseCase := usecase.NewProductSorterUseCase(reg)
	
	
	sorters := sorterUseCase.GetAvailableSorters()
	
	
	if len(sorters) != 3 {
		t.Errorf("Sorter count mismatch: got %d, want %d", len(sorters), 3)
	}
	
	
	sorterMap := make(map[string]bool)
	for _, s := range sorters {
		sorterMap[s] = true
	}
	
	
	if !sorterMap["MockSorter1"] {
		t.Error("MockSorter1 not found in GetAvailableSorters result")
	}
	if !sorterMap["MockSorter2"] {
		t.Error("MockSorter2 not found in GetAvailableSorters result")
	}
	if !sorterMap["MockSorter3"] {
		t.Error("MockSorter3 not found in GetAvailableSorters result")
	}
	
	
	cfg := config.NewConfig()
	cfg.DisabledSorters = []string{"MockSorter1", "MockSorter3"}
	
	
	sorterUseCase.SetConfig(cfg)
	
	
	sorters = sorterUseCase.GetAvailableSorters()
	
	
	if len(sorters) != 1 {
		t.Errorf("Sorter count mismatch after disabling: got %d, want %d", len(sorters), 1)
	}
	
	
	if sorters[0] != "MockSorter2" {
		t.Errorf("Enabled sorter mismatch: got %s, want %s", sorters[0], "MockSorter2")
	}
}

func TestProductSorterUseCaseGetSetConfig(t *testing.T) {
	
	reg := registry.NewSorterRegistry()
	
	
	sorterUseCase := usecase.NewProductSorterUseCase(reg)
	
	
	cfg := sorterUseCase.GetConfig()
	
	
	if len(cfg.DisabledSorters) != 0 {
		t.Errorf("Default DisabledSorters not empty: got %v", cfg.DisabledSorters)
	}
	
	if cfg.DefaultPageSize != 10 {
		t.Errorf("Default DefaultPageSize mismatch: got %d, want %d", cfg.DefaultPageSize, 10)
	}
	
	
	newCfg := config.NewConfig()
	newCfg.DisabledSorters = []string{"Sorter1", "Sorter2"}
	newCfg.DefaultPageSize = 20
	
	
	sorterUseCase.SetConfig(newCfg)
	
	
	cfg = sorterUseCase.GetConfig()
	
	
	if len(cfg.DisabledSorters) != 2 {
		t.Errorf("DisabledSorters length mismatch: got %d, want %d", len(cfg.DisabledSorters), 2)
	}
	
	if cfg.DisabledSorters[0] != "Sorter1" || cfg.DisabledSorters[1] != "Sorter2" {
		t.Errorf("DisabledSorters mismatch: got %v, want %v", cfg.DisabledSorters, []string{"Sorter1", "Sorter2"})
	}
	
	if cfg.DefaultPageSize != 20 {
		t.Errorf("DefaultPageSize mismatch: got %d, want %d", cfg.DefaultPageSize, 20)
	}
}

func TestProductSorterUseCaseGetRegistry(t *testing.T) {
	
	reg := registry.NewSorterRegistry()
	
	
	sorterUseCase := usecase.NewProductSorterUseCase(reg)
	
	
	retrievedReg := sorterUseCase.GetRegistry()
	
	
	if retrievedReg != reg {
		t.Error("Retrieved registry is not the same as the original")
	}
} 