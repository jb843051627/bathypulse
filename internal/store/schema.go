package store

import "context"

func (d *DB) initSchema(ctx context.Context) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS stations (id TEXT PRIMARY KEY, code TEXT NOT NULL UNIQUE, basin TEXT NOT NULL, latitude REAL NOT NULL, longitude REAL NOT NULL, depth_meters INTEGER NOT NULL, status TEXT NOT NULL, last_seen TEXT NOT NULL, version INTEGER NOT NULL DEFAULT 1)`,
		`CREATE TABLE IF NOT EXISTS waveforms (id TEXT PRIMARY KEY, station_id TEXT NOT NULL, captured_at TEXT NOT NULL, duration_ms INTEGER NOT NULL, sample_rate INTEGER NOT NULL, peak REAL NOT NULL, rms REAL NOT NULL, payload BLOB NOT NULL, checksum TEXT NOT NULL, state TEXT NOT NULL, sequence INTEGER NOT NULL, FOREIGN KEY(station_id) REFERENCES stations(id))`,
		`CREATE TABLE IF NOT EXISTS samples (station_id TEXT NOT NULL, captured_at TEXT NOT NULL, sequence INTEGER NOT NULL, value REAL NOT NULL, quality INTEGER NOT NULL, PRIMARY KEY(station_id, sequence))`,
		`CREATE TABLE IF NOT EXISTS events (id TEXT PRIMARY KEY, station_id TEXT NOT NULL, kind TEXT NOT NULL, started_at TEXT NOT NULL, ended_at TEXT NOT NULL, magnitude REAL NOT NULL, confidence REAL NOT NULL, state TEXT NOT NULL, waveform_count INTEGER NOT NULL, revision INTEGER NOT NULL DEFAULT 1)`,
		`CREATE TABLE IF NOT EXISTS event_waveforms (event_id TEXT NOT NULL, waveform_id TEXT NOT NULL, PRIMARY KEY(event_id, waveform_id))`,
		`CREATE TABLE IF NOT EXISTS alerts (id TEXT PRIMARY KEY, event_id TEXT NOT NULL, station_id TEXT NOT NULL, level TEXT NOT NULL, state TEXT NOT NULL, message TEXT NOT NULL, created_at TEXT NOT NULL, acked_at TEXT, revision INTEGER NOT NULL DEFAULT 1)`,
		`CREATE TABLE IF NOT EXISTS maintenance (id TEXT PRIMARY KEY, station_id TEXT NOT NULL, window_start TEXT NOT NULL, window_end TEXT NOT NULL, state TEXT NOT NULL, reason TEXT NOT NULL, revision INTEGER NOT NULL DEFAULT 1)`,
		`CREATE TABLE IF NOT EXISTS calibration_profiles (id TEXT PRIMARY KEY, station_id TEXT NOT NULL, version INTEGER NOT NULL, points BLOB NOT NULL, enabled INTEGER NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS telemetry_packets (id TEXT PRIMARY KEY, station_id TEXT NOT NULL, type TEXT NOT NULL, received_at TEXT NOT NULL, sequence INTEGER NOT NULL, payload BLOB NOT NULL, signature TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS audit_records (id TEXT PRIMARY KEY, entity_type TEXT NOT NULL, entity_id TEXT NOT NULL, action TEXT NOT NULL, actor TEXT NOT NULL, occurred_at TEXT NOT NULL, payload TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS station_snapshots (station_id TEXT NOT NULL, observed_at TEXT NOT NULL, payload BLOB NOT NULL, PRIMARY KEY(station_id, observed_at))`,
		`CREATE TABLE IF NOT EXISTS retention_rules (id TEXT PRIMARY KEY, station_id TEXT NOT NULL, keep_days INTEGER NOT NULL, enabled INTEGER NOT NULL, updated_at TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS station_sequences (station_id TEXT PRIMARY KEY, next_sequence INTEGER NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS leases (name TEXT PRIMARY KEY, owner TEXT NOT NULL, expires_at TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS analysis_runs (id TEXT PRIMARY KEY, station_id TEXT NOT NULL, started_at TEXT NOT NULL, finished_at TEXT NOT NULL, state TEXT NOT NULL, samples INTEGER NOT NULL, events INTEGER NOT NULL, score REAL NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS dispatch_jobs (id TEXT PRIMARY KEY, station_id TEXT NOT NULL, kind TEXT NOT NULL, state TEXT NOT NULL, attempt INTEGER NOT NULL, created_at TEXT NOT NULL, updated_at TEXT NOT NULL, payload TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS incidents (id TEXT PRIMARY KEY, station_id TEXT NOT NULL, alert_id TEXT NOT NULL, state TEXT NOT NULL, summary TEXT NOT NULL, opened_at TEXT NOT NULL, closed_at TEXT, revision INTEGER NOT NULL)`,
		`CREATE INDEX IF NOT EXISTS idx_waveforms_station_time ON waveforms(station_id, captured_at)`,
		`CREATE INDEX IF NOT EXISTS idx_events_station_time ON events(station_id, started_at)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_state ON alerts(state, created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_samples_station_time ON samples(station_id, captured_at)`,
		`CREATE INDEX IF NOT EXISTS idx_packets_station_seq ON telemetry_packets(station_id, sequence)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_entity_time ON audit_records(entity_type, entity_id, occurred_at)`,
	}
	for _, statement := range statements {
		if _, err := d.SQL.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}
