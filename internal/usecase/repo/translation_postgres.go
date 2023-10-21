package repo

import (
	"context"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/pkg/postgres"
)

const _defaultEntityCap = 64

// TranslationRepo -.
type TranslationRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *TranslationRepo {
	return &TranslationRepo{pg}
}

// GetHistory -.
func (r *TranslationRepo) GetHistory(ctx context.Context) ([]entity.Translation, error) {
	//sql, _, err := r.Db.
	//	Select("source, destination, original, translation").
	//	From("history").
	//	ToSql()
	//if err != nil {
	//	return nil, fmt.Errorf("TranslationRepo - GetHistory - r.Db: %w", err)
	//}
	//
	//rows, err := r.Pool.Query(ctx, sql)
	//if err != nil {
	//	return nil, fmt.Errorf("TranslationRepo - GetHistory - r.Pool.Query: %w", err)
	//}
	//defer rows.Close()
	//
	//entities := make([]entity.Translation, 0, _defaultEntityCap)
	//
	//for rows.Next() {
	//	e := entity.Translation{}
	//
	//	err = rows.Scan(&e.Source, &e.Destination, &e.Original, &e.Translation)
	//	if err != nil {
	//		return nil, fmt.Errorf("TranslationRepo - GetHistory - rows.Scan: %w", err)
	//	}
	//
	//	entities = append(entities, e)
	//}

	//return entities, nil
	return nil, nil
}

// Store -.
func (r *TranslationRepo) Store(ctx context.Context, t entity.Translation) error {
	//sql, args, err := r.Db.
	//	Insert("history").
	//	Columns("source, destination, original, translation").
	//	Values(t.Source, t.Destination, t.Original, t.Translation).
	//	ToSql()
	//if err != nil {
	//	return fmt.Errorf("TranslationRepo - Store - r.Db: %w", err)
	//}
	//
	//_, err = r.Pool.Exec(ctx, sql, args...)
	//if err != nil {
	//	return fmt.Errorf("TranslationRepo - Store - r.Pool.Exec: %w", err)
	//}

	return nil
}
