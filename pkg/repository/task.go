package repository

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/gihyodocker/taskapp/pkg/model"
)

type Task interface {
	Upsert(ctx context.Context, m *model.Task) error
	DeleteByID(ctx context.Context, id string) (int64, error)
	FindByID(ctx context.Context, id string) (*model.Task, error)
	FindAll(ctx context.Context) ([]*model.Task, error)
}

func NewTask(db *sql.DB) Task {
	return &task{
		db: db,
	}
}

type task struct {
	db *sql.DB
}

func (r task) Upsert(ctx context.Context, m *model.Task) error {
	return m.Upsert(
		ctx,
		r.db,
		boil.Whitelist("title", "content", "status", "updated"),
		boil.Infer(),
	)
}

func (r task) DeleteByID(ctx context.Context, id string) (int64, error) {
	return model.Tasks(
		qm.Where("id = ?", id),
	).DeleteAll(ctx, r.db)
}

func (r task) FindByID(ctx context.Context, id string) (*model.Task, error) {
	return model.FindTask(ctx, r.db, id)
}

func (r task) FindAll(ctx context.Context) ([]*model.Task, error) {
	return model.Tasks(qm.OrderBy("updated DESC")).All(ctx, r.db)
}
