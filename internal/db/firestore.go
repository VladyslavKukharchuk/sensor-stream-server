package db

import (
	"context"

	"cloud.google.com/go/firestore"
)

func NewFirestore(ctx context.Context, projectID string) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return client, nil
}
