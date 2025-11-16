package db

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

func NewFirestore(ctx context.Context, projectID, firestoreDatabaseID string) (*firestore.Client, error) {
	client, err := firestore.NewClientWithDatabase(ctx, projectID, firestoreDatabaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to creates a new firestore client: %w", err)
	}

	return client, nil
}
