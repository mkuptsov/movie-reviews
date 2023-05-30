package users

import (
	"context"

	"github.com/cloudmachinery/movie-reviews/internal/modules/log"
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
	err := s.repo.Create(ctx, user)
	if err != nil {
		return err
	}

	logger := log.FromContext(ctx)
	logger.Info("user created")
	return nil
}

func (s *Service) GetUserWithPassword(ctx context.Context, email string) (*UserWithPassword, error) {
	return s.repo.GetUserWithPassword(ctx, email)
}

func (s *Service) GetUserById(ctx context.Context, id int) (*User, error) {
	return s.repo.GetUserById(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	logger := log.FromContext(ctx)
	logger.Info("user deleted")
	return nil
}

func (s *Service) Update(ctx context.Context, id int, bio string) error {
	return s.repo.Update(ctx, id, bio)
}

func (s *Service) UpdateUserRole(ctx context.Context, id int, roleName string) error {
	err := s.repo.UpdateUserRole(ctx, id, roleName)
	if err != nil {
		return err
	}

	logger := log.FromContext(ctx)
	logger.Info("user role updated")
	return nil
}
