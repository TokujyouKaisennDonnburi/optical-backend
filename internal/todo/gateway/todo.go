package gateway

import (
	"github.com/jmoiron/sqlx"
)

type TodoPsqlRepository struct {
	db *sqlx.DB
}

func NewTodoPsqlRepository(db *sqlx.DB) *TodoPsqlRepository {
	if db == nil {
		panic("db is nil")
	}
	return &TodoPsqlRepository{
		db: db,
	}
}
