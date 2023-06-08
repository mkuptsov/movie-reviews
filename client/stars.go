package client

import "github.com/cloudmachinery/movie-reviews/contracts"

func (c *Client) CreateStar(req *contracts.AuthenticatedRequest[*contracts.CreateStarRequest]) (*contracts.Star, error) {
	var star contracts.Star

	_, err := c.client.R().
		SetResult(&star).
		SetAuthToken(req.AccessToken).
		SetBody(req.Request).
		Post(c.path("/api/stars"))

	return &star, err
}

func (c *Client) GetStarByID(id int) (*contracts.Star, error) {
	var star contracts.Star

	_, err := c.client.R().
		SetResult(&star).
		Get(c.path("/api/stars/%d", id))

	return &star, err
}
