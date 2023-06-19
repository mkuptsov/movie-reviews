package tests

import (
	"testing"
	"time"

	"github.com/cloudmachinery/movie-reviews/client"
	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/stretchr/testify/require"
)

func moviesAPIChecks(t *testing.T, c *client.Client) {
	var starWars, kingsMan, trainspotting *contracts.MovieDetails

	t.Run("movies.Create: success", func(t *testing.T) {
		cases := []struct {
			req  *contracts.CreateMovieRequest
			addr **contracts.MovieDetails
		}{
			{
				req: &contracts.CreateMovieRequest{
					Title:       "Star Wars",
					ReleaseDate: time.Date(1977, time.May, 25, 0, 0, 0, 0, time.UTC),
					Description: `Amid a galactic civil war, Rebel Alliance spies have stolen plans to the Galactic
					 Empire's Death Star, a massive space station capable of destroying entire planets. Imperial 
					 Senator Princess Leia Organa of Alderaan, secretly one of the Rebellion's leaders, has obtained 
					 its schematics, but her ship is intercepted by an Imperial Star Destroyer under the command of 
					 the ruthless Empire agent Darth Vader. Before she is captured, Leia hides the plans in the memory 
					 system of astromech droid R2-D2, who flees in an escape pod to the nearby desert planet Tatooine 
					 alongside his companion, protocol droid C-3PO.`,
					Genres: []int{Action.ID, Drama.ID},
					Cast: []*contracts.MovieCreditInfo{
						{
							StarID: lucas.ID,
							Role:   "director",
						},
						{
							StarID:  hamill.ID,
							Role:    "actor",
							Details: contracts.Ptr("char1, char2"),
						},
					},
				},
				addr: &starWars,
			},
			{
				req: &contracts.CreateMovieRequest{
					Title:       "Kingsman: The Secret Service",
					ReleaseDate: time.Date(2014, time.December, 13, 0, 0, 0, 0, time.UTC),
					Description: `In 1997, probationary secret agent Lee Unwin sacrifices himself in the Middle 
					East to save his superior, Harry Hart. Blaming himself for Lee's death, Harry returns to London 
					and gives Lee's young son Gary "Eggsy" a medal engraved with an emergency assistance number.`,
					Genres: []int{Action.ID},
					Cast: []*contracts.MovieCreditInfo{
						{
							StarID:  hamill.ID,
							Role:    "actor",
							Details: contracts.Ptr("char3"),
						},
					},
				},
				addr: &kingsMan,
			},
			{
				req: &contracts.CreateMovieRequest{
					Title:       "Trainspotting",
					ReleaseDate: time.Date(1996, time.February, 23, 0, 0, 0, 0, time.UTC),
					Description: `In Scotland, Mark Renton, a 26-year-old unemployed heroin addict, lives with his 
					parents in the Edinburgh ward of Leith and regularly takes drugs with his "friends": treacherous, 
					womanizing James Bond fanatic Simon "Sick Boy" Williamson; docile and bumbling Daniel "Spud" 
					Murphy and Swanney, "Mother Superior", their dealer. Renton's other friends, aggressive, alcoholic 
					psychopath Francis "Franco" Begbie and honest footballer and recreational speed user Tommy Mackenzie, 
					who both abstain from heroin, warn him about his dangerous drug habit. `,
					Genres: []int{Drama.ID},
					Cast: []*contracts.MovieCreditInfo{
						{
							StarID:  hamill.ID,
							Role:    "actor",
							Details: contracts.Ptr("char4"),
						},
					},
				},
				addr: &trainspotting,
			},
		}

		for _, cc := range cases {

			movie, err := c.CreateMovie(contracts.NewAuthenticated(cc.req, johnDoeToken))
			require.NoError(t, err)

			*cc.addr = movie
			require.NotEmpty(t, movie.ID)
			require.NotEmpty(t, movie.CreatedAt)
			require.NotEmpty(t, movie.Genres)
			require.Equal(t, len(cc.req.Genres), len(movie.Genres))
			require.NotEmpty(t, movie.Cast)
			require.Equal(t, len(cc.req.Cast), len(movie.Cast))
		}
	})

	t.Run("movies.Create: unauthorized", func(t *testing.T) {
		req := &contracts.CreateMovieRequest{
			Title:       starWars.Title,
			ReleaseDate: time.Date(1977, time.May, 25, 0, 0, 0, 0, time.UTC),
		}

		_, err := c.CreateMovie(contracts.NewAuthenticated(req, ""))
		requireUnauthorizedError(t, err, "invalid or missing token")
	})

	t.Run("movies.GetMovieByID: success", func(t *testing.T) {
		movie, err := c.GetMovieByID(starWars.ID)
		require.NoError(t, err)
		require.Equal(t, starWars.ID, movie.ID)
		require.Equal(t, len(starWars.Genres), len(movie.Genres))
		for i, genre := range starWars.Genres {
			require.Equal(t, *genre, *movie.Genres[i])
		}
		require.Equal(t, len(starWars.Cast), len(movie.Cast))
		for i, cast := range starWars.Cast {
			require.Equal(t, *cast, *movie.Cast[i])
		}
	})

	t.Run("movies.GetMovieByID: not found", func(t *testing.T) {
		_, err := c.GetMovieByID(fakeID)
		requireNotFoundError(t, err, "movie", "id", fakeID)
	})

	t.Run("movies.GetAll: success", func(t *testing.T) {
		req := contracts.GetMoviesRequest{}
		res, err := c.GetMovies(&req)
		require.NoError(t, err)

		require.Equal(t, 3, res.Total)
		require.Equal(t, 1, res.Page)
		require.Equal(t, testPaginationSize, res.Size)
		require.Equal(t, []*contracts.Movie{&starWars.Movie, &kingsMan.Movie}, res.Items)

		req.Page = res.Page + 1
		res, err = c.GetMovies(&req)
		require.NoError(t, err)

		require.Equal(t, 3, res.Total)
		require.Equal(t, 2, req.Page)
		require.Equal(t, testPaginationSize, res.Size)
		require.Equal(t, []*contracts.Movie{&trainspotting.Movie}, res.Items)
	})

	t.Run("stars.GetAll: by movie ID success", func(t *testing.T) {
		req := contracts.GetStarsRequest{
			MovieID: contracts.Ptr(kingsMan.ID),
		}
		res, err := c.GetStars(&req)
		require.NoError(t, err)
		require.Equal(t, len(kingsMan.Cast), res.Total)
		require.Equal(t, 1, res.Page)
		require.Equal(t, testPaginationSize, res.Size)
		require.Equal(t, []*contracts.Star{&hamill.Star}, res.Items)
	})

	t.Run("movies.Update: the same genre success", func(t *testing.T) {
		req := &contracts.UpdateMovieRequest{
			ID:          trainspotting.ID,
			Title:       trainspotting.Title,
			ReleaseDate: trainspotting.ReleaseDate,
			Description: "updated description",
			Version:     0,
			Genres:      []int{Drama.ID},
			Cast: []*contracts.MovieCreditInfo{
				{
					StarID: lucas.ID,
					Role:   "producer",
				},
				{
					StarID: hamill.ID,
					Role:   "director",
				},
			},
		}
		err := c.UpdateMovie(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err)

		res, err := c.GetMovieByID(trainspotting.ID)
		require.NoError(t, err)

		require.Equal(t, req.Description, res.Description)
		require.Equal(t, len(req.Genres), len(res.Genres))
		for i, genreID := range req.Genres {
			require.Equal(t, genreID, res.Genres[i].ID)
		}
		require.Equal(t, len(req.Cast), len(res.Cast))
		for i, mc := range req.Cast {
			require.Equal(t, mc.StarID, res.Cast[i].Star.ID)
			require.Equal(t, mc.Role, res.Cast[i].Role)
		}
	})

	t.Run("movies.Update: different genres success", func(t *testing.T) {
		req := &contracts.UpdateMovieRequest{
			ID:          trainspotting.ID,
			Title:       trainspotting.Title,
			ReleaseDate: trainspotting.ReleaseDate,
			Description: "updated description",
			Version:     1,
			Genres:      []int{Action.ID},
			Cast: []*contracts.MovieCreditInfo{
				{
					StarID: hamill.ID,
					Role:   "director",
				},
			},
		}
		err := c.UpdateMovie(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err)

		res, err := c.GetMovieByID(trainspotting.ID)

		require.NoError(t, err)
		require.Equal(t, req.Description, res.Description)
		require.Equal(t, len(req.Genres), len(res.Genres))
		for i, genreID := range req.Genres {
			require.Equal(t, genreID, res.Genres[i].ID)
		}
		require.Equal(t, len(req.Cast), len(res.Cast))
		for i, mc := range req.Cast {
			require.Equal(t, mc.StarID, res.Cast[i].Star.ID)
			require.Equal(t, mc.Role, res.Cast[i].Role)
		}
	})

	t.Run("movies.Update: unathorized", func(t *testing.T) {
		req := &contracts.UpdateMovieRequest{
			ID:          trainspotting.ID,
			Title:       trainspotting.Title,
			ReleaseDate: trainspotting.ReleaseDate,
			Description: "updated description",
			Version:     0,
		}
		err := c.UpdateMovie(contracts.NewAuthenticated(req, ""))
		requireUnauthorizedError(t, err, "invalid or missing token")
	})

	t.Run("movies.Update: not found", func(t *testing.T) {
		req := &contracts.UpdateMovieRequest{
			ID:          fakeID,
			Title:       trainspotting.Title,
			ReleaseDate: trainspotting.ReleaseDate,
			Description: "updated description",
			Version:     0,
		}
		err := c.UpdateMovie(contracts.NewAuthenticated(req, johnDoeToken))
		requireNotFoundError(t, err, "movie", "id", fakeID)
	})

	t.Run("movies.Update: invalid version", func(t *testing.T) {
		req := &contracts.UpdateMovieRequest{
			ID:          trainspotting.ID,
			Title:       trainspotting.Title,
			ReleaseDate: trainspotting.ReleaseDate,
			Description: "updated description",
			Version:     0,
		}
		err := c.UpdateMovie(contracts.NewAuthenticated(req, johnDoeToken))
		requireVersionMismatch(t, err, "movie", "id", req.ID, req.Version)
	})

	t.Run("movies.Delete: unauthorized", func(t *testing.T) {
		req := &contracts.DeleteMovieRequest{
			ID: trainspotting.ID,
		}
		err := c.DeleteMovie(contracts.NewAuthenticated(req, ""))
		requireUnauthorizedError(t, err, "invalid or missing token")
	})

	t.Run("movies.Delete: not found", func(t *testing.T) {
		req := &contracts.DeleteMovieRequest{
			ID: fakeID,
		}
		err := c.DeleteMovie(contracts.NewAuthenticated(req, johnDoeToken))
		requireNotFoundError(t, err, "movie", "id", fakeID)
	})

	t.Run("movies.Delete: success", func(t *testing.T) {
		req := &contracts.DeleteMovieRequest{
			ID: trainspotting.ID,
		}
		err := c.DeleteMovie(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err)

		_, err = c.GetMovieByID(trainspotting.ID)
		requireNotFoundError(t, err, "movie", "id", trainspotting.ID)
	})
}
