//go:build integration

package integration

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const (
	envTestDBDSN         = "TEST_DB_DSN"
	envTestRedisAddr     = "TEST_REDIS_ADDR"
	envTestRedisPassword = "TEST_REDIS_PASSWORD"
	envTestRedisDB       = "TEST_REDIS_DB"
)

var migrationsOnce sync.Once

func testDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := os.Getenv(envTestDBDSN)
	if dsn == "" {
		t.Skip("TEST_DB_DSN is not set; skipping integration tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Fatalf("failed to ping db: %v", err)
	}

	migrationsOnce.Do(func() {
		if err := applyMigrations(pool); err != nil {
			t.Fatalf("failed to apply migrations: %v", err)
		}
	})

	return pool
}

func applyMigrations(pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	root, err := repoRoot()
	if err != nil {
		return err
	}
	files := []string{
		filepath.Join(root, "migrations", "001_init.sql"),
		filepath.Join(root, "migrations", "002_indexes.sql"),
	}

	for _, path := range files {
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := execSQL(ctx, pool, string(raw)); err != nil {
			return err
		}
	}
	return nil
}

func execSQL(ctx context.Context, pool *pgxpool.Pool, sql string) error {
	statements := strings.Split(sql, ";")
	for _, stmt := range statements {
		trimmed := strings.TrimSpace(stmt)
		if trimmed == "" {
			continue
		}
		if _, err := pool.Exec(ctx, trimmed); err != nil {
			return err
		}
	}
	return nil
}

func truncateTables(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := pool.Exec(ctx, `
		TRUNCATE TABLE location_check_incidents, location_checks, incidents CASCADE
	`); err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
}

func testRedis(t *testing.T) *redis.Client {
	t.Helper()

	addr := os.Getenv(envTestRedisAddr)
	if addr == "" {
		t.Skip("TEST_REDIS_ADDR is not set; skipping redis integration tests")
	}

	db := 0
	if value := os.Getenv(envTestRedisDB); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			db = parsed
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv(envTestRedisPassword),
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		t.Fatalf("failed to ping redis: %v", err)
	}

	if err := client.FlushDB(ctx).Err(); err != nil {
		_ = client.Close()
		t.Fatalf("failed to flush redis db: %v", err)
	}

	return client
}

func repoRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Dir(filepath.Dir(wd)), nil
}
