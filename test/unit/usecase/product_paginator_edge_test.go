package usecase_test

import (
	"testing"

	"assessment/adapter/registry"
	"assessment/domain/model"
	"assessment/usecase"
)

func TestPaginateEdgeCases(t *testing.T) {

	reg := registry.NewSorterRegistry()

	mockSorter := NewMockSorter("MockSorter")

	reg.RegisterSorter(mockSorter)

	sorterUseCase := usecase.NewProductSorterUseCase(reg)

	emptyProducts := model.ProductList{}

	options := usecase.PaginationOptions{
		Page:     1,
		PageSize: 10,
	}

	result, err := sorterUseCase.SortAndPaginateProducts(emptyProducts, "MockSorter", options)
	if err != nil {
		t.Fatalf("SortAndPaginateProducts failed with empty list: %v", err)
	}

	if len(result.Items) != 0 {
		t.Errorf("Item count mismatch for empty list: got %d, want %d", len(result.Items), 0)
	}

	if result.Page != 1 {
		t.Errorf("Page mismatch for empty list: got %d, want %d", result.Page, 1)
	}

	if result.TotalItems != 0 {
		t.Errorf("TotalItems mismatch for empty list: got %d, want %d", result.TotalItems, 0)
	}

	if result.TotalPages != 0 {
		t.Errorf("TotalPages mismatch for empty list: got %d, want %d", result.TotalPages, 0)
	}

	if result.HasNext {
		t.Error("HasNext should be false for empty list")
	}

	if result.HasPrev {
		t.Error("HasPrev should be false for empty list")
	}

	options = usecase.PaginationOptions{
		Page:     1,
		PageSize: 0,
	}

	singleProduct := model.ProductList{
		{
			ID:   1,
			Name: "Test Product",
		},
	}

	result, err = sorterUseCase.SortAndPaginateProducts(singleProduct, "MockSorter", options)
	if err != nil {
		t.Fatalf("SortAndPaginateProducts failed with zero page size: %v", err)
	}

	if result.PageSize != 10 {
		t.Errorf("PageSize mismatch for zero page size: got %d, want %d", result.PageSize, 10)
	}
}
