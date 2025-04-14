package postgres

import (
	"context"
	"database/sql"
	"errors"
	"final/internal/domain"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type ReceptionRepository struct {
	db *sqlx.DB
}

func NewReceptionRepository(db *sqlx.DB) *ReceptionRepository {
	return &ReceptionRepository{db: db}
}

func (r *ReceptionRepository) CreateReception(ctx context.Context, reception *domain.Reception) (*domain.Reception, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error(err)
		return nil, domain.ErrCreateReception
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Errorf("rollback failed: %v", rollbackErr)
			}
		} else {
			commitErr := tx.Commit()
			if commitErr != nil {
				log.Errorf("commit failed: %v", commitErr)
			}
		}
	}()

	receptionID, err := r.getLastReceptionIDInProgressForPVZTx(ctx, tx, reception.PvzId)
	if receptionID != nil {
		return nil, domain.ErrCreateReceptionBecauseOfPreviousReception
	} else if err != nil {
		return nil, domain.ErrCreateReception
	}

	query := `INSERT INTO receptions (id, date_time, pvz_id, status)
              VALUES (COALESCE($1, gen_random_uuid()), COALESCE($2, now()), $3,  COALESCE($4, status('in_progress')))
              RETURNING id, date_time, status`

	row := tx.QueryRowContext(ctx, query, reception.ID, reception.DateTime, reception.PvzId, reception.Status)
	err = row.Scan(&reception.ID, &reception.DateTime, &reception.Status)
	if err != nil {
		log.Error(err)
		return nil, domain.ErrCreateReception
	}
	return reception, nil
}

func (r *ReceptionRepository) getLastReceptionIDInProgressForPVZTx(ctx context.Context, tx *sqlx.Tx, pvzID *string) (*string, error) {
	query := `SELECT id FROM receptions WHERE pvz_id = $1 AND status = 'in_progress'`

	var receptionID string
	err := tx.QueryRowContext(ctx, query, *pvzID).Scan(&receptionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Error(err)
		return nil, err
	}
	return &receptionID, nil
}

func (r *ReceptionRepository) getLastReceptionIDInProgressForPVZ(ctx context.Context, pvzID *string) (*string, error) {
	query := `SELECT id FROM receptions WHERE pvz_id = $1 AND status = 'in_progress'`

	var receptionID string
	err := r.db.QueryRowContext(ctx, query, *pvzID).Scan(&receptionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Error(err)
		return nil, err
	}
	return &receptionID, nil
}

func (r *ReceptionRepository) CloseReception(ctx context.Context, pvzID *string) (*domain.Reception, error) {
	receptionID, err := r.getLastReceptionIDInProgressForPVZ(ctx, pvzID)
	if err != nil {
		return nil, domain.ErrCloseReception
	} else if receptionID == nil {
		return nil, domain.ErrReceptionNotFound
	}

	updateQuery := `
					UPDATE receptions 
					SET status = 'close' 
					WHERE id = $1 
					RETURNING id, date_time, pvz_id, status`
	var reception domain.Reception
	err = r.db.QueryRowContext(ctx, updateQuery, *receptionID).Scan(&reception.ID, &reception.DateTime, &reception.PvzId, &reception.Status)
	if err != nil {
		log.Error(err)
		return nil, domain.ErrCloseReception
	}
	return &reception, nil
}
