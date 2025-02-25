package service

import (
	"assessment/domain/model"
)

type Sorter interface {
	Sort(products model.ProductList) model.ProductList

	Name() string
}

type SorterRegistry interface {
	RegisterSorter(sorter Sorter)

	GetSorter(name string) (Sorter, bool)

	GetAllSorters() []Sorter

	UnregisterSorter(name string) bool
}
