package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/store"
)

func (s *ObservatoryService) SetRetention(ctx context.Context, rule store.RetentionRule) error {
	if rule.StationID == "" || rule.KeepDays < 1 || rule.KeepDays > 3650 {
		return fmt.Errorf("retention rule: %w", model.ErrInvalid)
	}
	if _, err := s.stations.Get(ctx, rule.StationID); err != nil {
		return err
	}
	retention := s.db.Retention()
	if err := retention.Disable(ctx, rule.StationID); err != nil {
		return err
	}
	rule.Enabled = true
	rule.UpdatedAt = s.clock.Now()
	return retention.Save(ctx, rule)
}

func (s *ObservatoryService) RunRetention(ctx context.Context, stationID string) (int64, error) {
	rule, err := s.db.Retention().Get(ctx, stationID)
	if err != nil {
		return 0, err
	}
	cutoff := s.clock.Now().Add(-time.Duration(rule.KeepDays) * 24 * time.Hour)
	count, err := s.db.Retention().RemoveWaveforms(ctx, stationID, cutoff)
	if err != nil {
		return 0, err
	}
	s.cacheMu.Lock()
	delete(s.waveCache, stationID)
	s.cacheMu.Unlock()
	return count, nil
}
