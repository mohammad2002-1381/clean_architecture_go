package app

import "go-ca/internal/domain"

func MapPaginatedResult[S any, T any](source *domain.PaginatedResult[S], mapper func(*S) T) *domain.PaginatedResult[T] {
	if source == nil {
		return nil
	}

	dtos := make([]T, 0, len(source.Items))
	for i := range source.Items {
		// Use index to safely get the pointer of the item
		dtos = append(dtos, mapper(&source.Items[i]))
	}

	return &domain.PaginatedResult[T]{
		Items:      dtos,
		PageSize:   source.PageSize,
		PageNumber: source.PageNumber,
		TotalCount: source.TotalCount,
	}
}
