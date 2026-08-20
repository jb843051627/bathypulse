package store

import (
	"context"
	"database/sql"
)

func (d *DB) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := d.SQL.Begin()
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func execTx(ctx context.Context, tx *sql.Tx, query string, args ...any) error {
	_, err := tx.ExecContext(ctx, query, args...)
	return err
}
