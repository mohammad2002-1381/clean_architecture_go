// domain/product.go
package domain

type Product struct {
	BaseEntity[int32]
	Name        string  `gorm:"column:name;not null" json:"name"`
	Description string  `gorm:"column:description" json:"description"`
	Price       float64 `gorm:"column:price;not null" json:"price"`
	Stock       int32   `gorm:"column:stock;default:0" json:"stock"`
	IsActive    bool    `gorm:"column:is_active;default:true" json:"is_active"`
}

func NewProduct(name, description string, price float64, stock int32) *Product {
	return &Product{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		IsActive:    true,
	}
}
