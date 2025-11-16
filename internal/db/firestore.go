package db

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

func NewFirestore(ctx context.Context, projectID string) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to creates a new firestore client: %w", err)
	}

	return client, nil
}
