package persistence_test

import (
	"testing"
	"time"

	"assessment/domain/model"
	"assessment/infrastructure/persistence"
)

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

func TestInMemoryProductRepositoryGetAll(t *testing.T) {

	repo := persistence.NewInMemoryProductRepository()

	products := createTestProducts()
	err := repo.Save(products)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	retrievedProducts, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(retrievedProducts) != 3 {
		t.Errorf("Product count mismatch: got %d, want %d", len(retrievedProducts), 3)
	}

	for i := range products {
		if products[i] == retrievedProducts[i] {
			t.Error("Retrieved products are the same objects as the original")
		}
	}

	for i, p := range products {
		if p.ID != retrievedProducts[i].ID {
			t.Errorf("Product ID mismatch at index %d: got %d, want %d", i, retrievedProducts[i].ID, p.ID)
		}
		if p.Name != retrievedProducts[i].Name {
			t.Errorf("Product Name mismatch at index %d: got %s, want %s", i, retrievedProducts[i].Name, p.Name)
		}
		if p.Price != retrievedProducts[i].Price {
			t.Errorf("Product Price mismatch at index %d: got %f, want %f", i, retrievedProducts[i].Price, p.Price)
		}
		if !p.Created.Equal(retrievedProducts[i].Created) {
			t.Errorf("Product Created mismatch at index %d: got %v, want %v", i, retrievedProducts[i].Created, p.Created)
		}
		if p.SalesCount != retrievedProducts[i].SalesCount {
			t.Errorf("Product SalesCount mismatch at index %d: got %d, want %d", i, retrievedProducts[i].SalesCount, p.SalesCount)
		}
		if p.ViewsCount != retrievedProducts[i].ViewsCount {
			t.Errorf("Product ViewsCount mismatch at index %d: got %d, want %d", i, retrievedProducts[i].ViewsCount, p.ViewsCount)
		}
	}

	retrievedProducts[0].Name = "Modified Name"

	retrievedProducts2, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if retrievedProducts2[0].Name == "Modified Name" {
		t.Error("Modifying retrieved products affected the repository")
	}
}

func TestInMemoryProductRepositoryGetByIDs(t *testing.T) {

	repo := persistence.NewInMemoryProductRepository()

	products := createTestProducts()
	err := repo.Save(products)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	retrievedProducts, err := repo.GetByIDs([]int{1, 3})
	if err != nil {
		t.Fatalf("GetByIDs failed: %v", err)
	}

	if len(retrievedProducts) != 2 {
		t.Errorf("Product count mismatch: got %d, want %d", len(retrievedProducts), 2)
	}

	productMap := make(map[int]*model.Product)
	for _, p := range retrievedProducts {
		productMap[p.ID] = p
	}

	if _, exists := productMap[1]; !exists {
		t.Error("Product with ID 1 not found")
	}
	if _, exists := productMap[3]; !exists {
		t.Error("Product with ID 3 not found")
	}

	for _, p := range retrievedProducts {
		for _, original := range products {
			if p.ID == original.ID && p == original {
				t.Error("Retrieved products are the same objects as the original")
			}
		}
	}

	retrievedProducts, err = repo.GetByIDs([]int{4, 5})
	if err != nil {
		t.Fatalf("GetByIDs failed: %v", err)
	}

	if len(retrievedProducts) != 0 {
		t.Errorf("Product count mismatch for non-existent IDs: got %d, want %d", len(retrievedProducts), 0)
	}
}

func TestInMemoryProductRepositorySave(t *testing.T) {

	repo := persistence.NewInMemoryProductRepository()

	products := createTestProducts()
	err := repo.Save(products)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	retrievedProducts, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(retrievedProducts) != 3 {
		t.Errorf("Product count mismatch: got %d, want %d", len(retrievedProducts), 3)
	}

	products[0].Name = "Modified Name"

	retrievedProducts2, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if retrievedProducts2[0].Name == "Modified Name" {
		t.Error("Modifying original products affected the repository")
	}

	newProducts := model.ProductList{
		{
			ID:         4,
			Name:       "Product 4",
			Price:      40.0,
			Created:    time.Now(),
			SalesCount: 400,
			ViewsCount: 4000,
		},
	}

	err = repo.Save(newProducts)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	retrievedProducts3, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(retrievedProducts3) != 1 {
		t.Errorf("Product count mismatch after second save: got %d, want %d", len(retrievedProducts3), 1)
	}

	if retrievedProducts3[0].ID != 4 {
		t.Errorf("Product ID mismatch: got %d, want %d", retrievedProducts3[0].ID, 4)
	}
}
