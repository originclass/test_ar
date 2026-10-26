// 数据库访问层:连接 PostgreSQL、自动建表、问题的增查改。
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store 封装所有数据库操作
type Store struct {
	pool *pgxpool.Pool
}

// NewStore 建立连接池并初始化表结构
func NewStore(ctx context.Context) (*Store, error) {
	// 优先读取环境变量 PATROL_DB_DSN,便于换环境;默认连本机 PostgreSQL(3306 为本机自定义端口)
	dsn := os.Getenv("PATROL_DB_DSN")
	if dsn == "" {
		dsn = "host=localhost port=3306 user=postgres password=123456 dbname=patrol_ar sslmode=disable TimeZone=Asia/Shanghai"
	}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	if err := migrate(ctx, pool); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return &Store{pool: pool}, nil
}

// migrate 启动时自动建表,保证服务拿到空库也能直接运行
func migrate(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS issues (
	id          TEXT PRIMARY KEY,
	title       TEXT NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	priority    TEXT NOT NULL,
	status      TEXT NOT NULL,
	pos_x       DOUBLE PRECISION NOT NULL DEFAULT 0,
	pos_y       DOUBLE PRECISION NOT NULL DEFAULT 0,
	pos_z       DOUBLE PRECISION NOT NULL DEFAULT 0,
	created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
)`)
	return err
}

// Create 保存一条新问题
func (s *Store) Create(ctx context.Context, issue *Issue) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO issues (id, title, description, priority, status, pos_x, pos_y, pos_z, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $9)`,
		issue.ID, issue.Title, issue.Description, issue.Priority, issue.Status,
		issue.Position.X, issue.Position.Y, issue.Position.Z, issue.CreatedAt)
	return err
}

// List 返回问题列表,按创建时间倒序(最新上报的排最前)
func (s *Store) List(ctx context.Context) ([]Issue, error) {
	rows, err := s.pool.Query(ctx, `
SELECT id, title, description, priority, status, pos_x, pos_y, pos_z, created_at, updated_at
FROM issues ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	issues := []Issue{}
	for rows.Next() {
		var it Issue
		if err := rows.Scan(&it.ID, &it.Title, &it.Description, &it.Priority, &it.Status,
			&it.Position.X, &it.Position.Y, &it.Position.Z, &it.CreatedAt, &it.UpdatedAt); err != nil {
			return nil, err
		}
		issues = append(issues, it)
	}
	return issues, rows.Err()
}

// UpdateStatus 修改问题状态;返回是否确实更新到了记录(用于判断 404)
func (s *Store) UpdateStatus(ctx context.Context, id, status string) (bool, error) {
	tag, err := s.pool.Exec(ctx,
		`UPDATE issues SET status = $1, updated_at = $2 WHERE id = $3`,
		status, time.Now(), id)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}
