package sorters

import (
	"assessment/adapter/sorter"
	"assessment/domain/model"
	"testing"
	"time"
)

func createTestProducts() model.ProductList {
	return model.ProductList{
		{
			ID:         1,
			Name:       "Product A",
			Price:      10.99,
			Created:    time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
			SalesCount: 100,
			ViewsCount: 1000,
		},
		{
			ID:         2,
			Name:       "Product B",
			Price:      5.99,
			Created:    time.Date(2021, 2, 20, 0, 0, 0, 0, time.UTC),
			SalesCount: 200,
			ViewsCount: 1000,
		},
		{
			ID:         3,
			Name:       "Product C",
			Price:      15.99,
			Created:    time.Date(2019, 3, 25, 0, 0, 0, 0, time.UTC),
			SalesCount: 50,
			ViewsCount: 1000,
		},
	}
}

func TestPriceSorter(t *testing.T) {
	products := createTestProducts()

	originalProducts := make([]*model.Product, len(products))
	for i, p := range products {
		originalProducts[i] = p
	}

	ascSorter := sorter.NewPriceSorter(true)
	ascSorted := ascSorter.Sort(products)

	for i := 0; i < len(ascSorted)-1; i++ {
		if ascSorted[i].Price > ascSorted[i+1].Price {
			t.Errorf("Products not sorted by price in ascending order: %f > %f",
				ascSorted[i].Price, ascSorted[i+1].Price)
		}
	}

	descSorter := sorter.NewPriceSorter(false)
	descSorted := descSorter.Sort(products)

	for i := 0; i < len(descSorted)-1; i++ {
		if descSorted[i].Price < descSorted[i+1].Price {
			t.Errorf("Products not sorted by price in descending order: %f < %f",
				descSorted[i].Price, descSorted[i+1].Price)
		}
	}

	for i, p := range products {
		if p != originalProducts[i] {
			t.Errorf("Original products slice was modified during sorting at index %d", i)
		}
	}
}

func TestSalesPerViewSorter(t *testing.T) {
	products := createTestProducts()

	ascSorter := sorter.NewSalesPerViewSorter(true)
	ascSorted := ascSorter.Sort(products)

	for i := 0; i < len(ascSorted)-1; i++ {
		ratio1 := float64(ascSorted[i].SalesCount) / float64(ascSorted[i].ViewsCount)
		ratio2 := float64(ascSorted[i+1].SalesCount) / float64(ascSorted[i+1].ViewsCount)
		if ratio1 > ratio2 {
			t.Errorf("Products not sorted by sales per view in ascending order: %f > %f",
				ratio1, ratio2)
		}
	}

	descSorter := sorter.NewSalesPerViewSorter(false)
	descSorted := descSorter.Sort(products)

	for i := 0; i < len(descSorted)-1; i++ {
		ratio1 := float64(descSorted[i].SalesCount) / float64(descSorted[i].ViewsCount)
		ratio2 := float64(descSorted[i+1].SalesCount) / float64(descSorted[i+1].ViewsCount)
		if ratio1 < ratio2 {
			t.Errorf("Products not sorted by sales per view in descending order: %f < %f",
				ratio1, ratio2)
		}
	}
}

func TestCreationDateSorter(t *testing.T) {
	products := createTestProducts()

	ascSorter := sorter.NewDateSorter(true)
	ascSorted := ascSorter.Sort(products)

	for i := 0; i < len(ascSorted)-1; i++ {
		if ascSorted[i].Created.After(ascSorted[i+1].Created) {
			t.Errorf("Products not sorted by creation date in ascending order: %v after %v",
				ascSorted[i].Created, ascSorted[i+1].Created)
		}
	}

	descSorter := sorter.NewDateSorter(false)
	descSorted := descSorter.Sort(products)

	for i := 0; i < len(descSorted)-1; i++ {
		if descSorted[i].Created.Before(descSorted[i+1].Created) {
			t.Errorf("Products not sorted by creation date in descending order: %v before %v",
				descSorted[i].Created, descSorted[i+1].Created)
		}
	}
}

func TestNameSorter(t *testing.T) {
	products := createTestProducts()

	ascSorter := sorter.NewNameSorter(true)
	ascSorted := ascSorter.Sort(products)

	for i := 0; i < len(ascSorted)-1; i++ {
		if ascSorted[i].Name > ascSorted[i+1].Name {
			t.Errorf("Products not sorted by name in ascending order: %s > %s",
				ascSorted[i].Name, ascSorted[i+1].Name)
		}
	}

	descSorter := sorter.NewNameSorter(false)
	descSorted := descSorter.Sort(products)

	for i := 0; i < len(descSorted)-1; i++ {
		if descSorted[i].Name < descSorted[i+1].Name {
			t.Errorf("Products not sorted by name in descending order: %s < %s",
				descSorted[i].Name, descSorted[i+1].Name)
		}
	}
}
