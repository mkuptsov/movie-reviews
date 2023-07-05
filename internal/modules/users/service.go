package users

import (
	"context"

	"github.com/mkuptsov/movie-reviews/internal/log"
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
	logger.Info("user created",
		"email", user.Email)
	return nil
}

func (s *Service) GetUserWithPassword(ctx context.Context, email string) (*UserWithPassword, error) {
	return s.repo.GetUserWithPassword(ctx, email)
}

func (s *Service) GetUserByID(ctx context.Context, id int) (*User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *Service) GetUserByUserName(ctx context.Context, userName string) (*User, error) {
	return s.repo.GetUserByUserName(ctx, userName)
}

func (s *Service) DeleteUser(ctx context.Context, id int) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	logger := log.FromContext(ctx)
	logger.Info("user deleted",
		"user_id", id)
	return nil
}

func (s *Service) Update(ctx context.Context, id int, bio string) error {
	return s.repo.Update(ctx, id, bio)
}

func (s *Service) SetUserRole(ctx context.Context, id int, roleName string) error {
	err := s.repo.SetUserRole(ctx, id, roleName)
	if err != nil {
		return err
	}

	logger := log.FromContext(ctx)
	logger.Info("user role updated",
		"user_id", id,
		"new_role", roleName)
	return nil
}
