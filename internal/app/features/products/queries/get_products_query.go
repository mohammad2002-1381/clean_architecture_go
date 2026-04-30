// app/features/products/queries/get_products_query.go
package queries

import (
	"clean_architecture_go/internal/app/features/products/dto"
	"clean_architecture_go/internal/domain"
	"clean_architecture_go/internal/domain/extensions"
	"clean_architecture_go/internal/infra"
	errs "clean_architecture_go/internal/pkg/error"
	"context"
	"net/http"
)

type GetProductsQuery struct {
	extensions.PagedQuery
}

type GetProductsQueryHandler struct {
	BaseRepo *infra.BaseRepository[domain.Product, int32]
}

func NewGetProductsQueryHandler(baseRepo *infra.BaseRepository[domain.Product, int32]) *GetProductsQueryHandler {
	return &GetProductsQueryHandler{
		BaseRepo: baseRepo,
	}
}

func (h *GetProductsQueryHandler) Handle(ctx context.Context, req GetProductsQuery) (*extensions.PagedResult[*dto.ProductDTO], *errs.RequestError) {
	// Build conditions
	condition := "is_active = ?"
	args := []interface{}{true}

	// if req.Search != "" {
	// 	condition += " AND (name LIKE ? OR description LIKE ?)"
	// 	args = append(args, "%"+req.Search+"%", "%"+req.Search+"%")
	// }

	// if req.MinPrice > 0 {
	// 	condition += " AND price >= ?"
	// 	args = append(args, req.MinPrice)
	// }

	// if req.MaxPrice > 0 {
	// 	condition += " AND price <= ?"
	// 	args = append(args, req.MaxPrice)
	// }

	// if req.CategoryID > 0 {
	// 	condition += " AND category_id = ?"
	// 	args = append(args, req.CategoryID)
	// }

	// Use repository only - no direct DB
	result, err := h.BaseRepo.GetPaged(ctx, req.PagedQuery, condition, args...)
	if err != nil {
		return nil, errs.NewRequestError(http.StatusInternalServerError, err.Error())
	}

	// Map to DTOs
	dtos := make([]*dto.ProductDTO, len(result.Items))
	for i, p := range result.Items {
		dtos[i] = &dto.ProductDTO{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
			IsActive:    p.IsActive,
		}
	}
	pagedResult := extensions.NewPagedResult(dtos, result.TotalCount, req.PagedQuery)
	return &pagedResult, nil

}
