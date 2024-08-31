package pgrepo

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

// Repository is ...
type Repository struct {
	logger *zap.Logger
	db     *pgxpool.Pool
}

// NewRepository is ...
func NewRepository(db *pgxpool.Pool, logger *zap.Logger) *Repository {
	logger = logger.Named("repo")
	return &Repository{
		logger: logger,
		db:     db,
	}
}
