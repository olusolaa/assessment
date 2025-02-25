package registry

import (
	"sync"

	"assessment/domain/service"
)

type SorterRegistry struct {
	sorters map[string]service.Sorter
	mutex   sync.RWMutex
}

func NewSorterRegistry() *SorterRegistry {
	return &SorterRegistry{
		sorters: make(map[string]service.Sorter),
	}
}

func (r *SorterRegistry) RegisterSorter(sorter service.Sorter) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.sorters[sorter.Name()] = sorter
}

func (r *SorterRegistry) GetSorter(name string) (service.Sorter, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	sorter, exists := r.sorters[name]
	return sorter, exists
}

func (r *SorterRegistry) GetAllSorters() []service.Sorter {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	sorters := make([]service.Sorter, 0, len(r.sorters))
	for _, sorter := range r.sorters {
		sorters = append(sorters, sorter)
	}

	return sorters
}

func (r *SorterRegistry) UnregisterSorter(name string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.sorters[name]
	if exists {
		delete(r.sorters, name)
		return true
	}

	return false
}
