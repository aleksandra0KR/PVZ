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

func TestCreateProduct(t *testing.T) {
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
	receptionRepo := NewReceptionRepository(sqlxDB)
	productRepo := NewProductRepository(sqlxDB, receptionRepo)
	productType := "электроника"
	pvzId := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	receptionId := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	productId := "1"
	dateTime := time.Now()

	t.Run("Successful_Create_Product", func(t *testing.T) {
		ctx := context.Background()
		inputProduct := &domain.InputProduct{
			Type:  &productType,
			PvzId: &pvzId,
		}

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT id FROM receptions").
			WithArgs(pvzId).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(receptionId))

		mock.ExpectQuery("INSERT INTO products").
			WithArgs(receptionId, productType).
			WillReturnRows(sqlmock.NewRows([]string{"id", "date_time"}).AddRow(productId, dateTime))

		mock.ExpectCommit()

		product, err := productRepo.CreateProduct(ctx, inputProduct)
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, receptionId, *product.ReceptionID)
		assert.Equal(t, productType, *product.Type)
		assert.Equal(t, &productId, product.ID)
		assert.Equal(t, &dateTime, product.DateTime)
	})

	t.Run("Failure_Reception_Not_Found", func(t *testing.T) {
		ctx := context.Background()
		inputProduct := &domain.InputProduct{
			Type:  &productType,
			PvzId: &pvzId,
		}

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT id FROM receptions").
			WithArgs(pvzId).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectRollback()

		product, err := productRepo.CreateProduct(ctx, inputProduct)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrReceptionNotFound, err)
		assert.Nil(t, product)
	})

	t.Run("Failure_DB_Error", func(t *testing.T) {
		ctx := context.Background()
		inputProduct := &domain.InputProduct{
			Type:  &productType,
			PvzId: &pvzId,
		}

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT id FROM receptions").
			WithArgs(pvzId).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(receptionId))

		mock.ExpectQuery("INSERT INTO products").
			WithArgs(receptionId, productType).
			WillReturnError(sql.ErrConnDone)

		mock.ExpectRollback()

		product, err := productRepo.CreateProduct(ctx, inputProduct)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrCreateProduct, err)
		assert.Nil(t, product)
	})
}

func TestDeleteLastProductForPVZ(t *testing.T) {
	ctx := context.Background()
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
	receptionRepo := NewReceptionRepository(sqlxDB)
	productRepo := NewProductRepository(sqlxDB, receptionRepo)
	productType := "электроника"
	pvzId := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	receptionId := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	productId := "1"
	dateTime := time.Now()

	t.Run("Successful_Delete_Product", func(t *testing.T) {
		inputProduct := &domain.InputProduct{
			Type:  &productType,
			PvzId: &pvzId,
		}

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT id FROM receptions").
			WithArgs(pvzId).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(receptionId))

		mock.ExpectQuery("INSERT INTO products").
			WithArgs(receptionId, productType).
			WillReturnRows(sqlmock.NewRows([]string{"id", "date_time"}).AddRow(productId, dateTime))

		mock.ExpectCommit()

		product, err := productRepo.CreateProduct(ctx, inputProduct)
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, receptionId, *product.ReceptionID)
		assert.Equal(t, productType, *product.Type)
		assert.Equal(t, &productId, product.ID)
		assert.Equal(t, &dateTime, product.DateTime)
	})

	t.Run("Failure_Reception_Not_Found", func(t *testing.T) {
		t.Parallel()
		db, mock, err = sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				log.Error(err)
			}
		}(db)

		sqlxDB = sqlx.NewDb(db, "sqlmock")
		receptionRepo = NewReceptionRepository(sqlxDB)
		productRepo = NewProductRepository(sqlxDB, receptionRepo)
		pvzId = "1"

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT id FROM receptions").
			WithArgs(pvzId).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectRollback()

		err = productRepo.DeleteLastProductForPVZ(ctx, &pvzId)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrReceptionNotFound, err)
	})

	t.Run("Failure_No_Products_In_Reception", func(t *testing.T) {
		t.Parallel()
		db, mock, err = sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				log.Error(err)
			}
		}(db)

		sqlxDB = sqlx.NewDb(db, "sqlmock")
		receptionRepo = NewReceptionRepository(sqlxDB)
		productRepo = NewProductRepository(sqlxDB, receptionRepo)
		pvzId = "1"
		receptionId = "1"

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT id FROM receptions").
			WithArgs(pvzId).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(receptionId))

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM products").
			WithArgs(receptionId).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectRollback()

		err = productRepo.DeleteLastProductForPVZ(ctx, &pvzId)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrNoProductsInReception, err)
	})

	t.Run("Failure_DB_Error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT id FROM receptions").
			WithArgs(pvzId).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(receptionId))

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM products").
			WithArgs(receptionId).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectExec("DELETE FROM products").
			WithArgs(receptionId).
			WillReturnError(sql.ErrConnDone)

		mock.ExpectRollback()

		err = productRepo.DeleteLastProductForPVZ(ctx, &pvzId)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrDeleteProduct, err)
	})
}

func TestGetAmountOfProductsForReceptionTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			log.Error(err)
		}
	}(db)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	receptionRepo := NewReceptionRepository(sqlxDB)
	productRepo := NewProductRepository(sqlxDB, receptionRepo)
	receptionId := "1"

	ctx := context.Background()
	t.Run("Successful_Get_Amount_Of_Products", func(t *testing.T) {
		mock.ExpectBegin()
		tx, err := sqlxDB.Beginx()
		assert.NoError(t, err)

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM products").
			WithArgs(receptionId).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		count, err := productRepo.getAmountOfProductsForReceptionTx(ctx, tx, &receptionId)
		assert.NoError(t, err)
		assert.Equal(t, 5, count)
	})

	t.Run("Failure_DB_Error", func(t *testing.T) {
		mock.ExpectBegin()
		tx, err := sqlxDB.Beginx()
		assert.NoError(t, err)

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM products").
			WithArgs(receptionId).
			WillReturnError(sql.ErrConnDone)

		count, err := productRepo.getAmountOfProductsForReceptionTx(ctx, tx, &receptionId)
		assert.Error(t, err)
		assert.Equal(t, 0, count)
	})
}
