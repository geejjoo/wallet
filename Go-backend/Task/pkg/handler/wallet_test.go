package handler

import (
	"bytes"
	"fmt"
	"github.com/geejjoo/task"
	"github.com/geejjoo/task/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockService struct {
}

func (m *mockService) CreateWallet() (string, float64, error) {
	return "test", 100, nil
}

func (m *mockService) GetWallet(id string) (string, float64, error) {
	switch id {
	case "1":
		return "test", 100, nil
	case "":
		return "", 0, fmt.Errorf("Invalid wallet ID")
	default:
		return "test", 100, nil
	}
}

func (m *mockService) UpdateWallet(wallet *task.UpdateWallet) error {
	if wallet.Amount <= 0 {
		return fmt.Errorf("amount should be positive")
	}
	if wallet.FromID == wallet.ToID {
		return fmt.Errorf("Send wallet and target wallet are the same")
	}
	if wallet.Amount == 555 {
		return service.ToIdError
	}
	if wallet.Amount > 100 {
		return service.BalanceError
	}
	if wallet.FromID == "invalid_id" {
		return service.FromIdError
	}

	return nil
}

func (m *mockService) GetHistory(id string) ([]task.History, error) {
	if id == "nonexistent" {
		return []task.History{}, service.WalletNotFoundError
	}
	if id == "1" {
		return []task.History{}, service.DatabaseError
	}

	return []task.History{}, nil
}

func TestCreateWallet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := &Handler{
		services: &service.Service{&mockService{}},
	}

	r := gin.New()
	r.POST("/wallet/", handler.createWallet)

	t.Run("Valid request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/wallet/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"id":"test"`)
		assert.Contains(t, w.Body.String(), `"balance":100`)
	})
}

func TestGetWallet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := &Handler{
		services: &service.Service{&mockService{}},
	}

	r := gin.New()
	r.GET("/wallet/:id", handler.getWallet)

	t.Run("Valid request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/wallet/123", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"id":"test"`)
		assert.Contains(t, w.Body.String(), `"balance":100`)
	})

	t.Run("Invalid request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/wallet/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

}

func TestUpdateWallet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := &Handler{
		services: &service.Service{&mockService{}},
	}

	r := gin.New()
	r.POST("/wallet/:id/send", handler.updateWallet)

	t.Run("Valid request", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := []byte(`{"to_id":"456", "amount":50}`)
		req, _ := http.NewRequest("POST", "/wallet/123/send", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid request - Incorrect wallet ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := []byte(`{"to_id":"456", "amount":50}`)
		req, _ := http.NewRequest("POST", "/wallet/invalid_id/send", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Incorrect wallet ID")
	})

	t.Run("Invalid request - Incorrect target ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := []byte(`{"to_id":"1", "amount":555}`)
		req, _ := http.NewRequest("POST", "/wallet/1234/send", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Incorrect target ID")
	})

	t.Run("Invalid request - Not enough money", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := []byte(`{"to_id":"456", "amount":150}`)
		req, _ := http.NewRequest("POST", "/wallet/12345/send", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Not enough money")
	})

	t.Run("Invalid request - negative amount", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := []byte(`{"to_id":"456", "amount":-1}`)
		req, _ := http.NewRequest("POST", "/wallet/123456/send", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Amount should be positive")
	})
}

func TestGetHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := &Handler{
		services: &service.Service{&mockService{}},
	}

	r := gin.New()
	r.GET("/wallet/:id/history", handler.getHistory)

	t.Run("Valid request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/wallet/123/history", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

	})

	t.Run("Wallet not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/wallet/nonexistent/history", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Internal server error", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/wallet/1/history", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

}
