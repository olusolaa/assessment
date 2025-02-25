package registry_test

import (
	"fmt"
	"sync"
	"testing"

	"assessment/adapter/registry"
	"assessment/domain/model"
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

func TestSorterRegistryRegisterAndGet(t *testing.T) {

	reg := registry.NewSorterRegistry()

	mockSorter := NewMockSorter("MockSorter")

	reg.RegisterSorter(mockSorter)

	sorter, exists := reg.GetSorter("MockSorter")

	if !exists {
		t.Error("Sorter was not registered")
	}

	if sorter.Name() != "MockSorter" {
		t.Errorf("Sorter name mismatch: got %s, want %s", sorter.Name(), "MockSorter")
	}

	_, exists = reg.GetSorter("NonExistentSorter")

	if exists {
		t.Error("Non-existent sorter was found")
	}
}

func TestSorterRegistryGetAllSorters(t *testing.T) {

	reg := registry.NewSorterRegistry()

	mockSorter1 := NewMockSorter("MockSorter1")
	mockSorter2 := NewMockSorter("MockSorter2")
	mockSorter3 := NewMockSorter("MockSorter3")

	reg.RegisterSorter(mockSorter1)
	reg.RegisterSorter(mockSorter2)
	reg.RegisterSorter(mockSorter3)

	sorters := reg.GetAllSorters()

	if len(sorters) != 3 {
		t.Errorf("Sorter count mismatch: got %d, want %d", len(sorters), 3)
	}

	sorterMap := make(map[string]bool)
	for _, s := range sorters {
		sorterMap[s.Name()] = true
	}

	if !sorterMap["MockSorter1"] {
		t.Error("MockSorter1 not found in GetAllSorters result")
	}
	if !sorterMap["MockSorter2"] {
		t.Error("MockSorter2 not found in GetAllSorters result")
	}
	if !sorterMap["MockSorter3"] {
		t.Error("MockSorter3 not found in GetAllSorters result")
	}
}

func TestSorterRegistryUnregisterSorter(t *testing.T) {

	reg := registry.NewSorterRegistry()

	mockSorter := NewMockSorter("MockSorter")

	reg.RegisterSorter(mockSorter)

	success := reg.UnregisterSorter("MockSorter")

	if !success {
		t.Error("Sorter unregistration failed")
	}

	_, exists := reg.GetSorter("MockSorter")
	if exists {
		t.Error("Sorter still exists after unregistration")
	}

	success = reg.UnregisterSorter("NonExistentSorter")

	if success {
		t.Error("Unregistration of non-existent sorter succeeded")
	}
}

func TestSorterRegistryConcurrentAccess(t *testing.T) {

	reg := registry.NewSorterRegistry()

	numOperations := 100

	var wg sync.WaitGroup
	wg.Add(numOperations * 3)

	for i := 0; i < numOperations; i++ {
		go func(id int) {

			sorterName := fmt.Sprintf("MockSorter%d", id)

			reg.RegisterSorter(NewMockSorter(sorterName))
			wg.Done()

			sorter, exists := reg.GetSorter(sorterName)
			if !exists || sorter.Name() != sorterName {
				t.Errorf("Concurrent get failed for %s", sorterName)
			}
			wg.Done()

			success := reg.UnregisterSorter(sorterName)
			if !success {
				t.Errorf("Concurrent unregister failed for %s", sorterName)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	sorters := reg.GetAllSorters()
	if len(sorters) != 0 {
		t.Errorf("Registry not empty after concurrent operations: got %d sorters", len(sorters))
	}
}
