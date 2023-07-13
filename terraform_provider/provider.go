package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/mkuptsov/movie-reviews/client"
	"github.com/mkuptsov/movie-reviews/contracts"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: Provider})
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"admin_email": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MOVIE_REVIEWS_ADMIN_EMAIL", nil),
			},
			"admin_password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MOVIE_REVIEWS_ADMIN_PASSWORD", nil),
			},
			"api_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MOVIE_REVIEWS_API_URL", "http://localhost:8080"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"moviereviews_user_role": movieReviewsUserRole(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"moviereviews_user": movieReviewsUser(),
		},
		ConfigureContextFunc: func(ctx context.Context, rd *schema.ResourceData) (any, diag.Diagnostics) {
			url := rd.Get("api_url").(string)
			email := rd.Get("admin_email").(string)
			password := rd.Get("admin_password").(string)
			c := client.New(url)
			req := &contracts.LoginUserRequest{
				Email:    email,
				Password: password,
			}

			res, err := c.LoginUser(req)
			if err != nil {
				return nil, diag.FromErr(fmt.Errorf("cannot login with email %s: %w", email, err))
			}
			return clientAndToken{
				client: c,
				token:  res.AccessToken,
			}, nil
		},
	}
}

type clientAndToken struct {
	client *client.Client
	token  string
}
