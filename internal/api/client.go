package api

import (
	"context"
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

// Client provides methods for interacting with GitHub API
type Client struct {
	graphql *api.GraphQLClient
}

// NewClient creates a new API client using gh authentication
func NewClient() (*Client, error) {
	gql, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("create GraphQL client: %w", err)
	}

	return &Client{graphql: gql}, nil
}

// NewClientWithOptions creates a client with custom options (for testing)
func NewClientWithOptions(opts api.ClientOptions) (*Client, error) {
	gql, err := api.NewGraphQLClient(opts)
	if err != nil {
		return nil, fmt.Errorf("create GraphQL client: %w", err)
	}

	return &Client{graphql: gql}, nil
}

// graphQLString is a helper to create graphql.String values
func graphQLString(s string) graphql.String {
	return graphql.String(s)
}

// graphQLInt is a helper to create graphql.Int values
func graphQLInt(i int) graphql.Int {
	return graphql.Int(i)
}

// graphQLID is a helper to create graphql.ID values
func graphQLID(s string) graphql.ID {
	return graphql.ID(s)
}

// queryWithContext executes a GraphQL query with context
func (c *Client) queryWithContext(ctx context.Context, name string, query interface{}, variables map[string]interface{}) error {
	err := c.graphql.QueryWithContext(ctx, name, query, variables)
	if err != nil {
		return handleError(err)
	}
	return nil
}

// mutateWithContext executes a GraphQL mutation with context
func (c *Client) mutateWithContext(ctx context.Context, name string, mutation interface{}, variables map[string]interface{}) error {
	err := c.graphql.MutateWithContext(ctx, name, mutation, variables)
	if err != nil {
		return handleError(err)
	}
	return nil
}
