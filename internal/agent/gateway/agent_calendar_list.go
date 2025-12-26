package gateway

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent"
	"github.com/TokujouKaisenDonburi/optical-backend/internal/agent/transact"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CalendarQueryModel struct {
	Id       uuid.UUID      `db:"id"`
	Name     string         `db:"name"`
	Color    string         `db:"color"`
	ImageUrl sql.NullString `db:"image_url"`
	Members  string         `db:"members"`
	Options  string         `db:"options"`
}

func (r *AgentQueryPsqlRepository) FindCalendarByUserId(
	ctx context.Context,
	userId uuid.UUID,
) ([]agent.AnalyzableCalendar, error) {
	calendars := []agent.AnalyzableCalendar{}
	err := transact.Transact(ctx, func(tx *sqlx.Tx) error {
		var models []CalendarQueryModel
		query := `
			SELECT
				c.id,
				c.name,
				c.color,
				ci.url AS image_url,
				COALESCE((
					SELECT JSON_AGG(JSON_BUILD_OBJECT(
						'user_id', cm.user_id, 
						'user_name', u.name, 
						'joined_at', cm.joined_at
					))
					FROM
						calendar_members cm
					JOIN users u
						ON cm.user_id = u.id
					WHERE cm.calendar_id = c.id
				), '[]') AS members,
				COALESCE((
					SELECT JSON_AGG(JSON_BUILD_OBJECT(
						'option_id', co.option_id,
						'option_name', o.name
					))
					FROM
						calendar_options co
					JOIN options o
						ON co.option_id = o.id
					WHERE co.calendar_id = c.id
				), '[]') AS options
			FROM (
				SELECT * FROM calendars
				WHERE calendars.id IN (
					SELECT calendar_id FROM calendar_members
					WHERE calendar_members.user_id = $1
				)
			) c
			LEFT JOIN calendar_images ci
				ON c.image_id = ci.id
			ORDER BY c.id
		`
		err := tx.SelectContext(ctx, &models, query, userId)
		if err != nil {
			return err
		}
		for _, model := range models {
			var members []agent.AnalyzableMember
			if err := json.Unmarshal([]byte(model.Members), &members); err != nil {
				return err
			}
			var options []agent.AnalyzableOption
			if err := json.Unmarshal([]byte(model.Options), &options); err != nil {
				return err
			}
			calendars = append(calendars, agent.AnalyzableCalendar{
				Id:       model.Id.String(),
				Name:     model.Name,
				Color:    model.Color,
				ImageUrl: model.ImageUrl.String,
				Members:  members,
				Options:  options,
			})
		}
		return nil
	})
	return calendars, err
}

func (r *AgentQueryPsqlRepository) FindCalendarByIdAndUserId(
	ctx context.Context,
	userId, calendarId uuid.UUID,
) (*agent.AnalyzableCalendar, error) {
	calendar := agent.AnalyzableCalendar{}
	err := transact.Transact(ctx, func(tx *sqlx.Tx) error {
		var model CalendarQueryModel
		query := `
			SELECT
				c.id,
				c.name,
				c.color,
				ci.url AS image_url,
				COALESCE((
					SELECT JSON_AGG(JSON_BUILD_OBJECT(
						'user_id', cm.user_id, 
						'user_name', u.name, 
						'joined_at', cm.joined_at
					))
					FROM
						calendar_members cm
					JOIN users u
						ON cm.user_id = u.id
					WHERE cm.calendar_id = c.id
				), '[]') AS members,
				COALESCE((
					SELECT JSON_AGG(JSON_BUILD_OBJECT(
						'option_id', co.option_id,
						'option_name', o.name
					))
					FROM
						calendar_options co
					JOIN options o
						ON co.option_id = o.id
					WHERE co.calendar_id = c.id
				), '[]') AS options
			FROM (
				SELECT * FROM calendars
				WHERE calendars.id = $2
				AND calendars.id IN (
					SELECT calendar_id FROM calendar_members
					WHERE calendar_members.user_id = $1
				)
			) c
			LEFT JOIN calendar_images ci
				ON c.image_id = ci.id
			ORDER BY c.id
		`
		err := tx.GetContext(ctx, &model, query, userId, calendarId)
		if err != nil {
			return err
		}
		var members []agent.AnalyzableMember
		if err := json.Unmarshal([]byte(model.Members), &members); err != nil {
			return err
		}
		var options []agent.AnalyzableOption
		if err := json.Unmarshal([]byte(model.Options), &options); err != nil {
			return err
		}
		calendar = agent.AnalyzableCalendar{
			Id:       model.Id.String(),
			Name:     model.Name,
			Color:    model.Color,
			ImageUrl: model.ImageUrl.String,
			Members:  members,
			Options:  options,
		}
		return nil
	})
	return &calendar, err
}
