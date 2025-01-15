package logging

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetLogger(t *testing.T) {
	tests := []struct {
		name        string
		setupLogger func()
		expectPanic bool
	}{
		{
			name: "validLogger",
			setupLogger: func() {
				logger = logrus.New()
			},
			expectPanic: false,
		},
		{
			name: "loggerNotInitialized",
			setupLogger: func() {
				logger = nil
			},
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupLogger()

			if tt.expectPanic {
				assert.PanicsWithError(t, "logger not initialized", func() {
					GetLogger()
				})
			} else {
				assert.NotPanics(t, func() {
					l := GetLogger()
					assert.NotNil(t, l)
				})
			}
		})
	}
}
