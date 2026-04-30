package extentions

import (
	"clean_architecture_go/internal/domain"
)

type Mapper[E domain.BaseEntity[any], T comparable] interface {
	Map(input *E) T
}
