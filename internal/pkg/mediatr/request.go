package mediatr

import (
	"clean_architecture_go/internal/pkg/error"
	"context"
)

type Request[TResponse any] interface{}

type RequestHandler[TRequest Request[TResponse], TResponse any] interface {
	Handle(ctx context.Context, request TRequest) (*TResponse, *error.RequestError)
}
