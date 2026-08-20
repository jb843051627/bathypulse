package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jb843051627/bathypulse/internal/model"
	_ "modernc.org/sqlite"
)

type DB struct {
	SQL  *sql.DB
	Path string
}

func Open(path string) (*DB, error) {
	if path == "" {
		return nil, fmt.Errorf("database path is empty")
	}
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	sqlDB, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(8)
	sqlDB.SetMaxIdleConns(8)
	db := &DB{SQL: sqlDB, Path: path}
	if err := db.configure(context.Background()); err != nil {
		sqlDB.Close()
		return nil, err
	}
	if err := db.initSchema(context.Background()); err != nil {
		sqlDB.Close()
		return nil, err
	}
	return db, nil
}

func (d *DB) configure(ctx context.Context) error {
	_, err := d.SQL.ExecContext(ctx, "PRAGMA busy_timeout = 5000")
	return err
}

func (d *DB) Close() error {
	return d.SQL.Close()
}

func (d *DB) Ping(ctx context.Context) error {
	return d.SQL.PingContext(ctx)
}

func formatTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339Nano)
}

func parseTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

func mapNotFound(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return model.ErrNotFound
	}
	return err
}

func stationRepo(d *DB) StationRepository         { return StationRepository{db: d} }
func waveformRepo(d *DB) WaveformRepository       { return WaveformRepository{db: d} }
func eventRepo(d *DB) EventRepository             { return EventRepository{db: d} }
func alertRepo(d *DB) AlertRepository             { return AlertRepository{db: d} }
func maintenanceRepo(d *DB) MaintenanceRepository { return MaintenanceRepository{db: d} }
func sampleRepo(d *DB) SampleRepository           { return SampleRepository{db: d} }
func queryRepo(d *DB) QueryRepository             { return QueryRepository{db: d} }
func calibrationRepo(d *DB) CalibrationRepository { return CalibrationRepository{db: d} }
func packetRepo(d *DB) PacketRepository           { return PacketRepository{db: d} }
func auditRepo(d *DB) AuditRepository             { return AuditRepository{db: d} }
func snapshotRepo(d *DB) SnapshotRepository       { return SnapshotRepository{db: d} }
func retentionRepo(d *DB) RetentionRepository     { return RetentionRepository{db: d} }
func sequenceRepo(d *DB) SequenceRepository       { return SequenceRepository{db: d} }
func leaseRepo(d *DB) LeaseRepository             { return LeaseRepository{db: d} }

func (d *DB) Stations() StationRepository         { return stationRepo(d) }
func (d *DB) Waveforms() WaveformRepository       { return waveformRepo(d) }
func (d *DB) Events() EventRepository             { return eventRepo(d) }
func (d *DB) Alerts() AlertRepository             { return alertRepo(d) }
func (d *DB) Maintenance() MaintenanceRepository  { return maintenanceRepo(d) }
func (d *DB) Samples() SampleRepository           { return sampleRepo(d) }
func (d *DB) Queries() QueryRepository            { return queryRepo(d) }
func (d *DB) Calibrations() CalibrationRepository { return calibrationRepo(d) }
func (d *DB) Packets() PacketRepository           { return packetRepo(d) }
func (d *DB) Audits() AuditRepository             { return auditRepo(d) }
func (d *DB) Snapshots() SnapshotRepository       { return snapshotRepo(d) }
func (d *DB) Retention() RetentionRepository      { return retentionRepo(d) }
func (d *DB) Sequences() SequenceRepository       { return sequenceRepo(d) }
func (d *DB) Leases() LeaseRepository             { return leaseRepo(d) }
