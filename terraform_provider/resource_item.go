package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mkuptsov/movie-reviews/client"
	"github.com/mkuptsov/movie-reviews/contracts"
	"github.com/mkuptsov/movie-reviews/internal/apperrors"
)

func movieReviewsUserRole() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		CreateContext: userRoleCreateOrUpdate,
		ReadContext:   userRoleRead,
		UpdateContext: userRoleCreateOrUpdate,
		DeleteContext: userRoleDelete,
	}
}

func userRoleCreateOrUpdate(ctx context.Context, rd *schema.ResourceData, i any) diag.Diagnostics {
	cwt := i.(clientAndToken)

	accessToken := cwt.token
	c := cwt.client
	userID := rd.Get("user_id").(int)
	role := rd.Get("role").(string)
	req := contracts.NewAuthenticated(&contracts.SetUserRoleRequest{
		UserID: userID,
		Role:   role,
	}, accessToken)

	err := c.SetUserRole(req)
	if err != nil {
		diag.FromErr(fmt.Errorf("cannot set role userID=%d role=%s: %w", userID, role, err))
	}

	rd.Set("user_id", userID)
	rd.Set("user_role", role)
	rd.SetId(strconv.Itoa(userID))

	return nil
}

func userRoleRead(ctx context.Context, rd *schema.ResourceData, i any) diag.Diagnostics {
	cwt := i.(clientAndToken)

	userID := rd.Get("user_id").(int)
	user, err := cwt.client.GetUserByID(userID)
	if cerr, ok := err.(*client.Error); ok {
		if cerr.Code == int(apperrors.NotFoundCode) {
			rd.SetId("")
		}
	}
	if err != nil {
		return diag.FromErr(fmt.Errorf("cannot get user with ID %d: %w", userID, err))
	}

	rd.Set("user_id", user.ID)
	rd.Set("user_role", user.Role)
	rd.SetId(strconv.Itoa(user.ID))

	return nil
}

func userRoleDelete(ctx context.Context, rd *schema.ResourceData, i any) diag.Diagnostics {
	cwt := i.(clientAndToken)

	accessToken := cwt.token
	c := cwt.client
	userID := rd.Get("user_id").(int)
	req := contracts.NewAuthenticated(&contracts.SetUserRoleRequest{
		UserID: userID,
		Role:   "user",
	}, accessToken)

	err := c.SetUserRole(req)
	if err != nil {
		diag.FromErr(fmt.Errorf("cannot set role userID=%d role=user: %w", userID, err))
	}

	rd.Set("user_id", userID)
	rd.Set("user_role", "user")
	rd.SetId(strconv.Itoa(userID))

	return nil
}
