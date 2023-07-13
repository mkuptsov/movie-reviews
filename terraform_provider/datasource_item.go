package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mkuptsov/movie-reviews/client"
	"github.com/mkuptsov/movie-reviews/internal/apperrors"
)

func movieReviewsUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		ReadContext: userRead,
	}
}

func userRead(ctx context.Context, rd *schema.ResourceData, i any) diag.Diagnostics {
	cwt := i.(clientAndToken)

	username := rd.Get("username").(string)
	user, err := cwt.client.GetUserByUserName(username)
	if cerr, ok := err.(*client.Error); ok {
		if cerr.Code == int(apperrors.NotFoundCode) {
			rd.Set("role", "")
		}
	}
	if err != nil {
		return diag.FromErr(fmt.Errorf("cannot get user with username %s: %w", username, err))
	}

	rd.Set("role", user.Role)
	rd.SetId(strconv.Itoa(user.ID))

	return nil
}
