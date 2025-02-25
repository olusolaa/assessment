package usecase

import (
	"assessment/domain/model"
	"math"
)

type PaginationOptions struct {
	Page     int
	PageSize int
}

type PaginatedResult struct {
	Items      model.ProductList
	Page       int
	PageSize   int
	TotalItems int
	TotalPages int
	HasNext    bool
	HasPrev    bool
}

func (ps *ProductSorterUseCase) SortAndPaginateProducts(
	products model.ProductList,
	sorterName string,
	options PaginationOptions,
) (*PaginatedResult, error) {

	sortedProducts, err := ps.SortProducts(products, sorterName)
	if err != nil {
		return nil, err
	}

	return paginate(sortedProducts, options), nil
}

func paginate(products model.ProductList, options PaginationOptions) *PaginatedResult {
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

	if endIndex > totalItems {
		endIndex = totalItems
	}

	var pageItems model.ProductList
	if startIndex < totalItems {
		pageItems = products[startIndex:endIndex]
	} else {
		pageItems = model.ProductList{}
	}

	return &PaginatedResult{
		Items:      pageItems,
		Page:       options.Page,
		PageSize:   options.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    options.Page < totalPages,
		HasPrev:    options.Page > 1,
	}
}
