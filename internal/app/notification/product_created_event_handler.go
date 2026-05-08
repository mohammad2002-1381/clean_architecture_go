package notification

import (
	"context"
	"fmt"
	"go-ca/internal/domain"
	"strconv"
)

type ProductCreatedEventHandler struct {
}

func NewCreateProductCommandHandler() ProductCreatedEventHandler {
	return ProductCreatedEventHandler{}
}

func (p *ProductCreatedEventHandler) Handle(ctx context.Context, event *domain.ProductCreatedEvent) error {
	productIdStr := strconv.Itoa(int(event.Product.ID))

	fmt.Printf("product created: %s - %s", productIdStr, event.Product.Name)

	return nil
}
