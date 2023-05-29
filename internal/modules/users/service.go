package users

import (
	"context"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, user *UserWithPassword) error {
	return s.repo.Create(ctx, user)
}

func (s *Service) GetUserWithPassword(ctx context.Context, email string) (*UserWithPassword, error) {
	return s.repo.GetUserWithPassword(ctx, email)
}

func (s *Service) GetUserById(ctx context.Context, id int) (*User, error) {
	return s.repo.GetUserById(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int, bio string) error {
	return s.repo.Update(ctx, id, bio)
}

func (s *Service) UpdateUserRole(ctx context.Context, id int, roleName string) error {
	return s.repo.UpdateUserRole(ctx, id, roleName)
}
