package service

import (
	"context"
	"fmt"

	"github.com/jb843051627/bathypulse/internal/model"
	"github.com/jb843051627/bathypulse/internal/store"
	"github.com/jb843051627/bathypulse/internal/validation"
)

func (s *ObservatoryService) AcceptPacket(ctx context.Context, packet model.TelemetryPacket) (*model.TelemetryPacket, error) {
	if err := validation.Packet(packet); err != nil {
		return nil, err
	}
	if _, err := s.stations.Get(ctx, packet.StationID); err != nil {
		return nil, fmt.Errorf("packet station: %w", err)
	}
	if err := storePackets(s.db).Save(ctx, packet); err != nil {
		return nil, fmt.Errorf("store packet: %w", err)
	}
	return &packet, nil
}

func (s *ObservatoryService) GetPacket(ctx context.Context, id string) (*model.TelemetryPacket, error) {
	return storePackets(s.db).Get(ctx, id)
}

func (s *ObservatoryService) ListPackets(ctx context.Context, stationID string, limit int) ([]model.TelemetryPacket, error) {
	if limit <= 0 {
		limit = 100
	}
	return storePackets(s.db).List(ctx, stationID, limit)
}

func storePackets(db *store.DB) store.PacketRepository {
	return db.Packets()
}
