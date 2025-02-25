package services

import (
	"fmt"
	"assessment/config"
	"assessment/domain/model"
	"assessment/adapter/registry"
	"assessment/adapter/sorter"
)


type ProductService struct {
	registry *registry.SorterRegistry
	config   *config.Config
}


func NewProductService(registry *registry.SorterRegistry) *ProductService {
	return &ProductService{
		registry: registry,
		config:   config.NewConfig(),
	}
}


func (s *ProductService) SetConfig(cfg *config.Config) {
	s.config = cfg
}


func (s *ProductService) GetConfig() *config.Config {
	return s.config
}



func (s *ProductService) SortProducts(products model.ProductList, sorterName string) (model.ProductList, error) {
	sorter, exists := s.registry.GetSorter(sorterName)
	if !exists {
		return nil, fmt.Errorf("sorter %s not found", sorterName)
	}

	return sorter.Sort(products), nil
}


func (s *ProductService) SortAndPaginateProducts(
	products model.ProductList, 
	sorterName string, 
	options PaginationOptions,
) (PaginatedResult, error) {
	sortedProducts, err := s.SortProducts(products, sorterName)
	if err != nil {
		return PaginatedResult{}, err
	}

	return PaginateProducts(sortedProducts, options), nil
}


func (s *ProductService) GetAvailableSorters() []string {
	sorters := s.registry.GetAllSorters()
	names := make([]string, len(sorters))
	
	for i, sorter := range sorters {
		names[i] = sorter.Name()
	}
	
	return names
}


func (s *ProductService) GetEnabledSorters() []string {
	allSorters := s.registry.GetAllSorters()
	var enabledNames []string
	
	for _, sorter := range allSorters {
		name := sorter.Name()
		
		for _, cfg := range s.config.Sorters {
			if cfg.Name == name && cfg.Enabled {
				enabledNames = append(enabledNames, name)
				break
			}
		}
	}
	
	return enabledNames
}



func (s *ProductService) CalculateSalesPerViewRatio(product *model.Product) float64 {
	if product.ViewsCount == 0 {
		return 0
	}
	return float64(product.SalesCount) / float64(product.ViewsCount)
}


func (s *ProductService) InitializeDefaultSorters() {
	
	priceConfig, priceExists := s.config.GetSorterConfig("price")
	salesConfig, salesExists := s.config.GetSorterConfig("sales_per_view")
	dateConfig, dateExists := s.config.GetSorterConfig("creation_date")
	nameConfig, nameExists := s.config.GetSorterConfig("name")
	
	
	if priceExists && priceConfig.Enabled {
		s.registry.RegisterSorter(sorter.NewPriceSorter(priceConfig.Ascending))
		
		s.registry.RegisterSorter(sorter.NewPriceSorter(!priceConfig.Ascending))
	}
	
	
	if salesExists && salesConfig.Enabled {
		s.registry.RegisterSorter(sorter.NewSalesPerViewSorter(salesConfig.Ascending))
		s.registry.RegisterSorter(sorter.NewSalesPerViewSorter(!salesConfig.Ascending))
	}
	
	
	if dateExists && dateConfig.Enabled {
		s.registry.RegisterSorter(sorter.NewDateSorter(dateConfig.Ascending))
		s.registry.RegisterSorter(sorter.NewDateSorter(!dateConfig.Ascending))
	}
	
	
	if nameExists && nameConfig.Enabled {
		s.registry.RegisterSorter(sorter.NewNameSorter(nameConfig.Ascending))
		s.registry.RegisterSorter(sorter.NewNameSorter(!nameConfig.Ascending))
	}
} 
