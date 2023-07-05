package client

import "github.com/mkuptsov/movie-reviews/contracts"

func (c *Client) GetUserByUserName(userName string) (*contracts.User, error) {
	var u contracts.User

	_, err := c.client.R().
		SetResult(&u).
		Get(c.path("/api/users/username/%s", userName))

	return &u, err
}

func (c *Client) GetUserByID(id int) (*contracts.User, error) {
	var u contracts.User

	_, err := c.client.R().
		SetResult(&u).
		Get(c.path("/api/users/%d", id))

	return &u, err
}

func (c *Client) UpdateUser(req *contracts.AuthenticatedRequest[*contracts.UpdateUserRequest]) error {
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetBody(req.Request).
		Put(c.path("/api/users/%d", req.Request.UserID))

	return err
}

func (c *Client) DeleteUser(req *contracts.AuthenticatedRequest[*contracts.DeleteUserRequest]) error {
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		Delete(c.path("/api/users/%d", req.Request.UserID))

	return err
}

func (c *Client) SetUserRole(req *contracts.AuthenticatedRequest[*contracts.SetUserRoleRequest]) error {
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		Put(c.path("/api/users/%d/role/%s", req.Request.UserID, req.Request.Role))

	return err
}
