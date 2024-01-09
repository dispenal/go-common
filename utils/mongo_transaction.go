package common_utils

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

func ExecSession(ctx context.Context, client *mongo.Client, fn func(mongo.SessionContext) error) error {
	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// Start a transaction.
	if err := session.StartTransaction(); err != nil {
		return err
	}

	// Execute the callback with the session.
	if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := fn(sc); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	// Commit the transaction.
	if err := session.CommitTransaction(ctx); err != nil {
		return err
	}

	return nil
}

func ExecSessionWithRetry(ctx context.Context, client *mongo.Client, fn func(mongo.SessionContext) error) error {
	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// Start a transaction.
	if err := session.StartTransaction(); err != nil {
		return err
	}

	// retry 3 times
	for i := 0; i < 3; i++ {
		// Execute the callback with the session.
		if err := mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
			if err := fn(sc); err != nil {
				return err
			}
			return nil
		}); err != nil {
			if err := session.AbortTransaction(ctx); err != nil {
				return err
			}
			continue
		}
		break
	}

	// Commit the transaction.
	if err := session.CommitTransaction(ctx); err != nil {
		return err
	}

	return nil
}
