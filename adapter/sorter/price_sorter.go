package sorter

import (
	"sort"

	"assessment/domain/model"
)

type PriceSorter struct {
	ascending bool
}

func NewPriceSorter(ascending bool) *PriceSorter {
	return &PriceSorter{
		ascending: ascending,
	}
}

func (s *PriceSorter) Sort(products model.ProductList) model.ProductList {

	result := products.Clone()

	sort.Slice(result, func(i, j int) bool {
		if s.ascending {
			return result[i].Price < result[j].Price
		}
		return result[i].Price > result[j].Price
	})

	return result
}

func (s *PriceSorter) Name() string {
	if s.ascending {
		return "Price (ascending)"
	}
	return "Price (descending)"
}
