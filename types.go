package migrationmanager

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

// Migration represents a migration
type Migration struct {
	bun.BaseModel `bun:"managed_migrations"`
	ID            int          `json:"id" sql:",pk"`
	Name          string       `json:"name"`
	CreatedAt     time.Time    `json:"createdAt"`
	Up            func() error `bun:"-"`
	Down          func() error `bun:"-"`
}

// BeforeInsert hook
func (m *Migration) BeforeInsert(ctx context.Context) (context.Context, error) {
	m.CreatedAt = time.Now()
	return ctx, nil
}
