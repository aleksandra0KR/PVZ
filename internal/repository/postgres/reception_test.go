package postgres

import (
	"context"
	"database/sql"
	"final/internal/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateReception(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			log.Errorf("failed to close stub database connection: %v", err)
		}
	}(db)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewReceptionRepository(sqlxDB)

	pvzId := "pvz-123"
	receptionId := "reception-123"
	dateTime := time.Now()
	status := "in_progress"
	ctx := context.Background()
	t.Run("Successful_Create_Reception", func(t *testing.T) {
		inputReception := &domain.Reception{
			ID:       nil,
			DateTime: nil,
			PvzId:    &pvzId,
			Status:   &status,
		}

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT id FROM receptions").
			WithArgs(pvzId).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		mock.ExpectQuery("INSERT INTO receptions").
			WithArgs(nil, nil, pvzId, "in_progress").
			WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "status"}).AddRow(receptionId, dateTime, "in_progress"))

		reception, err := repo.CreateReception(ctx, inputReception)
		assert.NoError(t, err)
		assert.NotNil(t, reception)
		assert.Equal(t, receptionId, *reception.ID)
		assert.Equal(t, &pvzId, reception.PvzId)
		assert.Equal(t, &status, reception.Status)
	})

	t.Run("Failure_Previous_Reception_In_Progress", func(t *testing.T) {
		inputReception := &domain.Reception{
			ID:       nil,
			DateTime: nil,
			PvzId:    &pvzId,
			Status:   &status,
		}

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT id FROM receptions").
			WithArgs(pvzId).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(receptionId))

		reception, err := repo.CreateReception(ctx, inputReception)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrCreateReceptionBecauseOfPreviousReception, err)
		assert.Nil(t, reception)
	})

	t.Run("Failure_DB_Error", func(t *testing.T) {
		inputReception := &domain.Reception{
			ID:       nil,
			DateTime: nil,
			PvzId:    &pvzId,
			Status:   &status,
		}

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT id FROM receptions").
			WithArgs(pvzId).
			WillReturnError(sql.ErrConnDone)

		reception, err := repo.CreateReception(ctx, inputReception)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrCreateReception, err)
		assert.Nil(t, reception)
	})
}

func TestCloseReception(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			log.Error(err)
		}
	}(db)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewReceptionRepository(sqlxDB)

	pvzId := "pvz-123"
	receptionId := "reception-123"
	dateTime := time.Now()
	status := "close"

	ctx := context.Background()
	t.Run("Successful_Close_Reception", func(t *testing.T) {
		mock.ExpectQuery("SELECT id FROM receptions").
			WithArgs(pvzId).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(receptionId))

		mock.ExpectQuery("UPDATE receptions").
			WithArgs(receptionId).
			WillReturnRows(sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status"}).AddRow(receptionId, dateTime, pvzId, "close"))

		reception, err := repo.CloseReception(ctx, &pvzId)
		assert.NoError(t, err)
		assert.NotNil(t, reception)
		assert.Equal(t, receptionId, *reception.ID)
		assert.Equal(t, &pvzId, reception.PvzId)
		assert.Equal(t, &status, reception.Status)
	})

	t.Run("Failure_Reception_Not_Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id FROM receptions").
			WithArgs(pvzId).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		reception, err := repo.CloseReception(ctx, &pvzId)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrReceptionNotFound, err)
		assert.Nil(t, reception)
	})

	t.Run("Failure_DB_Errorr", func(t *testing.T) {
		mock.ExpectQuery("SELECT id FROM receptions").
			WithArgs(pvzId).
			WillReturnError(sql.ErrConnDone)

		reception, err := repo.CloseReception(ctx, &pvzId)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrCloseReception, err)
		assert.Nil(t, reception)
	})
}
