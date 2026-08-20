package service

import (
	"context"
	"sort"

	"github.com/jb843051627/bathypulse/internal/store"
)

func (s *ObservatoryService) AuditSummary(ctx context.Context, entityType, entityID string) (map[string]int, error) {
	items, err := s.EntityAudit(ctx, entityType, entityID)
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, item := range items {
		counts[item.Action]++
	}
	return counts, nil
}

func (s *ObservatoryService) AuditActions(ctx context.Context, entityType, entityID string) ([]string, error) {
	items, err := s.EntityAudit(ctx, entityType, entityID)
	if err != nil {
		return nil, err
	}
	actions := make([]string, 0, len(items))
	for _, item := range items {
		actions = append(actions, item.Action)
	}
	sort.Strings(actions)
	return actions, nil
}

func auditRecordKey(item store.AuditRecord) string {
	return item.EntityType + ":" + item.EntityID + ":" + item.Action
}
