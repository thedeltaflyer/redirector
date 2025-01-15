package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandleGet(t *testing.T) {
	tests := []struct {
		name         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "normal_request",
			expectedCode: http.StatusOK,
			expectedBody: "Welcome to lnk.now. This is a private URL shortener for perpetualtag.com.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup Gin context and recorder
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			rc := &RootController{}
			r.GET("/", rc.HandleGet)

			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			// Perform request
			r.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}
