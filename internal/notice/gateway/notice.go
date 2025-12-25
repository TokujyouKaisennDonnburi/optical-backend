package gateway

import "github.com/jmoiron/sqlx"

type NoticePsqlRepository struct {
	db *sqlx.DB
}

func NewNoticePsqlRepository(db *sqlx.DB) *NoticePsqlRepository {
	return &NoticePsqlRepository{
		db: db,
	}
}
