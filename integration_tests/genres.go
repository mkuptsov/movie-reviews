package tests

import (
	"net/http"
	"testing"

	"github.com/cloudmachinery/movie-reviews/client"
	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/stretchr/testify/require"
)

var (
	Action *contracts.Genre
	Drama  *contracts.Genre
	Spooky *contracts.Genre
)

func genresAPIChecks(t *testing.T, c *client.Client) {
	t.Run("genres.GetGenres: empty", func(t *testing.T) {
		genres, err := c.GetGenres()
		require.NoError(t, err)
		require.Empty(t, genres)
	})

	t.Run("genres.CreateGenre: success: Action by Admin, Drama and Spooky by John Doe", func(t *testing.T) {
		cases := []struct {
			name  string
			token string
			addr  **contracts.Genre
		}{
			{"Action", adminToken, &Action},
			{"Drama", johnDoeToken, &Drama},
			{"Spooky", johnDoeToken, &Spooky},
		}

		for _, cc := range cases {
			req := &contracts.CreateGenreRequest{
				Name: cc.name,
			}
			g, err := c.CreateGenre(contracts.NewAuthenticated(req, cc.token))
			require.NoError(t, err)

			*cc.addr = g
			require.NotEmpty(t, g.ID)
			require.Equal(t, req.Name, g.Name)
		}
	})

	t.Run("genres.CreateGenre: short name", func(t *testing.T) {
		req := &contracts.CreateGenreRequest{
			Name: "Oh",
		}
		_, err := c.CreateGenre(contracts.NewAuthenticated(req, johnDoeToken))
		requireBadRequestError(t, err, "Name")
	})

	t.Run("genres.CreateGenre: existing name", func(t *testing.T) {
		req := &contracts.CreateGenreRequest{
			Name: Action.Name,
		}
		_, err := c.CreateGenre(contracts.NewAuthenticated(req, johnDoeToken))
		requireAlreadyExistsError(t, err, "genre", "name", Action.Name)
	})

	t.Run("genres.GetGenres: three genres", func(t *testing.T) {
		genres, err := c.GetGenres()

		require.NoError(t, err)
		require.Equal(t, []*contracts.Genre{Action, Drama, Spooky}, genres)
	})

	t.Run("genres.GetGenre: success", func(t *testing.T) {
		g, err := c.GetGenreByID(Spooky.ID)

		require.NoError(t, err)
		require.Equal(t, Spooky, g)
	})

	t.Run("genres.GetGenre: not found", func(t *testing.T) {
		_, err := c.GetGenreByID(fakeID)
		requireNotFoundError(t, err, "genre", "id", fakeID)
	})

	t.Run("genres.UpdateGenre: success", func(t *testing.T) {
		req := &contracts.UpdateGenreRequest{
			ID:   Spooky.ID,
			Name: "Horror",
		}
		err := c.UpdateGenre(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err)

		Spooky = getGenre(t, c, Spooky.ID)
		require.Equal(t, req.Name, Spooky.Name)
	})

	t.Run("genres.UpdateGenre: unauthorized", func(t *testing.T) {
		req := &contracts.UpdateGenreRequest{
			ID:   Spooky.ID,
			Name: "Horror",
		}
		err := c.UpdateGenre(contracts.NewAuthenticated(req, ""))
		requireUnauthorizedError(t, err, "invalid or missing token")
	})

	t.Run("genres.UpdateGenre: not found", func(t *testing.T) {
		req := &contracts.UpdateGenreRequest{
			ID:   fakeID,
			Name: "Horror",
		}
		err := c.UpdateGenre(contracts.NewAuthenticated(req, johnDoeToken))
		requireNotFoundError(t, err, "genre", "id", fakeID)
	})

	t.Run("genres.DeleteGenre: success", func(t *testing.T) {
		req := &contracts.DeleteGenreRequest{
			ID: Spooky.ID,
		}
		err := c.DeleteGenre(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err)

		Spooky = getGenre(t, c, Spooky.ID)
		require.Nil(t, Spooky)
	})

	t.Run("genres.DeleteGenre: unauthorized", func(t *testing.T) {
		req := &contracts.DeleteGenreRequest{
			ID: Drama.ID,
		}
		err := c.DeleteGenre(contracts.NewAuthenticated(req, ""))
		requireUnauthorizedError(t, err, "invalid or missing token")
	})

	t.Run("genres.GetGenres: two genres", func(t *testing.T) {
		genres, err := c.GetGenres()
		require.NoError(t, err)
		require.Equal(t, []*contracts.Genre{Action, Drama}, genres)
	})
}

func getGenre(t *testing.T, c *client.Client, id int) *contracts.Genre {
	u, err := c.GetGenreByID(id)
	if err != nil {
		cerr, ok := err.(*client.Error)
		require.True(t, ok)
		require.Equal(t, http.StatusNotFound, cerr.Code)
		return nil
	}

	return u
}
