package main

import (
	"fmt"
	"os"

	"assessment/adapter/registry"
	"assessment/adapter/sorter"
	"assessment/domain/model"
	"assessment/infrastructure/config"
	"assessment/infrastructure/persistence"
	"assessment/usecase"
)

func main() {
	productData := []map[string]interface{}{
		{
			"id":          1,
			"name":        "Alabaster Table",
			"price":       12.99,
			"created":     "2019-01-04",
			"sales_count": 32,
			"views_count": 730,
		},
		{
			"id":          2,
			"name":        "Zebra Table",
			"price":       44.49,
			"created":     "2012-01-04",
			"sales_count": 301,
			"views_count": 3279,
		},
		{
			"id":          3,
			"name":        "Coffee Table",
			"price":       10.00,
			"created":     "2014-05-28",
			"sales_count": 1048,
			"views_count": 20123,
		},
	}

	products := make(model.ProductList, 0, len(productData))
	for _, data := range productData {
		createdStr := data["created"].(string)
		created, err := model.ParseTime(createdStr)
		if err != nil {
			fmt.Printf("Error parsing date %s: %v\n", createdStr, err)
			continue
		}

		product := &model.Product{
			ID:         data["id"].(int),
			Name:       data["name"].(string),
			Price:      data["price"].(float64),
			Created:    created,
			SalesCount: data["sales_count"].(int),
			ViewsCount: data["views_count"].(int),
		}
		products = append(products, product)
	}

	repo := persistence.NewInMemoryProductRepository()

	if err := repo.Save(products); err != nil {
		fmt.Printf("Error saving products: %v\n", err)
		return
	}

	cfg := config.NewConfig()
	configFile := "infrastructure/config/sample_config.json"
	if _, err := os.Stat(configFile); err == nil {
		if err := cfg.LoadFromFile(configFile); err != nil {
			fmt.Printf("Warning: Failed to load config file: %v\n", err)
		} else {
			fmt.Println("Configuration loaded from", configFile)
		}
	} else {
		fmt.Println("Using default configuration")
	}

	reg := registry.NewSorterRegistry()

	sorterUseCase := usecase.NewProductSorterUseCase(reg)

	sorterUseCase.SetConfig(cfg)

	sorter.InitializeDefaultSorters(reg, cfg)

	reg.RegisterSorter(sorter.NewNameSorter(true))

	fmt.Println("Available sorters:")
	for _, name := range sorterUseCase.GetAvailableSorters() {
		fmt.Printf("- %s\n", name)
	}
	fmt.Println()

	repoProducts, err := repo.GetAll()
	if err != nil {
		fmt.Printf("Error getting products: %v\n", err)
		return
	}

	priceSortedProducts, err := sorterUseCase.SortProducts(repoProducts, "Price (ascending)")
	if err != nil {
		fmt.Printf("Error sorting by price: %v\n", err)
	} else {
		fmt.Println("Products sorted by price (ascending):")
		printProducts(priceSortedProducts)
	}

	spvSortedProducts, err := sorterUseCase.SortProducts(repoProducts, "Sales per View (descending)")
	if err != nil {
		fmt.Printf("Error sorting by sales per view: %v\n", err)
	} else {
		fmt.Println("\nProducts sorted by sales per view (descending):")
		printProducts(spvSortedProducts)
	}

	dateSortedProducts, err := sorterUseCase.SortProducts(repoProducts, "Creation Date (ascending)")
	if err != nil {
		fmt.Printf("Error sorting by creation date: %v\n", err)
	} else {
		fmt.Println("\nProducts sorted by creation date (ascending):")
		printProducts(dateSortedProducts)
	}

	nameSortedProducts, err := sorterUseCase.SortProducts(repoProducts, "Name (ascending)")
	if err != nil {
		fmt.Printf("Error sorting by name: %v\n", err)
	} else {
		fmt.Println("\nProducts sorted by name (ascending):")
		printProducts(nameSortedProducts)
	}

	fmt.Println("\nDemonstrating pagination (page 1, 2 items per page):")
	paginationOptions := usecase.PaginationOptions{Page: 1, PageSize: 2}
	paginatedResult, err := sorterUseCase.SortAndPaginateProducts(repoProducts, "Price (ascending)", paginationOptions)
	if err != nil {
		fmt.Printf("Error paginating products: %v\n", err)
	} else {
		fmt.Printf("Page %d of %d (Total items: %d)\n",
			paginatedResult.Page, paginatedResult.TotalPages, paginatedResult.TotalItems)
		printProducts(paginatedResult.Items)

		if paginatedResult.HasNext {
			fmt.Println("Has next page: Yes")
		} else {
			fmt.Println("Has next page: No")
		}

		if paginatedResult.HasPrev {
			fmt.Println("Has previous page: Yes")
		} else {
			fmt.Println("Has previous page: No")
		}
	}
}

func printProducts(products model.ProductList) {
	for _, p := range products {
		fmt.Println(p.String())
	}
}
