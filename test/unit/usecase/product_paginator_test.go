package usecase_test

import (
	"fmt"
	"testing"
	"time"

	"assessment/adapter/registry"
	"assessment/domain/model"
	"assessment/infrastructure/config"
	"assessment/usecase"
)

func createPaginationTestProducts() model.ProductList {

	date, _ := time.Parse("2006-01-02", "2020-01-01")

	products := make(model.ProductList, 10)
	for i := 0; i < 10; i++ {
		products[i] = &model.Product{
			ID:         i + 1,
			Name:       fmt.Sprintf("Product %d", i+1),
			Price:      float64(10 * (i + 1)),
			Created:    date,
			SalesCount: 100 * (i + 1),
			ViewsCount: 1000 * (i + 1),
		}
	}

	return products
}

func TestPagination(t *testing.T) {

	reg := registry.NewSorterRegistry()

	mockSorter := NewMockSorter("MockSorter")

	reg.RegisterSorter(mockSorter)

	sorterUseCase := usecase.NewProductSorterUseCase(reg)

	products := createPaginationTestProducts()

	testCases := []struct {
		name            string
		page            int
		pageSize        int
		expectedItems   int
		expectedPage    int
		expectedTotal   int
		expectedPages   int
		expectedHasNext bool
		expectedHasPrev bool
	}{
		{
			name:            "First page with items",
			page:            1,
			pageSize:        3,
			expectedItems:   3,
			expectedPage:    1,
			expectedTotal:   10,
			expectedPages:   4,
			expectedHasNext: true,
			expectedHasPrev: false,
		},
		{
			name:            "Middle page",
			page:            2,
			pageSize:        3,
			expectedItems:   3,
			expectedPage:    2,
			expectedTotal:   10,
			expectedPages:   4,
			expectedHasNext: true,
			expectedHasPrev: true,
		},
		{
			name:            "Last page with items",
			page:            4,
			pageSize:        3,
			expectedItems:   1,
			expectedPage:    4,
			expectedTotal:   10,
			expectedPages:   4,
			expectedHasNext: false,
			expectedHasPrev: true,
		},
		{
			name:            "Page beyond total pages",
			page:            5,
			pageSize:        3,
			expectedItems:   1,
			expectedPage:    4,
			expectedTotal:   10,
			expectedPages:   4,
			expectedHasNext: false,
			expectedHasPrev: true,
		},
		{
			name:            "Invalid page (negative)",
			page:            -1,
			pageSize:        3,
			expectedItems:   3,
			expectedPage:    1,
			expectedTotal:   10,
			expectedPages:   4,
			expectedHasNext: true,
			expectedHasPrev: false,
		},
		{
			name:            "Invalid page size (negative)",
			page:            1,
			pageSize:        -1,
			expectedItems:   10,
			expectedPage:    1,
			expectedTotal:   10,
			expectedPages:   1,
			expectedHasNext: false,
			expectedHasPrev: false,
		},
		{
			name:            "Page size larger than total items",
			page:            1,
			pageSize:        20,
			expectedItems:   10,
			expectedPage:    1,
			expectedTotal:   10,
			expectedPages:   1,
			expectedHasNext: false,
			expectedHasPrev: false,
		},
		{
			name:            "Page too high with remainder",
			page:            3,
			pageSize:        5,
			expectedItems:   5,
			expectedPage:    2,
			expectedTotal:   10,
			expectedPages:   2,
			expectedHasNext: false,
			expectedHasPrev: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			options := usecase.PaginationOptions{
				Page:     tc.page,
				PageSize: tc.pageSize,
			}

			result, err := sorterUseCase.SortAndPaginateProducts(products, "MockSorter", options)
			if err != nil {
				t.Fatalf("SortAndPaginateProducts failed: %v", err)
			}

			if len(result.Items) != tc.expectedItems {
				t.Errorf("Item count mismatch: got %d, want %d", len(result.Items), tc.expectedItems)
			}

			if result.Page != tc.expectedPage {
				t.Errorf("Page mismatch: got %d, want %d", result.Page, tc.expectedPage)
			}

			expectedPageSize := tc.pageSize
			if tc.pageSize < 1 {
				expectedPageSize = 10
			}

			if result.PageSize != expectedPageSize {
				t.Errorf("PageSize mismatch: got %d, want %d", result.PageSize, expectedPageSize)
			}

			if result.TotalItems != tc.expectedTotal {
				t.Errorf("TotalItems mismatch: got %d, want %d", result.TotalItems, tc.expectedTotal)
			}

			if result.TotalPages != tc.expectedPages {
				t.Errorf("TotalPages mismatch: got %d, want %d", result.TotalPages, tc.expectedPages)
			}

			if result.HasNext != tc.expectedHasNext {
				t.Errorf("HasNext mismatch: got %v, want %v", result.HasNext, tc.expectedHasNext)
			}

			if result.HasPrev != tc.expectedHasPrev {
				t.Errorf("HasPrev mismatch: got %v, want %v", result.HasPrev, tc.expectedHasPrev)
			}
		})
	}
}

func TestPaginationWithDisabledSorter(t *testing.T) {

	reg := registry.NewSorterRegistry()

	mockSorter := NewMockSorter("MockSorter")

	reg.RegisterSorter(mockSorter)

	sorterUseCase := usecase.NewProductSorterUseCase(reg)

	cfg := config.NewConfig()
	cfg.DisabledSorters = []string{"MockSorter"}

	sorterUseCase.SetConfig(cfg)

	products := createPaginationTestProducts()

	options := usecase.PaginationOptions{
		Page:     1,
		PageSize: 3,
	}

	_, err := sorterUseCase.SortAndPaginateProducts(products, "MockSorter", options)

	if err == nil {
		t.Error("SortAndPaginateProducts did not return error for disabled sorter")
	}
}
