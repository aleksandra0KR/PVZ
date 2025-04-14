package postgres

import (
	"database/sql"
	"final/internal/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreatePVZ(t *testing.T) {
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
	pvzRepo := NewPvzRepository(sqlxDB)
	allowedCity := "Москва"

	t.Run("Successful_Create_PVZ", func(t *testing.T) {
		pvz := &domain.PVZ{
			City: &allowedCity,
		}

		mock.ExpectQuery("INSERT INTO pvz").
			WithArgs(nil, nil, "Москва").
			WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date"}).AddRow("1", time.Now()))

		result, err := pvzRepo.CreatePVZ(pvz)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Москва", *result.City)
	})

	t.Run("Failure_DB_Error", func(t *testing.T) {
		pvz := &domain.PVZ{
			City: &allowedCity,
		}

		mock.ExpectQuery("INSERT INTO pvz").
			WithArgs(nil, nil, "Москва").
			WillReturnError(sql.ErrConnDone)

		result, err := pvzRepo.CreatePVZ(pvz)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrCreatePVZ, err)
		assert.Nil(t, result)
	})
}

func TestGetAllPVZ(t *testing.T) {
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
	pvzRepo := NewPvzRepository(sqlxDB)

	t.Run("Successful_Get_All_PVZ", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "registration_date", "city"}).
			AddRow("1", time.Now(), "Москва").
			AddRow("2", time.Now(), "Казань")

		mock.ExpectQuery("SELECT \\* FROM pvz").WillReturnRows(rows)

		result, err := pvzRepo.GetAllPVZ()
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Москва", *result[0].City)
		assert.Equal(t, "Казань", *result[1].City)
	})

	t.Run("Failure_DB_Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT * FROM pvz").WillReturnError(sql.ErrConnDone)

		result, err := pvzRepo.GetAllPVZ()
		assert.Error(t, err)
		assert.Equal(t, domain.ErrGetPVZ, err)
		assert.Nil(t, result)
	})
}

func TestGetPVZInfo(t *testing.T) {
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
	pvzRepo := NewPvzRepository(sqlxDB)

	t.Run("Successful_Get_PVZ_Info", func(t *testing.T) {
		startDate := time.Now().AddDate(0, 0, -1)
		endDate := time.Now()
		offset := 0
		limit := 10

		rows := sqlmock.NewRows([]string{
			"pvz_id", "registration_date", "city",
			"reception_id", "date_time", "status",
			"product_id", "product_datetime", "type", "product_reception_id",
		}).
			AddRow("1", time.Now(), "Москва", "1", time.Now(), "in_progress", "1", time.Now(), "электроника", "1").
			AddRow("1", time.Now(), "Москва", "2", time.Now(), "close", nil, nil, nil, nil)

		mock.ExpectQuery("SELECT pvz.id as pvz_id, pvz.registration_date, pvz.city, receptions.id as reception_id, receptions.date_time, receptions.status, products.id as product_id, products.date_time as product_datetime, products.type, products.reception_id as product_reception_id FROM pvz LEFT JOIN receptions ON pvz.id = receptions.pvz_id LEFT JOIN products ON receptions.id = products.reception_id WHERE receptions.date_time >= \\$1 AND receptions.date_time <= \\$2 ORDER BY pvz.registration_date DESC LIMIT 10 OFFSET 0").
			WithArgs(startDate, endDate).
			WillReturnRows(rows)

		result, err := pvzRepo.GetPVZInfo(&startDate, &endDate, offset, limit)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Москва", *result[0].PVZ.City)
		assert.Len(t, result[0].Receptions, 1)
		assert.Equal(t, "in_progress", *result[0].Receptions[0].Reception.Status)
	})

	t.Run("Failure_DB_Error", func(t *testing.T) {
		startDate := time.Now().AddDate(0, 0, -1)
		endDate := time.Now()
		offset := 0
		limit := 10

		mock.ExpectQuery("SELECT pvz.id as pvz_id, pvz.registration_date, pvz.city, receptions.id as reception_id, receptions.date_time, receptions.status, products.id as product_id, products.date_time as product_datetime, products.type, products.reception_id as product_reception_id FROM pvz LEFT JOIN receptions ON pvz.id = receptions.pvz_id LEFT JOIN products ON receptions.id = products.reception_id WHERE receptions.date_time >= $1 AND receptions.date_time <= $2 ORDER BY pvz.registration_date DESC LIMIT $3 OFFSET $4").
			WithArgs(startDate, endDate, limit, offset).
			WillReturnError(sql.ErrConnDone)

		result, err := pvzRepo.GetPVZInfo(&startDate, &endDate, offset, limit)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
