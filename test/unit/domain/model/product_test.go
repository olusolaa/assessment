package model_test

import (
	"testing"
	"time"
	
	"assessment/domain/model"
)

func TestProductClone(t *testing.T) {
	
	created, _ := time.Parse("2006-01-02", "2022-01-01")
	original := model.ProductList{
		{
			ID:         1,
			Name:       "Test Product",
			Price:      10.99,
			Created:    created,
			SalesCount: 100,
			ViewsCount: 1000,
		},
		{
			ID:         2,
			Name:       "Another Product",
			Price:      20.99,
			Created:    created,
			SalesCount: 200,
			ViewsCount: 2000,
		},
	}
	
	
	cloned := original.Clone()
	
	
	if len(cloned) != len(original) {
		t.Errorf("Clone length mismatch: got %d, want %d", len(cloned), len(original))
	}
	
	
	for i, p := range original {
		if p.ID != cloned[i].ID {
			t.Errorf("Product ID mismatch at index %d: got %d, want %d", i, cloned[i].ID, p.ID)
		}
		if p.Name != cloned[i].Name {
			t.Errorf("Product Name mismatch at index %d: got %s, want %s", i, cloned[i].Name, p.Name)
		}
		if p.Price != cloned[i].Price {
			t.Errorf("Product Price mismatch at index %d: got %f, want %f", i, cloned[i].Price, p.Price)
		}
		if !p.Created.Equal(cloned[i].Created) {
			t.Errorf("Product Created mismatch at index %d: got %v, want %v", i, cloned[i].Created, p.Created)
		}
		if p.SalesCount != cloned[i].SalesCount {
			t.Errorf("Product SalesCount mismatch at index %d: got %d, want %d", i, cloned[i].SalesCount, p.SalesCount)
		}
		if p.ViewsCount != cloned[i].ViewsCount {
			t.Errorf("Product ViewsCount mismatch at index %d: got %d, want %d", i, cloned[i].ViewsCount, p.ViewsCount)
		}
	}
	
	
	cloned[0].Name = "Modified Name"
	cloned[0].Price = 99.99
	
	if original[0].Name == "Modified Name" {
		t.Error("Modifying clone affected original Name")
	}
	if original[0].Price == 99.99 {
		t.Error("Modifying clone affected original Price")
	}
}

func TestProductString(t *testing.T) {
	
	created, _ := time.Parse("2006-01-02", "2022-01-01")
	product := &model.Product{
		ID:         1,
		Name:       "Test Product",
		Price:      10.99,
		Created:    created,
		SalesCount: 100,
		ViewsCount: 1000,
	}
	
	
	expected := "ID: 1, Name: Test Product, Price: $10.99, Created: 2022-01-01, Sales/View: 0.100000"
	
	
	actual := product.String()
	
	
	if actual != expected {
		t.Errorf("String representation mismatch: got %s, want %s", actual, expected)
	}
	
	
	product.ViewsCount = 0
	expected = "ID: 1, Name: Test Product, Price: $10.99, Created: 2022-01-01, Sales/View: 0.000000"
	actual = product.String()
	
	if actual != expected {
		t.Errorf("String representation with zero views mismatch: got %s, want %s", actual, expected)
	}
}

func TestParseTime(t *testing.T) {
	
	validDate := "2022-01-01"
	expected, _ := time.Parse("2006-01-02", validDate)
	
	actual, err := model.ParseTime(validDate)
	if err != nil {
		t.Errorf("ParseTime returned error for valid date: %v", err)
	}
	
	if !actual.Equal(expected) {
		t.Errorf("ParseTime result mismatch: got %v, want %v", actual, expected)
	}
	
	
	invalidDate := "not-a-date"
	_, err = model.ParseTime(invalidDate)
	if err == nil {
		t.Error("ParseTime did not return error for invalid date")
	}
} 