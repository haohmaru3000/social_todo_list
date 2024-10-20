package storage

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type sqlStore struct {
	db     *gorm.DB
	tracer trace.Tracer
}

func NewSQLStore(db *gorm.DB) *sqlStore {
	return &sqlStore{
		db:     db,
		tracer: otel.Tracer("Item.Storage"),
	}
}
