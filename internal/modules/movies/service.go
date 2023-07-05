package movies

import (
	"context"

	"github.com/mkuptsov/movie-reviews/internal/log"
	"github.com/mkuptsov/movie-reviews/internal/modules/genres"
	"github.com/mkuptsov/movie-reviews/internal/modules/stars"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	repo         *Repository
	genreService *genres.Service
	starsService *stars.Service
}

func NewService(repo *Repository, genresService *genres.Service, starsService *stars.Service) *Service {
	return &Service{
		repo:         repo,
		genreService: genresService,
		starsService: starsService,
	}
}

func (s *Service) CreateMovie(ctx context.Context, movie *MovieDetails) error {
	err := s.repo.CreateMovie(ctx, movie)
	if err != nil {
		return err
	}

	logger := log.FromContext(ctx)
	logger.Info("movie created",
		"movie_id", movie.ID)

	return s.assemble(ctx, movie)
}

func (s *Service) GetMovieByID(ctx context.Context, id int) (*MovieDetails, error) {
	movie, err := s.repo.GetMovieByID(ctx, id)
	if err != nil {
		return nil, err
	}
	err = s.assemble(ctx, movie)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (s *Service) GetAllPaginated(ctx context.Context, starID *int, searchTerm *string, sortByRating *string, offset, limit int) ([]*Movie, int, error) {
	return s.repo.GetAllPaginated(ctx, starID, searchTerm, sortByRating, offset, limit)
}

func (s *Service) UpdateMovie(ctx context.Context, id int, movie *MovieDetails) error {
	err := s.repo.UpdateMovie(ctx, id, movie)
	if err != nil {
		return err
	}

	logger := log.FromContext(ctx)
	logger.Info("movie updated",
		"movie_id", id)

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

func (s *Service) assemble(ctx context.Context, movie *MovieDetails) error {
	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		var err error
		movie.Genres, err = s.genreService.GetGenresByMovieID(groupCtx, movie.ID)
		return err
	})

	group.Go(func() error {
		var err error
		movie.Cast, err = s.starsService.GetCastByMovieID(groupCtx, movie.ID)
		return err
	})

	return group.Wait()
}
