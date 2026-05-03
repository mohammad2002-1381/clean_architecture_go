package productapp

import (
	"context"
	"go-ca/internal/app/service"
	"go-ca/internal/domain"
)

type GetProductByIdQuery struct {
	ID uint
}

type GetProductByIdQueryHandler struct {
	productRepo        domain.BaseRepository[domain.Product, uint]
	currentUserService service.CurrentUserService
}

func NewGetProductByIdQueryHandler(
	productRepo domain.BaseRepository[domain.Product, uint],
	currentUserService service.CurrentUserService,
) GetProductByIdQueryHandler {
	return GetProductByIdQueryHandler{
		productRepo:        productRepo,
		currentUserService: currentUserService,
	}
}

func (q *GetProductByIdQueryHandler) Handle(ctx context.Context, request GetProductByIdQuery) (*ProductDto, error) {
	userID, err := q.currentUserService.GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	p, err := q.productRepo.Where(ctx, &domain.Product{
		UserID: userID,
		BaseEntity: domain.BaseEntity[uint]{
			ID: request.ID,
		},
	}).First(ctx)

	if err != nil {
		return nil, err
	}

	dto := NewProductDto(p)
	return &dto, nil
}
