package migrationmanager

import (
	"context"
	"time"
)

// Migration represents a migration
type Migration struct {
	tableName struct{}     `pg:"managed_migrations"`
	ID        int          `json:"id" sql:",pk"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"createdAt"`
	Up        func() error `pg:"-"`
	Down      func() error `pg:"-"`
}

// BeforeInsert hook
func (m *Migration) BeforeInsert(ctx context.Context) (context.Context, error) {
	m.CreatedAt = time.Now()
	return ctx, nil
}
