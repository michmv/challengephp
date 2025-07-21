package db

import (
	"challengephp/lib"
	"challengephp/src/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	TableEvents     = "events"
	TableUsers      = "users"
	TableEventTypes = "event_types"

	PrepareSql = []string{
		"ALTER TABLE " + TableEvents + " DISABLE TRIGGER ALL;",
		"TRUNCATE TABLE " + TableEvents + " RESTART IDENTITY;",
		"TRUNCATE TABLE " + TableEventTypes + " RESTART IDENTITY;",
		"TRUNCATE TABLE " + TableUsers + " RESTART IDENTITY;",
	}

	RecoveryDbAfterPrepare = []string{
		"ALTER TABLE " + TableEvents + " ENABLE TRIGGER ALL;",
		"SELECT setval('users_id_seq', (SELECT MAX(id) FROM " + TableUsers + "));",
		"SELECT setval('event_types_id_seq', (SELECT MAX(id) FROM " + TableEventTypes + "));",
		"SELECT setval('events_id_seq', (SELECT MAX(id) FROM " + TableEvents + "));",
	}

	CreateIndex = []string{
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_events_count ON " + TableEvents + " USING btree (id)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_events_user_timestamp ON " + TableEvents + " (user_id, timestamp DESC)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_events_timestamp_desc ON " + TableEvents + " (timestamp DESC)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_events_type_timestamp ON " + TableEvents + " (type_id, timestamp DESC)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_events_stats ON " + TableEvents + " (user_id, ((metadata->>'page')), type_id)",
		"CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_events_covering ON " + TableEvents + " (user_id, type_id, timestamp DESC) INCLUDE (id, metadata)",
	}

	DropIndex = []string{
		"ALTER TABLE " + TableEvents + " DISABLE TRIGGER ALL",
		"DROP INDEX CONCURRENTLY IF EXISTS idx_events_count",
		"DROP INDEX CONCURRENTLY IF EXISTS idx_events_user_timestamp",
		"DROP INDEX CONCURRENTLY IF EXISTS idx_events_timestamp_desc",
		"DROP INDEX CONCURRENTLY IF EXISTS idx_events_type_timestamp",
		"DROP INDEX CONCURRENTLY IF EXISTS idx_events_stats",
		"DROP INDEX CONCURRENTLY IF EXISTS idx_events_covering",
	}

	Init = []string{
		`DROP TABLE IF EXISTS ` + TableUsers + `;`,
		`CREATE TABLE ` + TableUsers + ` (
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR NOT NULL,
			created_at TIMESTAMP NOT NULL
		);`,
		`DROP TABLE IF EXISTS ` + TableEventTypes + `;`,
		`CREATE TABLE ` + TableEventTypes + ` (
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR NOT NULL UNIQUE
		);`,
		`DROP TABLE IF EXISTS ` + TableEvents + `;`,
		`CREATE TABLE ` + TableEvents + ` (
			id BIGSERIAL PRIMARY KEY,
			timestamp TIMESTAMP NOT NULL,
			metadata JSONB,
			user_id BIGINT NOT NULL,
			type_id BIGINT NOT NULL
		);`,
	}
)

func CreateDB(conf config.DB) (*pgxpool.Pool, lib.Error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&pool_max_conns=50", conf.User, conf.Password, conf.Host, conf.Port, conf.Database)
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, lib.Err(err)
	}
	return pool, nil
}
