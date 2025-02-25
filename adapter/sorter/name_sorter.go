package sorter

import (
	"sort"
	"strings"
	
	"assessment/domain/model"
)


type NameSorter struct {
	ascending bool
}


func NewNameSorter(ascending bool) *NameSorter {
	return &NameSorter{
		ascending: ascending,
	}
}


func (s *NameSorter) Sort(products model.ProductList) model.ProductList {
	
	result := products.Clone()
	
	
	sort.Slice(result, func(i, j int) bool {
		if s.ascending {
			return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
		}
		return strings.ToLower(result[i].Name) > strings.ToLower(result[j].Name)
	})
	
	return result
}


func (s *NameSorter) Name() string {
	if s.ascending {
		return "Name (ascending)"
	}
	return "Name (descending)"
} 
