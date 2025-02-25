package sorter_test

import (
	"testing"
	"time"
	
	"assessment/adapter/sorter"
	"assessment/domain/model"
)


func createTestProducts() model.ProductList {
	
	date1, _ := time.Parse("2006-01-02", "2020-01-01")
	date2, _ := time.Parse("2006-01-02", "2019-01-01")
	date3, _ := time.Parse("2006-01-02", "2021-01-01")
	
	
	return model.ProductList{
		{
			ID:         1,
			Name:       "B Product",
			Price:      20.0,
			Created:    date1,
			SalesCount: 100,
			ViewsCount: 1000,
		},
		{
			ID:         2,
			Name:       "C Product",
			Price:      10.0,
			Created:    date2,
			SalesCount: 200,
			ViewsCount: 1000,
		},
		{
			ID:         3,
			Name:       "A Product",
			Price:      30.0,
			Created:    date3,
			SalesCount: 300,
			ViewsCount: 1000,
		},
	}
}


func verifyOriginalUnchanged(t *testing.T, original, sorted model.ProductList) {
	
	if len(original) != 3 {
		t.Errorf("Original products length changed: got %d, want %d", len(original), 3)
	}
	
	
	if original[0].ID != 1 || original[1].ID != 2 || original[2].ID != 3 {
		t.Error("Original products were modified")
	}
	
	
	for i := range original {
		if &original[i] == &sorted[i] {
			t.Error("Sorted products are the same objects as the original")
		}
	}
}

func TestPriceSorter(t *testing.T) {
	
	products := createTestProducts()
	
	
	ascendingSorter := sorter.NewPriceSorter(true)
	descendingSorter := sorter.NewPriceSorter(false)
	
	
	sortedAscending := ascendingSorter.Sort(products)
	
	
	if sortedAscending[0].Price != 10.0 || sortedAscending[1].Price != 20.0 || sortedAscending[2].Price != 30.0 {
		t.Error("Products not sorted correctly by price (ascending)")
	}
	
	
	verifyOriginalUnchanged(t, products, sortedAscending)
	
	
	sortedDescending := descendingSorter.Sort(products)
	
	
	if sortedDescending[0].Price != 30.0 || sortedDescending[1].Price != 20.0 || sortedDescending[2].Price != 10.0 {
		t.Error("Products not sorted correctly by price (descending)")
	}
	
	
	verifyOriginalUnchanged(t, products, sortedDescending)
	
	
	if ascendingSorter.Name() != "Price (ascending)" {
		t.Errorf("Sorter name mismatch: got %s, want %s", ascendingSorter.Name(), "Price (ascending)")
	}
	
	if descendingSorter.Name() != "Price (descending)" {
		t.Errorf("Sorter name mismatch: got %s, want %s", descendingSorter.Name(), "Price (descending)")
	}
}

func TestDateSorter(t *testing.T) {
	
	products := createTestProducts()
	
	
	ascendingSorter := sorter.NewDateSorter(true)
	descendingSorter := sorter.NewDateSorter(false)
	
	
	sortedAscending := ascendingSorter.Sort(products)
	
	
	if sortedAscending[0].ID != 2 || sortedAscending[1].ID != 1 || sortedAscending[2].ID != 3 {
		t.Error("Products not sorted correctly by date (ascending)")
	}
	
	
	verifyOriginalUnchanged(t, products, sortedAscending)
	
	
	sortedDescending := descendingSorter.Sort(products)
	
	
	if sortedDescending[0].ID != 3 || sortedDescending[1].ID != 1 || sortedDescending[2].ID != 2 {
		t.Error("Products not sorted correctly by date (descending)")
	}
	
	
	verifyOriginalUnchanged(t, products, sortedDescending)
	
	
	if ascendingSorter.Name() != "Creation Date (ascending)" {
		t.Errorf("Sorter name mismatch: got %s, want %s", ascendingSorter.Name(), "Creation Date (ascending)")
	}
	
	if descendingSorter.Name() != "Creation Date (descending)" {
		t.Errorf("Sorter name mismatch: got %s, want %s", descendingSorter.Name(), "Creation Date (descending)")
	}
}

func TestNameSorter(t *testing.T) {
	
	products := createTestProducts()
	
	
	ascendingSorter := sorter.NewNameSorter(true)
	descendingSorter := sorter.NewNameSorter(false)
	
	
	sortedAscending := ascendingSorter.Sort(products)
	
	
	if sortedAscending[0].Name != "A Product" || sortedAscending[1].Name != "B Product" || sortedAscending[2].Name != "C Product" {
		t.Error("Products not sorted correctly by name (ascending)")
	}
	
	
	verifyOriginalUnchanged(t, products, sortedAscending)
	
	
	sortedDescending := descendingSorter.Sort(products)
	
	
	if sortedDescending[0].Name != "C Product" || sortedDescending[1].Name != "B Product" || sortedDescending[2].Name != "A Product" {
		t.Error("Products not sorted correctly by name (descending)")
	}
	
	
	verifyOriginalUnchanged(t, products, sortedDescending)
	
	
	if ascendingSorter.Name() != "Name (ascending)" {
		t.Errorf("Sorter name mismatch: got %s, want %s", ascendingSorter.Name(), "Name (ascending)")
	}
	
	if descendingSorter.Name() != "Name (descending)" {
		t.Errorf("Sorter name mismatch: got %s, want %s", descendingSorter.Name(), "Name (descending)")
	}
}

func TestSalesPerViewSorter(t *testing.T) {
	
	date, _ := time.Parse("2006-01-02", "2020-01-01")
	products := model.ProductList{
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
			SalesCount: 300,
			ViewsCount: 1000, 
		},
		{
			ID:         3,
			Name:       "Product 3",
			Price:      30.0,
			Created:    date,
			SalesCount: 200,
			ViewsCount: 1000, 
		},
	}
	
	
	ascendingSorter := sorter.NewSalesPerViewSorter(true)
	descendingSorter := sorter.NewSalesPerViewSorter(false)
	
	
	sortedAscending := ascendingSorter.Sort(products)
	
	
	if sortedAscending[0].ID != 1 || sortedAscending[1].ID != 3 || sortedAscending[2].ID != 2 {
		t.Error("Products not sorted correctly by sales per view (ascending)")
	}
	
	
	verifyOriginalUnchanged(t, products, sortedAscending)
	
	
	sortedDescending := descendingSorter.Sort(products)
	
	
	if sortedDescending[0].ID != 2 || sortedDescending[1].ID != 3 || sortedDescending[2].ID != 1 {
		t.Error("Products not sorted correctly by sales per view (descending)")
	}
	
	
	verifyOriginalUnchanged(t, products, sortedDescending)
	
	
	productsWithZeroViews := model.ProductList{
		{
			ID:         1,
			Name:       "Product 1",
			Price:      10.0,
			Created:    date,
			SalesCount: 100,
			ViewsCount: 0, 
		},
		{
			ID:         2,
			Name:       "Product 2",
			Price:      20.0,
			Created:    date,
			SalesCount: 300,
			ViewsCount: 1000, 
		},
	}
	
	
	sortedWithZeroViews := ascendingSorter.Sort(productsWithZeroViews)
	
	
	if sortedWithZeroViews[0].ID != 1 || sortedWithZeroViews[1].ID != 2 {
		t.Error("Products with zero views not sorted correctly")
	}
	
	
	if ascendingSorter.Name() != "Sales per View (ascending)" {
		t.Errorf("Sorter name mismatch: got %s, want %s", ascendingSorter.Name(), "Sales per View (ascending)")
	}
	
	if descendingSorter.Name() != "Sales per View (descending)" {
		t.Errorf("Sorter name mismatch: got %s, want %s", descendingSorter.Name(), "Sales per View (descending)")
	}
} 