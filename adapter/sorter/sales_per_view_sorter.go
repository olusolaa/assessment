package sorter

import (
	"sort"

	"assessment/domain/model"
)

type SalesPerViewSorter struct {
	ascending bool
}

func NewSalesPerViewSorter(ascending bool) *SalesPerViewSorter {
	return &SalesPerViewSorter{
		ascending: ascending,
	}
}

func (s *SalesPerViewSorter) Sort(products model.ProductList) model.ProductList {

	result := products.Clone()

	sort.Slice(result, func(i, j int) bool {

		spv1 := calculateSalesPerView(result[i])
		spv2 := calculateSalesPerView(result[j])

		if s.ascending {
			return spv1 < spv2
		}
		return spv1 > spv2
	})

	return result
}

func calculateSalesPerView(p *model.Product) float64 {
	if p.ViewsCount == 0 {
		return 0
	}
	return float64(p.SalesCount) / float64(p.ViewsCount)
}

func (s *SalesPerViewSorter) Name() string {
	if s.ascending {
		return "Sales per View (ascending)"
	}
	return "Sales per View (descending)"
}
