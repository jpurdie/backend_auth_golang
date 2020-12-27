package authapi

import (
	"context"
	"github.com/lib/pq"
	"time"
)

// Base contains common fields for all tables
type Base struct {
	ID        int       `json:"-"  db:"id"`
	CreatedAt pq.NullTime `json:"-"  db:"created_at"`
	UpdatedAt pq.NullTime `json:"-"  db:"updated_at"`
	DeletedAt pq.NullTime `json:"-"  db:"deleted_at"`
}

// ListQuery holds company/location data used for list db queries
type ListQuery struct {
	Query string
	ID    int
}

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (b *Base) BeforeInsert(ctx context.Context) (context.Context, error) {
	now := time.Now()
	b.CreatedAt.Time = now
	b.UpdatedAt.Time = now
	return ctx, nil
}

// BeforeUpdate hooks into update operations, setting updatedAt to current time
func (b *Base) BeforeUpdate(ctx context.Context) (context.Context, error) {
	b.UpdatedAt.Time = time.Now()
	return ctx, nil
}
