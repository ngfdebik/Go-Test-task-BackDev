package service

import (
	"context"

	"example.com/m/internal/config"
	"example.com/m/internal/user"
	"example.com/m/internal/user/db"
	"example.com/m/pkg/api/client/mongodb"
)

type Service struct {
	storage user.Storage
}

func NewStorage() *user.Storage {
	cfg := config.GetConfig()
	mongoDBClient, err := mongodb.NewClient(context.Background(), cfg.MongoDB.Host, cfg.MongoDB.Port, cfg.MongoDB.Username,
		cfg.MongoDB.Password, cfg.MongoDB.Database, cfg.MongoDB.AuthDB)
	if err != nil {
		panic(err)
	}
	strg := db.NewStorage(mongoDBClient, cfg.MongoDB.Collection)
	return &strg
}
func NewService(storage user.Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) Create(ctx context.Context, dto user.CreateUserDTO) (str string, err error) {
	str, err = s.storage.Create(context.Background(), dto)
	if err != nil {
		return "", err
	}
	return str, nil
}

func (s *Service) GuidExist(ctx context.Context, guid string) (flag bool, u user.CreateUserDTO) {
	u, err := s.storage.FindOne(context.Background(), guid)
	if err != nil {
		return false, u
	}
	return true, u
}

func (s *Service) FindRefresh(ctx context.Context, refToken string) (u user.CreateUserDTO, err error) {
	u, err = s.storage.FindRefresh(ctx, refToken)
	if err != nil {
		return u, err
	}
	return u, nil
}
func (s *Service) Update(ctx context.Context, user user.CreateUserDTO) error {
	err := s.storage.Update(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
