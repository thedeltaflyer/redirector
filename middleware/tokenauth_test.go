package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockKV struct {
	data  map[string][]byte
	error error
}

func (m *mockKV) Get(key []byte) ([]byte, error) {
	if m.error != nil {
		return nil, m.error
	}
	value, exists := m.data[string(key)]
	if !exists {
		return nil, nil
	}
	return value, nil
}

func (m *mockKV) Put(key []byte, value []byte) error               { return nil }
func (m *mockKV) ExclusivePut(key []byte, value []byte) error      { return nil }
func (m *mockKV) Replace(key []byte, value []byte) ([]byte, error) { return nil, nil }
func (m *mockKV) Delete(key []byte) error                          { return nil }

func TestTokenAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		kvMock         *mockKV
		expectedStatus int
	}{
		{
			name:           "No Authorization header",
			authHeader:     "",
			kvMock:         &mockKV{},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Authorization header format",
			authHeader:     "InvalidHeader",
			kvMock:         &mockKV{},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid scheme in Authorization header",
			authHeader:     "Basic token",
			kvMock:         &mockKV{},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Empty token",
			authHeader:     "Bearer ",
			kvMock:         &mockKV{},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Token not found in KV",
			authHeader:     "Bearer invalid_token",
			kvMock:         &mockKV{data: map[string][]byte{}},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "KV Get method returns error",
			authHeader:     "Bearer valid_token",
			kvMock:         &mockKV{error: errors.New("internal error")},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Valid token",
			authHeader:     "Bearer valid_token",
			kvMock:         &mockKV{data: map[string][]byte{"valid_token": []byte("data")}},
			expectedStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.Use(TokenAuthMiddleware(test.kvMock))
			r.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if test.authHeader != "" {
				req.Header.Set("Authorization", test.authHeader)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != test.expectedStatus {
				t.Errorf("expected status %d, got %d", test.expectedStatus, w.Code)
			}
		})
	}
}
