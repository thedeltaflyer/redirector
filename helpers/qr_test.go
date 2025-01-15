package helpers

import (
	"image/color"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

func TestGetQRParamsFromContext(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		wantConfig QRConfig
		wantErr    string
	}{
		{
			name:       "default values",
			query:      "",
			wantConfig: QRConfig{Size: 256, Level: qrcode.Medium, BgColor: color.White, FgColor: color.Black, Border: false},
			wantErr:    "",
		},
		{
			name:       "custom size within range",
			query:      "size=300",
			wantConfig: QRConfig{Size: 300, Level: qrcode.Medium, BgColor: color.White, FgColor: color.Black, Border: false},
			wantErr:    "",
		},
		{
			name:       "size below minimum",
			query:      "size=-200",
			wantConfig: QRConfig{Size: -164, Level: qrcode.Medium, BgColor: color.White, FgColor: color.Black, Border: false},
			wantErr:    "",
		},
		{
			name:       "size above maximum",
			query:      "size=5000",
			wantConfig: QRConfig{Size: 4096, Level: qrcode.Medium, BgColor: color.White, FgColor: color.Black, Border: false},
			wantErr:    "",
		},
		{
			name:       "level set to Low",
			query:      "level=L",
			wantConfig: QRConfig{Size: 256, Level: qrcode.Low, BgColor: color.White, FgColor: color.Black, Border: false},
			wantErr:    "",
		},
		{
			name:       "invalid level",
			query:      "level=Z",
			wantConfig: QRConfig{Size: 256, Level: qrcode.Low, BgColor: nil, FgColor: nil, Border: false},
			wantErr:    "invalid QR level (must be one of L,M,H,B): Z",
		},
		{
			name:       "custom background color",
			query:      "bg_color=#ff0000",
			wantConfig: QRConfig{Size: 256, Level: qrcode.Medium, BgColor: color.RGBA{R: 255, G: 0, B: 0, A: 255}, FgColor: color.Black, Border: false},
			wantErr:    "",
		},
		{
			name:       "invalid background color",
			query:      "bg_color=invalidcolor",
			wantConfig: QRConfig{Size: 256, Level: qrcode.Medium, BgColor: nil, FgColor: nil, Border: false},
			wantErr:    "invalid hex value: #invalidcolor",
		},
		{
			name:       "custom foreground color",
			query:      "fg_color=#00ff00",
			wantConfig: QRConfig{Size: 256, Level: qrcode.Medium, BgColor: color.White, FgColor: color.RGBA{R: 0, G: 255, B: 0, A: 255}, Border: false},
			wantErr:    "",
		},
		{
			name:       "invalid foreground color",
			query:      "fg_color=#12345g",
			wantConfig: QRConfig{Size: 256, Level: qrcode.Medium, BgColor: color.White, FgColor: nil, Border: false},
			wantErr:    "invalid hex value: #12345g",
		},
		{
			name:       "border enabled",
			query:      "border=true",
			wantConfig: QRConfig{Size: 256, Level: qrcode.Medium, BgColor: color.White, FgColor: color.Black, Border: true},
			wantErr:    "",
		},
		{
			name:       "multiple parameters",
			query:      "size=512&level=H&fg_color=#0000ff&bg_color=#ffffff&border=true",
			wantConfig: QRConfig{Size: 512, Level: qrcode.High, BgColor: color.RGBA{R: 255, G: 255, B: 255, A: 255}, FgColor: color.RGBA{R: 0, G: 0, B: 255, A: 255}, Border: true},
			wantErr:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/?"+tt.query, nil)

			gotConfig, err := GetQRParamsFromContext(c)
			if (err != nil && tt.wantErr == "") || (err == nil && tt.wantErr != "") || (err != nil && !strings.Contains(err.Error(), tt.wantErr)) {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			if gotConfig != tt.wantConfig {
				t.Errorf("expected config: %+v, got: %+v", tt.wantConfig, gotConfig)
			}
		})
	}
}
