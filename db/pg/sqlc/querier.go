// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"context"
)

type Querier interface {
	AddContact(ctx context.Context, arg AddContactParams) (Contact, error)
	AddExamine(ctx context.Context, arg AddExamineParams) (Examine, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteExamine(ctx context.Context, ownerID int64) error
	DeleteUser(ctx context.Context, id int64) error
	GetExamine(ctx context.Context, ownerID int64) ([]Examine, error)
	GetUser(ctx context.Context, username string) (User, error)
	LoginUser(ctx context.Context, arg LoginUserParams) (User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
