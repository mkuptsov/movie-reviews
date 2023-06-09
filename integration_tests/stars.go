package tests

import (
	"testing"
	"time"

	"github.com/cloudmachinery/movie-reviews/client"
	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/stretchr/testify/require"
)

func starsAPIChecks(t *testing.T, c *client.Client) {
	var lucas, hamill, mcgregor *contracts.Star

	t.Run("stars.Create: success", func(t *testing.T) {
		cases := []struct {
			req  *contracts.CreateStarRequest
			addr **contracts.Star
		}{
			{
				req: &contracts.CreateStarRequest{
					FirstName:  "George",
					MiddleName: contracts.Ptr("Walton"),
					LastName:   "Lucas",
					BirthDate:  time.Date(1944, time.May, 14, 0, 0, 0, 0, time.UTC),
					BirthPlace: contracts.Ptr("Modesto, California6 U.S."),
					Bio:        contracts.Ptr("Famous creator of Star Wars"),
				},
				addr: &lucas,
			},
			{
				req: &contracts.CreateStarRequest{
					FirstName:  "Mark",
					MiddleName: contracts.Ptr("Richard"),
					LastName:   "Hamill",
					BirthDate:  time.Date(1951, time.September, 25, 0, 0, 0, 0, time.UTC),
					BirthPlace: contracts.Ptr("Oakland, California6 U.S."),
				},
				addr: &hamill,
			},
			{
				req: &contracts.CreateStarRequest{
					FirstName:  "Ewan",
					MiddleName: contracts.Ptr("Gordon"),
					LastName:   "McGregor",
					BirthDate:  time.Date(1971, time.March, 31, 0, 0, 0, 0, time.UTC),
					BirthPlace: contracts.Ptr("Perth, Scotland"),
				},
				addr: &mcgregor,
			},
		}

		for _, cc := range cases {

			star, err := c.CreateStar(contracts.NewAuthenticated(cc.req, johnDoeToken))
			require.NoError(t, err)

			*cc.addr = star
			require.NotEmpty(t, star.ID)
			require.NotEmpty(t, star.CreatedAt)
		}
	})

	t.Run("stars.Create: unauthorized", func(t *testing.T) {
		req := &contracts.CreateStarRequest{
			FirstName:  "George",
			MiddleName: contracts.Ptr("Walton"),
			LastName:   "Lucas",
			BirthDate:  time.Date(1944, time.May, 14, 0, 0, 0, 0, time.UTC),
			BirthPlace: contracts.Ptr("Modesto, California6 U.S."),
			Bio:        contracts.Ptr("Famous creator of Star Wars"),
		}

		_, err := c.CreateStar(contracts.NewAuthenticated(req, ""))
		requireUnauthorizedError(t, err, "invalid or missing token")
	})

	t.Run("stars.GetStarByID: success", func(t *testing.T) {
		star, err := c.GetStarByID(hamill.ID)
		require.NoError(t, err)
		require.Equal(t, star.ID, hamill.ID)
	})

	t.Run("stars.GetStarByID: not found", func(t *testing.T) {
		_, err := c.GetStarByID(fakeID)
		requireNotFoundError(t, err, "star", "id", fakeID)
	})

	t.Run("stars.GetAll: success", func(t *testing.T) {
		req := contracts.GetStarsRequest{}
		res, err := c.GetStars(&req)
		require.NoError(t, err)

		require.Equal(t, 3, res.Total)
		require.Equal(t, 1, res.Page)
		require.Equal(t, testPaginationSize, res.Size)
		require.Equal(t, []*contracts.Star{lucas, hamill}, res.Items)

		req.Page = res.Page + 1
		res, err = c.GetStars(&req)
		require.NoError(t, err)

		require.Equal(t, 3, res.Total)
		require.Equal(t, 2, req.Page)
		require.Equal(t, testPaginationSize, res.Size)
		require.Equal(t, []*contracts.Star{mcgregor}, res.Items)
	})

	t.Run("stars.Update: success", func(t *testing.T) {
		req := &contracts.UpdateStarRequest{
			ID:         mcgregor.ID,
			FirstName:  mcgregor.FirstName,
			MiddleName: mcgregor.MiddleName,
			LastName:   mcgregor.LastName,
			BirthDate:  mcgregor.BirthDate,
			BirthPlace: mcgregor.BirthPlace,
			Bio:        contracts.Ptr("Updated bio"),
		}
		err := c.UpdateStar(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err)

		res, err := c.GetStarByID(mcgregor.ID)
		require.NoError(t, err)
		require.Equal(t, *req.Bio, *res.Bio)
	})

	t.Run("stars.Update: unathorized", func(t *testing.T) {
		req := &contracts.UpdateStarRequest{
			ID:         mcgregor.ID,
			FirstName:  mcgregor.FirstName,
			MiddleName: mcgregor.MiddleName,
			LastName:   mcgregor.LastName,
			BirthDate:  mcgregor.BirthDate,
			BirthPlace: mcgregor.BirthPlace,
			Bio:        contracts.Ptr("Updated bio"),
		}
		err := c.UpdateStar(contracts.NewAuthenticated(req, ""))
		requireUnauthorizedError(t, err, "invalid or missing token")
	})

	t.Run("stars.Update: not found", func(t *testing.T) {
		req := &contracts.UpdateStarRequest{
			ID:         fakeID,
			FirstName:  mcgregor.FirstName,
			MiddleName: mcgregor.MiddleName,
			LastName:   mcgregor.LastName,
			BirthDate:  mcgregor.BirthDate,
			BirthPlace: mcgregor.BirthPlace,
			Bio:        contracts.Ptr("Updated bio"),
		}
		err := c.UpdateStar(contracts.NewAuthenticated(req, johnDoeToken))
		requireNotFoundError(t, err, "star", "id", fakeID)
	})

	t.Run("stars.Delete: unauthorized", func(t *testing.T) {
		req := &contracts.DeleteStarRequest{
			ID: mcgregor.ID,
		}
		err := c.DeleteStar(contracts.NewAuthenticated(req, ""))
		requireUnauthorizedError(t, err, "invalid or missing token")
	})

	t.Run("stars.Delete: not found", func(t *testing.T) {
		req := &contracts.DeleteStarRequest{
			ID: fakeID,
		}
		err := c.DeleteStar(contracts.NewAuthenticated(req, johnDoeToken))
		requireNotFoundError(t, err, "star", "id", fakeID)
	})

	t.Run("stars.Delete: success", func(t *testing.T) {
		req := &contracts.DeleteStarRequest{
			ID: mcgregor.ID,
		}
		err := c.DeleteStar(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err)

		_, err = c.GetStarByID(mcgregor.ID)
		requireNotFoundError(t, err, "star", "id", mcgregor.ID)
	})
}
