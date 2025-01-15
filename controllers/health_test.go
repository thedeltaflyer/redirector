package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockKVWrapper struct {
	getFunc          func(key []byte) ([]byte, error)
	putFunc          func(key []byte, value []byte) error
	exclusivePutFunc func(key []byte, value []byte) error
	replaceFunc      func(key []byte, value []byte) ([]byte, error)
	deleteFunc       func(key []byte) error
	forceError       error
}

func (m *mockKVWrapper) Get(key []byte) ([]byte, error) {
	return m.getFunc(key)
}
func (m *mockKVWrapper) Put(key []byte, value []byte) error {
	return m.putFunc(key, value)
}
func (m *mockKVWrapper) ExclusivePut(key []byte, value []byte) error {
	return m.exclusivePutFunc(key, value)
}
func (m *mockKVWrapper) Replace(key []byte, value []byte) ([]byte, error) {
	return m.replaceFunc(key, value)
}
func (m *mockKVWrapper) Delete(key []byte) error {
	return m.deleteFunc(key)
}

func TestHealthController_HandleGet(t *testing.T) {
	tests := []struct {
		name           string
		kvGetFunc      func(key []byte) ([]byte, error)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "healthy",
			kvGetFunc: func(key []byte) ([]byte, error) {
				if string(key) == "health" {
					return []byte("ok"), nil
				}
				return nil, errors.New("invalid key")
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"ok"}`,
		},
		{
			name: "kv error",
			kvGetFunc: func(key []byte) ([]byte, error) {
				return nil, errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"database error","status":"error"}`,
		},
		{
			name: "incorrect health value",
			kvGetFunc: func(key []byte) ([]byte, error) {
				return []byte("not-ok"), nil
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"expected \"ok\", got: \"not-ok\"","status":"error"}`,
		},
	}

	gin.SetMode(gin.TestMode)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockKV := &mockKVWrapper{getFunc: tt.kvGetFunc}
			controller := &HealthController{KV: mockKV}
			router := gin.New()
			router.GET("/health", controller.HandleGet)

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()

			// Act
			router.ServeHTTP(rec, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.True(t, strings.Contains(rec.Body.String(), tt.expectedBody))
		})
	}
}
