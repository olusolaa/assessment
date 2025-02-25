package persistence

import (
	"sync"
	
	"assessment/domain/model"
)


type InMemoryProductRepository struct {
	products model.ProductList
	mutex    sync.RWMutex
}


func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make(model.ProductList, 0),
	}
}


func (r *InMemoryProductRepository) GetAll() (model.ProductList, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	
	return r.products.Clone(), nil
}


func (r *InMemoryProductRepository) GetByIDs(ids []int) (model.ProductList, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	
	idMap := make(map[int]bool)
	for _, id := range ids {
		idMap[id] = true
	}
	
	
	result := make(model.ProductList, 0)
	for _, product := range r.products {
		if idMap[product.ID] {
			
			result = append(result, &model.Product{
				ID:         product.ID,
				Name:       product.Name,
				Price:      product.Price,
				Created:    product.Created,
				SalesCount: product.SalesCount,
				ViewsCount: product.ViewsCount,
			})
		}
	}
	
	return result, nil
}


func (r *InMemoryProductRepository) Save(products model.ProductList) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	
	r.products = products.Clone()
	
	return nil
} 
