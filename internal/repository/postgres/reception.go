package postgres

import (
	"database/sql"
	"errors"
	"final/internal/domain"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ReceptionRepository struct {
	db *sqlx.DB
}

func NewReceptionRepository(db *sqlx.DB) *ReceptionRepository {
	return &ReceptionRepository{db: db}
}

func (r *ReceptionRepository) CreateReception(reception *domain.Reception) (*domain.Reception, error) {
	receptionID, err := r.getLastReceptionIDInProgressForPVZ(reception.PvzId)
	if err != nil {
		return nil, err
	} else if receptionID != nil {
		return nil, fmt.Errorf("previous reception is still in progress for pvz with id: %s", *reception.PvzId)
	}

	query := `INSERT INTO receptions (id, date_time, pvz_id, status)
              VALUES (COALESCE($1, gen_random_uuid()), COALESCE($2, now()), $3,  COALESCE($4, status('in_progress')))
              RETURNING id, date_time, status`

	row := r.db.QueryRow(query, reception.ID, reception.DateTime, reception.PvzId, reception.Status)
	err = row.Scan(&reception.ID, &reception.DateTime, &reception.Status)
	if err != nil {
		return nil, err
	}
	return reception, nil
}

func (r *ReceptionRepository) getLastReceptionIDInProgressForPVZTx(tx *sqlx.Tx, pvzID *string) (*string, error) {
	query := `SELECT id FROM receptions WHERE pvz_id = $1 AND status = 'in_progress'`

	var receptionID string
	err := tx.QueryRow(query, *pvzID).Scan(&receptionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &receptionID, nil
}

func (r *ReceptionRepository) getLastReceptionIDInProgressForPVZ(pvzID *string) (*string, error) {
	query := `SELECT id FROM receptions WHERE pvz_id = $1 AND status = 'in_progress'`

	var receptionID string
	err := r.db.QueryRow(query, *pvzID).Scan(&receptionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &receptionID, nil
}

func (r *ReceptionRepository) CloseReception(pvzID *string) (*domain.Reception, error) {
	receptionID, err := r.getLastReceptionIDInProgressForPVZ(pvzID)
	if err != nil {
		return nil, err
	} else if receptionID == nil {
		return nil, domain.ErrReceptionNotFound
	}

	updateQuery := `
					UPDATE receptions 
					SET status = 'close' 
					WHERE id = $1 
					RETURNING id, date_time, pvz_id, status`
	var reception domain.Reception
	err = r.db.QueryRow(updateQuery, *receptionID).Scan(&reception.ID, &reception.DateTime, &reception.PvzId, &reception.Status)
	if err != nil {
		return nil, err
	}
	return &reception, nil
}
