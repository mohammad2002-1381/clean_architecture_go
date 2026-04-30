package main

import (
	"fmt"
	"time"
)

type BaseEntity[TID comparable] struct {
	ID        TID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IBaseEntity[TID comparable] interface {
	Get()
}

func (b BaseEntity[TID]) Get() {}

type User struct {
	BaseEntity[int32]

	FirstName string
	LastName  string
	Email     string
}

type T struct {
	i int32
}

type BaseRepository[TEntity IBaseEntity[TID], TID comparable] struct {
}

func NewRepository[TEntity IBaseEntity[TID], TID comparable]() *BaseRepository[TEntity, TID] {
	return &BaseRepository[TEntity, TID]{}
}

func main() {
	v := NewRepository[T, int32]()
	fmt.Println(v)
}
