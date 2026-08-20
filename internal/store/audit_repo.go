package store

import (
	"context"
	"time"
)

type AuditRecord struct {
	ID         string
	EntityType string
	EntityID   string
	Action     string
	Actor      string
	OccurredAt time.Time
	Payload    string
}

type AuditRepository struct{ db *DB }

func (r AuditRepository) Append(ctx context.Context, record AuditRecord) error {
	_, err := r.db.SQL.ExecContext(ctx, `INSERT INTO audit_records(id, entity_type, entity_id, action, actor, occurred_at, payload) VALUES(?, ?, ?, ?, ?, ?, ?)`, record.ID, record.EntityType, record.EntityID, record.Action, record.Actor, formatTime(record.OccurredAt), record.Payload)
	return err
}

func (r AuditRepository) ListEntity(ctx context.Context, entityType, entityID string) ([]AuditRecord, error) {
	rows, err := r.db.SQL.QueryContext(ctx, `SELECT id, entity_type, entity_id, action, actor, occurred_at, payload FROM audit_records WHERE entity_type = ? AND entity_id = ? ORDER BY occurred_at`, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]AuditRecord, 0)
	for rows.Next() {
		var item AuditRecord
		var occurred string
		if err := rows.Scan(&item.ID, &item.EntityType, &item.EntityID, &item.Action, &item.Actor, &occurred, &item.Payload); err != nil {
			return nil, err
		}
		item.OccurredAt = parseTime(occurred)
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r AuditRepository) Count(ctx context.Context) (int, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT COUNT(*) FROM audit_records`)
	var count int
	return count, row.Scan(&count)
}
