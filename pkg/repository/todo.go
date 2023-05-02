package repository

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/gihyodocker/todoapp/pkg/model"
)

type TODO interface {
	Upsert(ctx context.Context, m *model.Todo) error
	DeleteByID(ctx context.Context, id string) (int64, error)
	FindByID(ctx context.Context, id string) (*model.Todo, error)
}

func NewTODO(db *sql.DB) TODO {
	return &todo{
		db: db,
	}
}

type todo struct {
	db *sql.DB
}

func (r todo) Upsert(ctx context.Context, m *model.Todo) error {
	return m.Upsert(
		ctx,
		r.db,
		boil.Whitelist("title", "content", "status", "updated_at"),
		boil.Infer(),
	)
}

func (r todo) DeleteByID(ctx context.Context, id string) (int64, error) {
	return model.Todos(
		qm.Where("id = ?", id),
	).DeleteAll(ctx, r.db)
}

func (r todo) FindByID(ctx context.Context, id string) (*model.Todo, error) {
	return model.FindTodo(ctx, r.db, id)
}
