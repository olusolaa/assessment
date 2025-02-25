package usecase_test

import (
	"testing"
	
	"assessment/adapter/registry"
	"assessment/domain/model"
	"assessment/infrastructure/config"
	"assessment/usecase"
)


func TestIsSorterEnabledEdgeCases(t *testing.T) {
	
	reg := registry.NewSorterRegistry()
	
	
	sorterUseCase := usecase.NewProductSorterUseCase(reg)
	
	
	cfg := config.NewConfig()
	cfg.DisabledSorters = nil
	
	
	sorterUseCase.SetConfig(cfg)
	
	
	products := model.ProductList{
		{
			ID:   1,
			Name: "Test Product",
		},
	}
	
	
	_, err := sorterUseCase.SortProducts(products, "NonExistentSorter")
	
	
	if err == nil {
		t.Error("SortProducts did not return error for non-existent sorter")
	} else if err.Error() != "sorter not found: NonExistentSorter" {
		t.Errorf("Unexpected error: %v", err)
	}
	
	
	mockSorter := NewMockSorter("MockSorter")
	
	
	reg.RegisterSorter(mockSorter)
	
	
	_, err = sorterUseCase.SortProducts(products, "MockSorter")
	if err != nil {
		t.Errorf("SortProducts failed with nil DisabledSorters: %v", err)
	}
	
	
	cfg = config.NewConfig()
	cfg.DisabledSorters = []string{}
	
	
	sorterUseCase.SetConfig(cfg)
	
	
	_, err = sorterUseCase.SortProducts(products, "MockSorter")
	if err != nil {
		t.Errorf("SortProducts failed with empty DisabledSorters: %v", err)
	}
	
	
	cfg = config.NewConfig()
	cfg.DisabledSorters = []string{"MockSorter"}
	
	
	sorterUseCase.SetConfig(cfg)
	
	
	_, err = sorterUseCase.SortProducts(products, "MockSorter")
	if err == nil {
		t.Error("SortProducts did not return error for disabled sorter")
	} else if err.Error() != "sorter is disabled: MockSorter" {
		t.Errorf("Unexpected error: %v", err)
	}
	
	
	
	sorterUseCase = usecase.NewProductSorterUseCase(reg)
	
	
	_, err = sorterUseCase.SortProducts(products, "MockSorter")
	if err != nil {
		t.Errorf("SortProducts failed with nil config: %v", err)
	}
} 