package productapp

import "go-ca/internal/domain"

type ItemDto struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	Description *string     `json:"description"`
	Product     *ProductDTO `json:"product"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
}

func NewItemDto(item *domain.Item) ItemDto {
	productDto := NewProductDTO(item.Product)
	return ItemDto{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Product:     &productDto,
		CreatedAt:   item.CreatedAt.Format(""),
		UpdatedAt:   item.UpdatedAt.Format(""),
	}
}
