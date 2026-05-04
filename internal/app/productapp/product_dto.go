package productapp

import (
	"time"

	"go-ca/internal/domain"
)

type ProductDTO struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Price       float64    `json:"price"`
	Items       []*ItemDto `json:"items"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
}

func NewProductDTO(product *domain.Product) ProductDTO {
	var itemsDto []*ItemDto

	for _, item := range product.Items {
		itemsDto = append(itemsDto, &ItemDto{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   item.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   item.UpdatedAt.Format(time.RFC3339),
		})
	}

	return ProductDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Items:       itemsDto, // <-- FIX: Assign the slice to the DTO
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}
}
