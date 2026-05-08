package productapp

import (
	"context"
	"go-ca/internal/app"
	"go-ca/internal/app/service"
	"go-ca/internal/domain"
)

type GetProductListQuery struct {
	domain.PaginationParams
}

type GetProductListQueryHandler struct {
	productRepo        domain.BaseRepository[*domain.Product, uint]
	currentUserService service.CurrentUserService
}

func NewGetProductListQueryHandler(
	productRepo domain.BaseRepository[*domain.Product, uint],
	currentUserService service.CurrentUserService,
) GetProductListQueryHandler {
	return GetProductListQueryHandler{
		productRepo:        productRepo,
		currentUserService: currentUserService,
	}
}

func (q *GetProductListQueryHandler) Handle(ctx context.Context, request GetProductListQuery) (*domain.PaginatedResult[ProductDTO], error) {
	userID, e := q.currentUserService.GetUserID(ctx)

	if e != nil {
		return nil, e
	}

	products, err := q.productRepo.Where(ctx, &domain.Product{
		UserID: userID,
	}).GetPaged(ctx, request.PaginationParams)

	if err != nil {
		return nil, err
	}

	return app.MapPaginatedResult(&products, NewProductDTO), nil
}
