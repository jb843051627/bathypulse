package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/store"
)

func (s *ObservatoryService) RecordAudit(ctx context.Context, entityType, entityID, action, actor, payload string) error {
	if entityType == "" || entityID == "" || action == "" || actor == "" {
		return fmt.Errorf("audit fields: %w", model.ErrInvalid)
	}
	record := store.AuditRecord{ID: fmt.Sprintf("audit-%d", s.clock.Now().UnixNano()), EntityType: entityType, EntityID: entityID, Action: action, Actor: actor, OccurredAt: s.clock.Now(), Payload: payload}
	return s.db.Audits().Append(ctx, record)
}

func (s *ObservatoryService) EntityAudit(ctx context.Context, entityType, entityID string) ([]store.AuditRecord, error) {
	return s.db.Audits().ListEntity(ctx, entityType, entityID)
}

func auditAge(at, now time.Time) time.Duration {
	return now.Sub(at)
}
