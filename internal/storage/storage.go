package storage

import (
	//"context"
	"errors"
)

//type Repository interface {
//	Create(ctx context.Context, author *Author) error
//	FindAll(ctx context.Context) (u []Author, err error)
//	FindOne(ctx context.Context, id string) (Author, error)
//	Update(ctx context.Context, user Author) error
//	Delete(ctx context.Context, id string) error
//}

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)
