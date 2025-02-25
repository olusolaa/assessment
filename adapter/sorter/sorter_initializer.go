package sorter

import (
	"assessment/domain/service"
	"assessment/infrastructure/config"
)


func InitializeDefaultSorters(registry service.SorterRegistry, cfg *config.Config) {
	
	registry.RegisterSorter(NewPriceSorter(true))
	registry.RegisterSorter(NewPriceSorter(false))
	
	
	registry.RegisterSorter(NewDateSorter(true))
	registry.RegisterSorter(NewDateSorter(false))
	
	
	registry.RegisterSorter(NewNameSorter(true))
	registry.RegisterSorter(NewNameSorter(false))
	
	
	registry.RegisterSorter(NewSalesPerViewSorter(true))
	registry.RegisterSorter(NewSalesPerViewSorter(false))
} 
