package sorter

import (
	"sort"

	"assessment/domain/model"
)

type DateSorter struct {
	ascending bool
}

func NewDateSorter(ascending bool) *DateSorter {
	return &DateSorter{
		ascending: ascending,
	}
}

func (s *DateSorter) Sort(products model.ProductList) model.ProductList {

	result := products.Clone()

	sort.Slice(result, func(i, j int) bool {
		if s.ascending {
			return result[i].Created.Before(result[j].Created)
		}
		return result[i].Created.After(result[j].Created)
	})

	return result
}

func (s *DateSorter) Name() string {
	if s.ascending {
		return "Creation Date (ascending)"
	}
	return "Creation Date (descending)"
}
