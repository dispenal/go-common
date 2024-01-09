package common_utils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func ExecSession(ctx context.Context, mongoClient *mongo.Client, fn func(session mongo.Session) error) error {
	session, err := mongoClient.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	err = session.StartTransaction()
	if err != nil {
		return err
	}

	err = fn(session)
	if err != nil {
		if rbErr := session.AbortTransaction(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return session.CommitTransaction(ctx)
}

func ExecSessionWithRetry(ctx context.Context, mongoClient *mongo.Client, fn func(session mongo.Session) error) error {
	var retryFunc = func() error {
		session, err := mongoClient.StartSession()
		if err != nil {
			return err
		}
		defer session.EndSession(ctx)

		err = session.StartTransaction()
		if err != nil {
			return err
		}

		err = fn(session)
		if err != nil {
			if rbErr := session.AbortTransaction(ctx); rbErr != nil {
				return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
			}
			return err
		}

		return session.CommitTransaction(ctx)
	}

	err := retryFunc()
	for i := 0; i < 3; i++ {
		if err == nil {
			break
		} else if strings.Contains(err.Error(), "not primary") ||
			strings.Contains(err.Error(), "no reachable servers") ||
			strings.Contains(err.Error(), "connection() error") ||
			strings.Contains(err.Error(), "connection reset by peer") ||
			strings.Contains(err.Error(), "connection refused") ||
			strings.Contains(err.Error(), "connection timed out") ||
			strings.Contains(err.Error(), "connection closed") ||
			strings.Contains(err.Error(), "connection reset") {
			time.Sleep(100 * time.Millisecond)
			err = retryFunc()
		} else {
			break
		}
	}

	return err
}
