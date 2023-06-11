package movies

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

func (s *Service) CreateMovie(ctx context.Context, movie *MovieDetails) error {
	err := s.repo.CreateMovie(ctx, movie)
	if err != nil {
		return nil
	}

	logger := log.FromContext(ctx)
	logger.Info("movie created",
		"movie_title", movie.Title)

	return nil
}

func (s *Service) GetMovieByID(ctx context.Context, id int) (*MovieDetails, error) {
	return s.repo.GetMovieByID(ctx, id)
}

func (s *Service) GetAllPaginated(ctx context.Context, offset, limit int) ([]*Movie, int, error) {
	return s.repo.GetAllPaginated(ctx, offset, limit)
}

func (s *Service) UpdateMovie(ctx context.Context, id int, movie *MovieDetails) error {
	err := s.repo.UpdateMovie(ctx, id, movie)
	if err != nil {
		return err
	}

	logger := log.FromContext(ctx)
	logger.Info("movie updated",
		"movie_title", movie.Title)

	return nil
}

func (s *Service) DeleteMovie(ctx context.Context, id int) error {
	err := s.repo.DeleteMovie(ctx, id)
	if err != nil {
		return err
	}

	logger := log.FromContext(ctx)
	logger.Info("movie deleted",
		"movie_id", id,
	)

	return nil
}
