package productapp

import (
	"context"

	"go-ca/internal/app/notification"
	"go-ca/internal/app/service"
	"go-ca/internal/domain"
)

type CreateProductCommand struct {
	Name        string               `json:"name"`
	Description *string              `json:"description"`
	Price       float64              `json:"price"`
	Items       []*CreateItemCommand `json:"items"`
}

type CreateItemCommand struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type CreateProductCommandHandler struct {
	productRepo        domain.BaseRepository[domain.Product, uint]
	currentUserService service.CurrentUserService
}

func NewCreateProductCommandHandler(
	productRepo domain.BaseRepository[domain.Product, uint],
	currentUserService service.CurrentUserService,
) CreateProductCommandHandler {
	return CreateProductCommandHandler{
		productRepo:        productRepo,
		currentUserService: currentUserService,
	}
}

func (c *CreateProductCommandHandler) Handle(ctx context.Context, request CreateProductCommand) (*ProductDTO, error) {
	event := notification.UserRegisteredEvent{
		NewEmail: "test",
	}

	userID, e := c.currentUserService.GetUserID(ctx)

	if e != nil {
		return nil, e
	}

	p := domain.NewProduct(request.Name,
		request.Description,
		request.Price,
		userID,
		&event,
	)

	var domainItems []*domain.Item
	for _, reqItem := range request.Items {
		newItem := domain.NewItem(reqItem.Name, reqItem.Description)
		domainItems = append(domainItems, &newItem)
	}

	p.SetItems(domainItems)

	err := c.productRepo.Add(ctx, &p)

	if err != nil {
		return nil, err
	}

	dto := NewProductDTO(&p)

	return &dto, nil
}
