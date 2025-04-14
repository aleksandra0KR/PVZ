package controller

import (
	"final/internal/usecase"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	services := usecase.Usecase{}
	handler := NewHandler(services)

	t.Run("Failure_No_Route", func(t *testing.T) {
		router := handler.Handle()

		req, _ := http.NewRequest("GET", "/nonexistent", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotImplemented, w.Code)
		assert.JSONEq(t, `{"message":"not implemented"}`, w.Body.String())
	})

	t.Run("Create_Pvz", func(t *testing.T) {
		router := handler.Handle()

		req, _ := http.NewRequest("POST", "/pvz", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Get_Pvz_List", func(t *testing.T) {
		router := handler.Handle()

		req, _ := http.NewRequest("GET", "/pvz", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Close_Last_Reception", func(t *testing.T) {
		router := handler.Handle()

		req, _ := http.NewRequest("POST", "/pvz/123/close_last_reception", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Delete_Last_Product_For_Pvz", func(t *testing.T) {
		router := handler.Handle()

		req, _ := http.NewRequest("POST", "/pvz/123/delete_last_product", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Create_Reception", func(t *testing.T) {
		router := handler.Handle()

		req, _ := http.NewRequest("POST", "/receptions", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Create_Product", func(t *testing.T) {
		router := handler.Handle()

		req, _ := http.NewRequest("POST", "/products", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
