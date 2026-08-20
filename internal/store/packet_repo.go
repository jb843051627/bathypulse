package store

import (
	"context"
	"database/sql"

	"github.com/jb843051627/bathypulse/internal/model"
)

type PacketRepository struct{ db *DB }

func (r PacketRepository) Save(ctx context.Context, packet model.TelemetryPacket) error {
	_, err := r.db.SQL.ExecContext(ctx, `INSERT INTO telemetry_packets(id, station_id, type, received_at, sequence, payload, signature) VALUES(?, ?, ?, ?, ?, ?, ?)`, packet.ID, packet.StationID, packet.Type, formatTime(packet.ReceivedAt), packet.Sequence, packet.Payload, packet.Signature)
	return err
}

func (r PacketRepository) Get(ctx context.Context, id string) (*model.TelemetryPacket, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT id, station_id, type, received_at, sequence, payload, signature FROM telemetry_packets WHERE id = ?`, id)
	var packet model.TelemetryPacket
	var kind, received string
	if err := row.Scan(&packet.ID, &packet.StationID, &kind, &received, &packet.Sequence, &packet.Payload, &packet.Signature); err != nil {
		return nil, mapNotFound(err)
	}
	packet.Type = model.PacketType(kind)
	packet.ReceivedAt = parseTime(received)
	packet.Payload = model.CloneBytes(packet.Payload)
	return &packet, nil
}

func (r PacketRepository) List(ctx context.Context, stationID string, limit int) ([]model.TelemetryPacket, error) {
	rows, err := r.db.SQL.QueryContext(ctx, `SELECT id, station_id, type, received_at, sequence, payload, signature FROM telemetry_packets WHERE station_id = ? ORDER BY sequence DESC LIMIT ?`, stationID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]model.TelemetryPacket, 0)
	for rows.Next() {
		var packet model.TelemetryPacket
		var kind, received string
		if err := rows.Scan(&packet.ID, &packet.StationID, &kind, &received, &packet.Sequence, &packet.Payload, &packet.Signature); err != nil {
			return nil, err
		}
		packet.Type = model.PacketType(kind)
		packet.ReceivedAt = parseTime(received)
		packet.Payload = model.CloneBytes(packet.Payload)
		result = append(result, packet)
	}
	return result, rows.Err()
}

func (r PacketRepository) Count(ctx context.Context, stationID string) (int, error) {
	row := r.db.SQL.QueryRowContext(ctx, `SELECT COUNT(*) FROM telemetry_packets WHERE station_id = ?`, stationID)
	var count int
	return count, row.Scan(&count)
}

var _ = sql.ErrNoRows
