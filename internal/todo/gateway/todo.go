package gateway

import (
	"context"
	"time"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/todo"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/apperr"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/TokujouKaisenDonburi/optical-backend/pkg/psql"
	"github.com/google/uuid"
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

func (r *TodoPsqlRepository) CreateList(
	ctx context.Context,
	list *todo.List,
) error {
	return db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		exists, err := psql.IsUserInCalendarMembers(ctx, tx, list.UserId, list.CalendarId)
		if err != nil {
			return err
		}
		if !exists {
			return apperr.ForbiddenError("the user is not in calendar members")
		}
		query := `
			INSERT INTO todo_lists(id, user_id, calendar_id, name, created_at, updated_at)
			VALUES(:id, :userId, :calendarId, :name, :createdAt, :updatedAt)
		`
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"id":         list.Id,
			"userId":     list.UserId,
			"calendarId": list.CalendarId,
			"name":       list.Name,
			"createdAt":  time.Now().UTC(),
			"updatedAt":  time.Now().UTC(),
		})
		return err
	})
}

func (r *TodoPsqlRepository) AddItem(
	ctx context.Context,
	listId uuid.UUID,
	todoItem *todo.Item,
) error {
	return db.RunInTx(r.db, func(tx *sqlx.Tx) error {
		todoList, err := psql.FindTodoListById(ctx, tx, listId)
		if err != nil {
			return err
		}
		exists, err := psql.IsUserInCalendarMembers(ctx, tx, todoItem.UserId, todoList.CalendarId)
		if err != nil {
			return err
		}
		if !exists {
			return apperr.ForbiddenError("the user is not in calendar members")
		}
		query := `
			INSERT INTO todo_items(id, todo_list_id, user_id, name, is_done, created_at, updated_at)
			VALUES(:id, :todoListId, :userId, :name, :isDone, :createdAt, :updatedAt)
		`
		_, err = tx.NamedExecContext(ctx, query, map[string]any{
			"id":         todoItem.Id,
			"todoListId": todoList.Id,
			"userId":     todoItem.UserId,
			"name":       todoItem.Name,
			"isDone":     todoItem.IsDone,
			"createdAt":  time.Now(),
			"updatedAt":  time.Now(),
		})
		return err
	})
}
