package common_utils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const rowCloseErrorMsg = "pq: unexpected Parse response 'C'"
const deadLockErrorMsg = "pq: unexpected Parse response 'D'"
const badConnectionErrMsg = "driver: bad connection"
const txAbortingErrMsg = "pq: Could not complete operation in a failed transaction"

func ExecTx(ctx context.Context, pgxPool *pgxpool.Pool, fn func(tx pgx.Tx) error) error {

	tx, err := pgxPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

func ExecTxWithRetry(ctx context.Context, pgxPool *pgxpool.Pool, fn func(tx pgx.Tx) error) error {
	var retryFunc = func() error {
		tx, err := pgxPool.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			return err
		}

		err = fn(tx)
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
			}
			return err
		}

		return tx.Commit(ctx)
	}

	err := retryFunc()
	for i := 0; i < 3; i++ {
		if err == nil {
			break
		} else if strings.Contains(err.Error(), badConnectionErrMsg) ||
			strings.Contains(err.Error(), deadLockErrorMsg) ||
			strings.Contains(err.Error(), rowCloseErrorMsg) ||
			strings.Contains(err.Error(), txAbortingErrMsg) {
			// immediately RETRY??
			time.Sleep(100 * time.Millisecond)
			LogInfo(fmt.Sprintf("retry transaction %d times \n", i+1))
			err = retryFunc()
		} else {
			// DON'T NEED TO RETRY THIS ERROR
			break
		}
	}

	return err
}
