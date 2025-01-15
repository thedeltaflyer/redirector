package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"redirector/helpers"
)

func Test_HandleGet(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		mode           string
		mockGetValue   []byte
		mockGetErr     error
		expectStatus   int
		expectLocation string
		expectJSON     gin.H
		expectText     string
		expectQR       bool
	}{
		{
			name:           "redirect_success",
			key:            "existingKey",
			mode:           "/",
			mockGetValue:   []byte("https://example.com"),
			expectStatus:   http.StatusTemporaryRedirect,
			expectLocation: "https://example.com",
		},
		{
			name:         "json_success",
			key:          "existingKey",
			mode:         "/json",
			mockGetValue: []byte("https://example.com"),
			expectStatus: http.StatusOK,
			expectJSON:   gin.H{"key": "existingKey", "url": "https://example.com"},
		},
		{
			name:         "text_success",
			key:          "existingKey",
			mode:         "/text",
			mockGetValue: []byte("https://example.com"),
			expectStatus: http.StatusOK,
			expectText:   "https://example.com",
		},
		{
			name:         "qr_success",
			key:          "existingKey",
			mode:         "/qr",
			mockGetValue: []byte("https://example.com"),
			expectStatus: http.StatusOK,
			expectQR:     true,
		},
		{
			name:         "qr_bad_param",
			key:          "existingKey",
			mode:         "/qr?size=bad_size",
			mockGetValue: []byte("https://example.com"),
			expectStatus: http.StatusBadRequest,
			expectQR:     false,
		},
		//{
		//	name:         "qr_too_long",
		//	key:          "existingKey",
		//	mode:         "/qr?level=B",
		//	mockGetValue: []byte("https://" + strings.Repeat("A", 100000000000) + ".com"),
		//	expectStatus: http.StatusInternalServerError,
		//	expectQR:     false,
		//},
		{
			name:         "key_not_found",
			key:          "missingKey",
			mode:         "/",
			mockGetValue: nil,
			expectStatus: http.StatusNotFound,
		},
		{
			name:         "invalid_mode",
			key:          "existingKey",
			mode:         "/invalid",
			mockGetValue: []byte("https://example.com"),
			expectStatus: http.StatusNotFound,
		},
		{
			name:         "internal_error",
			key:          "existingKey",
			mode:         "/",
			mockGetErr:   fmt.Errorf("database error"),
			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.Default()
			mockStore := &mockKVWrapper{getFunc: func(key []byte) ([]byte, error) {
				return test.mockGetValue, test.mockGetErr
			}}
			controller := &RedirectorController{KV: mockStore}
			router.GET("/:key/*mode", controller.HandleGet)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s%s", test.key, test.mode), nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, test.expectStatus, rec.Code)
			if test.expectLocation != "" {
				assert.Equal(t, test.expectLocation, rec.Header().Get("Location"))
			}
			if test.expectJSON != nil {
				var response gin.H
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, test.expectJSON, response)
			}
			if test.expectText != "" {
				assert.Equal(t, test.expectText, rec.Body.String())
			}
			if test.expectQR {
				assert.Equal(t, "image/png", rec.Header().Get("Content-Type"))
				assert.Greater(t, len(rec.Body.Bytes()), 0)
			}
		})
	}
}

func Test_HandlePost(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		body           gin.H
		mockPutErr     error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "valid_no_key",
			key:            "",
			body:           gin.H{"url": "https://example.com"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid_with_key",
			key:            "customKey",
			body:           gin.H{"url": "https://example.com"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "key_conflict",
			key:            "existingKey",
			body:           gin.H{"url": "https://example.com"},
			mockPutErr:     helpers.NewAlreadyExistsError([]byte("existingKey")),
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "key_too_long",
			key:            strings.Repeat("A", 200),
			body:           gin.H{"url": "https://example.com"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "bind_error",
			body:           gin.H{"incorrect_field": "https://example.com"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.Default()
			mockStore := &mockKVWrapper{exclusivePutFunc: func(key []byte, value []byte) error {
				return test.mockPutErr
			}}
			controller := &RedirectorController{KV: mockStore}
			router.POST("/:key", controller.HandlePost)
			router.POST("", controller.HandlePost)

			body, _ := json.Marshal(test.body)
			req := httptest.NewRequest(http.MethodPost, "/"+test.key, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, test.expectedStatus, rec.Code)
		})
	}
}

func Test_HandlePutWithKey(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		body           gin.H
		mockRepErr     error
		mockReplace    map[string]string
		expectedStatus int
	}{
		{
			name:           "valid_replace",
			key:            "existingKey",
			body:           gin.H{"url": "https://new-example.com"},
			mockReplace:    map[string]string{"existingKey": "https://example.com"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "key_not_exist",
			key:            "nonexistentKey",
			body:           gin.H{"url": "https://new-example.com"},
			mockRepErr:     helpers.NewDoesNotExistError([]byte("nonexistentKey")),
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "bind_error",
			key:            "existingKey",
			body:           gin.H{"incorrect_field": "https://new-example.com"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.Default()
			mockStore := &mockKVWrapper{replaceFunc: func(key []byte, value []byte) ([]byte, error) {
				return nil, test.mockRepErr
			}}
			controller := &RedirectorController{KV: mockStore}
			router.PUT("/:key", controller.HandlePutWithKey)

			body, _ := json.Marshal(test.body)
			req := httptest.NewRequest(http.MethodPut, "/"+test.key, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, test.expectedStatus, rec.Code)
		})
	}
}
