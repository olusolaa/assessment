package repository

import (
	"assessment/domain/model"
)


type ProductRepository interface {
	
	GetAll() (model.ProductList, error)
	
	
	GetByIDs(ids []int) (model.ProductList, error)
	
	
	Save(products model.ProductList) error
} 