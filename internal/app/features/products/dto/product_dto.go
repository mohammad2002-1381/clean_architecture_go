package dto

type ProductDTO struct {
	ID          int32   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int32   `json:"stock"`
	IsActive    bool    `json:"is_active"`
}
