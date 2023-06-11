package client

import "github.com/cloudmachinery/movie-reviews/contracts"

func (c *Client) CreateStar(req *contracts.AuthenticatedRequest[*contracts.CreateStarRequest]) (*contracts.StarDetails, error) {
	var star contracts.StarDetails

	_, err := c.client.R().
		SetResult(&star).
		SetAuthToken(req.AccessToken).
		SetBody(req.Request).
		Post(c.path("/api/stars"))

	return &star, err
}

func (c *Client) GetStarByID(id int) (*contracts.StarDetails, error) {
	var star contracts.StarDetails

	_, err := c.client.R().
		SetResult(&star).
		Get(c.path("/api/stars/%d", id))

	return &star, err
}

func (c *Client) GetStars(req *contracts.GetStarsRequest) (*contracts.PaginatedResponse[contracts.Star], error) {
	var res contracts.PaginatedResponse[contracts.Star]

	_, err := c.client.R().
		SetResult(&res).
		SetQueryParams(req.PaginatiedRequest.ToQueryParams()).
		Get(c.path("/api/stars"))

	return &res, err
}

func (c *Client) UpdateStar(req *contracts.AuthenticatedRequest[*contracts.UpdateStarRequest]) error {
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetBody(req.Request).
		Put(c.path("/api/stars/%d", req.Request.ID))

	return err
}

func (c *Client) DeleteStar(req *contracts.AuthenticatedRequest[*contracts.DeleteStarRequest]) error {
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		Delete(c.path("/api/stars/%d", req.Request.ID))

	return err
}
