package services

import (
	"assessment/domain/model"
	"math"
)

type PaginationOptions struct {
	Page     int
	PageSize int
}

func NewPaginationOptions() PaginationOptions {
	return PaginationOptions{
		Page:     1,
		PageSize: 10,
	}
}

type PaginatedResult struct {
	Items      model.ProductList
	TotalItems int
	TotalPages int
	Page       int
	PageSize   int
	HasNext    bool
	HasPrev    bool
}

func PaginateProducts(products model.ProductList, options PaginationOptions) PaginatedResult {
	totalItems := len(products)

	if options.Page < 1 {
		options.Page = 1
	}
	if options.PageSize < 1 {
		options.PageSize = 10
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(options.PageSize)))

	if options.Page > totalPages && totalPages > 0 {
		options.Page = totalPages
	}

	startIndex := (options.Page - 1) * options.PageSize
	endIndex := startIndex + options.PageSize

	if startIndex >= totalItems {
		startIndex = 0
		endIndex = 0
	}
	if endIndex > totalItems {
		endIndex = totalItems
	}

	var pageItems model.ProductList
	if startIndex < endIndex {
		pageItems = products[startIndex:endIndex]
	} else {
		pageItems = model.ProductList{}
	}

	return PaginatedResult{
		Items:      pageItems,
		TotalItems: totalItems,
		TotalPages: totalPages,
		Page:       options.Page,
		PageSize:   options.PageSize,
		HasNext:    options.Page < totalPages,
		HasPrev:    options.Page > 1,
	}
}
