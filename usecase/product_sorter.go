package usecase

import (
	"fmt"

	"assessment/domain/model"
	"assessment/domain/service"
	"assessment/infrastructure/config"
)

type ProductSorterUseCase struct {
	registry service.SorterRegistry
	config   *config.Config
}

func NewProductSorterUseCase(registry service.SorterRegistry) *ProductSorterUseCase {
	return &ProductSorterUseCase{
		registry: registry,
		config:   config.NewConfig(),
	}
}

func (ps *ProductSorterUseCase) SetConfig(cfg *config.Config) {
	ps.config = cfg
}

func (ps *ProductSorterUseCase) GetConfig() *config.Config {
	return ps.config
}

func (ps *ProductSorterUseCase) GetRegistry() service.SorterRegistry {
	return ps.registry
}

func (ps *ProductSorterUseCase) SortProducts(products model.ProductList, sorterName string) (model.ProductList, error) {
	sorter, exists := ps.registry.GetSorter(sorterName)
	if !exists {
		return nil, fmt.Errorf("sorter not found: %s", sorterName)
	}

	if !ps.isSorterEnabled(sorterName) {
		return nil, fmt.Errorf("sorter is disabled: %s", sorterName)
	}

	return sorter.Sort(products), nil
}

func (ps *ProductSorterUseCase) GetAvailableSorters() []string {
	sorters := ps.registry.GetAllSorters()
	names := make([]string, 0, len(sorters))

	for _, sorter := range sorters {
		if ps.isSorterEnabled(sorter.Name()) {
			names = append(names, sorter.Name())
		}
	}

	return names
}

func (ps *ProductSorterUseCase) isSorterEnabled(sorterName string) bool {
	if ps.config == nil {
		return true
	}

	for _, disabled := range ps.config.DisabledSorters {
		if disabled == sorterName {
			return false
		}
	}

	return true
}
