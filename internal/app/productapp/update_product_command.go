package productapp

import (
	"context"
	"go-ca/internal/app/service"
	"go-ca/internal/domain"
)

type UpdateProductCommand struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Price       float64 `json:"price"`
}

type UpdateProductCommandHandler struct {
	productRepo        domain.BaseRepository[*domain.Product, uint]
	currentUserService service.CurrentUserService
}

func NewUpdateProductCommandHandler(
	productRepo domain.BaseRepository[*domain.Product, uint],
	currentUserService service.CurrentUserService,
) UpdateProductCommandHandler {
	return UpdateProductCommandHandler{
		productRepo: productRepo,
	}
}

func (c *UpdateProductCommandHandler) Handle(ctx context.Context, request UpdateProductCommand) (*ProductDTO, error) {
	userID, e := c.currentUserService.GetUserID(ctx)

	if e != nil {
		return nil, e
	}

	p, err := c.productRepo.Where(ctx, &domain.Product{BaseEntity: domain.BaseEntity[uint]{
		ID: request.ID,
	}, UserID: userID}).First(ctx)

	if err != nil {
		return nil, err
	}
	
	p.SetName(request.Name)
	p.SetDescription(request.Description)
	p.SetPrice(request.Price)

	err = c.productRepo.Update(ctx, p)
	if err != nil {
		return nil, err
	}

	dto := NewProductDTO(p)
	return &dto, nil
}
