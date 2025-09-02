package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/samber/lo"

	"github.com/DanilaKorobkov/defi-monitoring/internal/domain"
)

// Executor allows repository work with sqlx.DB and sqlx.Tx as driver.
//
//nolint:revive // Unnecessary comments for technical interfaces.
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
}

type SubjectsRepository struct {
	db Executor
}

func NewSubjectsRepository(db Executor) *SubjectsRepository {
	return &SubjectsRepository{
		db: db,
	}
}

func (p SubjectsRepository) Add(ctx context.Context, subject domain.Subject) error {
	model, err := newSubjectModel(subject)
	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, queryAddSubject, model.TelegramUserID, model.Payload)
	if err != nil {
		return fmt.Errorf("ExecContext: %w", err)
	}

	return nil
}

func (p SubjectsRepository) GetAll(ctx context.Context) ([]domain.Subject, error) {
	var models []subjectModel

	err := p.db.SelectContext(ctx, &models, queryGetAllSubjects)
	if err != nil {
		return nil, fmt.Errorf("SelectContext: %w", err)
	}

	if len(models) == 0 {
		return nil, nil
	}

	subjects := lo.Map(models, func(item subjectModel, _ int) domain.Subject {
		return item.mustToSubject()
	})

	return subjects, nil
}

const queryGetAllSubjects = `
SELECT 
    telegram_user_id, 
    payload 
FROM 
    subjects
`

const queryAddSubject = `
INSERT INTO 
    subjects (telegram_user_id, payload)
VALUES 
    ($1, $2)
ON CONFLICT 
    (telegram_user_id)
DO UPDATE SET 
	payload=EXCLUDED.payload
`
