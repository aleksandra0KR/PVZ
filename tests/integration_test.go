//go:build integration
// +build integration

package tests

import (
	"bytes"
	"encoding/json"
	"final/internal/controller"
	"final/internal/database"
	"final/internal/domain"
	"final/internal/logger"
	auth "final/internal/midleware"
	"final/internal/repository"
	"final/internal/usecase"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestDB() (http.Handler, *usecase.Usecase, *sqlx.DB) {
	logger.InitLogger()
	if err := godotenv.Load("../test.env"); err != nil {
		log.Fatalf("error loading .env file")
	}

	db := database.InitializeDBPostgres(3, 10)
	err := auth.InitAuthFromConfig()
	if err != nil {
		log.Error(err)
	}

	repository := repository.NewRepository(db.GetDB())
	usecase := usecase.NewUseCase(repository)
	handler := controller.NewHandler(*usecase)
	router := handler.Handle()

	return router, usecase, db.GetDB()
}

func clearDatabase(db *sqlx.DB) error {
	_, err := db.Exec("DELETE FROM products")
	if err != nil {
		log.Error("error clearing database")
		return err
	}
	_, err = db.Exec("DELETE FROM receptions")
	if err != nil {
		log.Error("error clearing database")
		return err
	}
	_, err = db.Exec("DELETE FROM pvz")
	if err != nil {
		log.Error("error clearing database")
		return err
	}
	return nil
}

func getToken(t *testing.T, router http.Handler, role string) string {
	user := &domain.User{
		Role: &role,
	}
	userJSON, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var token string
	err := json.Unmarshal(w.Body.Bytes(), &token)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	return token
}

func TestCreatePVZ_AddReception_AddProducts_CloseReception(t *testing.T) {
	router, _, db := setupTestDB()
	defer func(db *sqlx.DB) {
		err := clearDatabase(db)
		if err != nil {
			log.Error(err)
		}
	}(db)

	// Get moderator and employee token
	tokenModerator := getToken(t, router, "moderator")
	tokenEmployee := getToken(t, router, "employee")

	// Create Pvz
	allowedCity := "Москва"
	newPVZ := &domain.PVZ{
		City: &allowedCity,
	}
	pvzJSON, _ := json.Marshal(newPVZ)

	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBuffer(pvzJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenModerator)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdPVZ domain.PVZ
	err := json.Unmarshal(w.Body.Bytes(), &createdPVZ)
	assert.NoError(t, err)
	assert.Equal(t, newPVZ.City, createdPVZ.City)

	// Create reception
	status := "in_progress"
	reception := &domain.Reception{
		PvzId:  createdPVZ.ID,
		Status: &status,
	}
	receptionJSON, _ := json.Marshal(reception)

	req = httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewBuffer(receptionJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenEmployee)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdReception domain.Reception
	err = json.Unmarshal(w.Body.Bytes(), &createdReception)
	assert.NoError(t, err)
	assert.Equal(t, reception.Status, createdReception.Status)

	// Create 50 products
	typeAllowed := "электроника"
	for i := 0; i < 50; i++ {
		product := &domain.InputProduct{
			PvzId: createdReception.PvzId,
			Type:  &typeAllowed,
		}
		productJSON, _ := json.Marshal(product)

		req = httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(productJSON))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tokenEmployee)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var createdProduct domain.Product
		err = json.Unmarshal(w.Body.Bytes(), &createdProduct)
		assert.NoError(t, err)
		assert.Equal(t, product.Type, createdProduct.Type)
	}

	// Close reception
	closeReceptionJSON, _ := json.Marshal(&domain.Reception{ID: createdReception.ID})

	req = httptest.NewRequest(http.MethodPost, "/pvz/"+*createdPVZ.ID+"/close_last_reception", bytes.NewBuffer(closeReceptionJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenEmployee)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var closedReception domain.Reception
	err = json.Unmarshal(w.Body.Bytes(), &closedReception)
	assert.NoError(t, err)
	assert.Equal(t, "close", *closedReception.Status)
}
