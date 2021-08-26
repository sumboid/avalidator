package main

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
)

func GetUserByEmail(ctx context.Context, c *graphql.Client, email string) (*UserModel, error) {
	req := graphql.NewRequest(`
		query ($email: String!) {
			users(where: {email: {_eq: $email}}, limit: 1) {
				id
				role
			}
		}
	`)

	req.Var("email", email)

	req.Header.Set("x-hasura-admin-secret", Config.GraphQL.Secret)

	var resp struct {
		Users []*UserModel `json:"users"`
	}

	if err := c.Run(ctx, req, &resp); err != nil {
		return nil, err
	}

	if resp.Users == nil {
		return nil, fmt.Errorf("Unexpected empty output")
	}

	if len(resp.Users) == 0 {
		return nil, CreateNotFoundError(email)
	}

	return resp.Users[0], nil
}

func InsertUser(ctx context.Context, c *graphql.Client, model *UserModel) (*UserModel, error) {
	req := graphql.NewRequest(`
		mutation($email: String!, $name: String!, $role: String!) {
			__typename
			insert_users(objects: {email: $email, name: $name, role: $role}) {
				returning {
					id
					role
					name
					email
					created_at
					updated_at
				}
			}
		}
	`)

	req.Var("email", model.Email)
	req.Var("name", model.Name)
	req.Var("role", model.Role)

	req.Header.Set("x-hasura-admin-secret", Config.GraphQL.Secret)

	var resp struct {
		InsertUsers struct {
			Returning []*UserModel `json:"returning"`
		} `json:"insert_users"`
	}

	if err := c.Run(ctx, req, &resp); err != nil {
		return nil, err
	}

	if len(resp.InsertUsers.Returning) == 0 {
		return nil, fmt.Errorf("Unexpected empty output")
	}

	return resp.InsertUsers.Returning[0], nil
}
