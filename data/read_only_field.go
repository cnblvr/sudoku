package data

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type IDint64Getter interface {
	ID() int64
}

type IDuuidGetter interface {
	ID() uuid.UUID
}

type CreatedAtGetter interface {
	CreatedAt() time.Time
}

type UpdatedAtGetter interface {
	UpdatedAt() time.Time
}
