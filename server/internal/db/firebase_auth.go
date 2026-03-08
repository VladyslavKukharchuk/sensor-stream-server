package db

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

func NewFirebaseAuth(ctx context.Context, projectID string) (*auth.Client, error) {
	config := &firebase.Config{ProjectID: projectID}

	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %w", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting firebase auth client: %w", err)
	}

	return client, nil
}
