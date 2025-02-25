package model

import (
	"fmt"
	"time"
)

type Product struct {
	ID         int
	Name       string
	Price      float64
	Created    time.Time
	SalesCount int
	ViewsCount int
}

type ProductList []*Product

func ParseTime(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func (pl ProductList) Clone() ProductList {
	cloned := make(ProductList, len(pl))
	for i, p := range pl {

		cloned[i] = &Product{
			ID:         p.ID,
			Name:       p.Name,
			Price:      p.Price,
			Created:    p.Created,
			SalesCount: p.SalesCount,
			ViewsCount: p.ViewsCount,
		}
	}
	return cloned
}

func (p *Product) String() string {
	salesPerView := float64(0)
	if p.ViewsCount > 0 {
		salesPerView = float64(p.SalesCount) / float64(p.ViewsCount)
	}
	return fmt.Sprintf("ID: %d, Name: %s, Price: $%.2f, Created: %s, Sales/View: %.6f",
		p.ID, p.Name, p.Price, p.Created.Format("2006-01-02"), salesPerView)
}
