package user

import "context"

type Storage interface {
	Create(ctx context.Context, user CreateUserDTO) (string, error)
	FindOne(ctx context.Context, id string) (CreateUserDTO, error)
	FindRefresh(ctx context.Context, refToken string) (CreateUserDTO, error)
	Update(ctx context.Context, user CreateUserDTO) error
	Delete(ctx context.Context, id string) error
}
