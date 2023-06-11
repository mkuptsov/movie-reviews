package stars

import (
	"context"

	"github.com/cloudmachinery/movie-reviews/internal/log"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateStar(ctx context.Context, star *StarDetails) error {
	err := s.repo.CreateStar(ctx, star)
	if err != nil {
		return nil
	}

	logger := log.FromContext(ctx)
	logger.Info("star created",
		"first_name", star.FirstName,
		"last_name", star.LastName)

	return nil
}

func (s *Service) GetStarByID(ctx context.Context, id int) (*StarDetails, error) {
	return s.repo.GetStarByID(ctx, id)
}

func (s *Service) GetAllPaginated(ctx context.Context, offset, limit int) ([]*Star, int, error) {
	return s.repo.GetAllPaginated(ctx, offset, limit)
}

func (s *Service) UpdateStar(ctx context.Context, id int, star *StarDetails) error {
	err := s.repo.UpdateStar(ctx, id, star)
	if err != nil {
		return err
	}

	logger := log.FromContext(ctx)
	logger.Info("star updated",
		"first_name", star.FirstName,
		"last_name", star.LastName)

	return nil
}

func (s *Service) DeleteStar(ctx context.Context, id int) error {
	err := s.repo.DeleteStar(ctx, id)
	if err != nil {
		return err
	}

	logger := log.FromContext(ctx)
	logger.Info("star deleted",
		"star_id", id,
	)

	return nil
}
