package store

import (
	"context"
	"time"
)

type Lease struct {
	Name      string
	Owner     string
	ExpiresAt time.Time
}

type LeaseRepository struct{ db *DB }

func (r LeaseRepository) Acquire(ctx context.Context, name, owner string, until time.Time) (bool, error) {
	result, err := r.db.SQL.ExecContext(ctx, `INSERT INTO leases(name, owner, expires_at) VALUES(?, ?, ?) ON CONFLICT(name) DO UPDATE SET owner = excluded.owner, expires_at = excluded.expires_at WHERE leases.expires_at < ? OR leases.owner = ?`, name, owner, formatTime(until), formatTime(time.Now().UTC()), owner)
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	return count > 0, err
}

func (r LeaseRepository) Release(ctx context.Context, name, owner string) error {
	_, err := r.db.SQL.ExecContext(ctx, `DELETE FROM leases WHERE name = ? AND owner = ?`, name, owner)
	return err
}

func (r LeaseRepository) Get(ctx context.Context, name string) (Lease, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT name, owner, expires_at FROM leases WHERE name = ?`, name)
	var lease Lease
	var expires string
	if err := row.Scan(&lease.Name, &lease.Owner, &expires); err != nil {
		return lease, mapNotFound(err)
	}
	lease.ExpiresAt = parseTime(expires)
	return lease, nil
}
