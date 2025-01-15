package helpers

import (
	"image/color"
	"testing"
)

func TestHexToRGBA(t *testing.T) {
	tests := []struct {
		name    string
		hex     string
		want    color.RGBA
		wantErr bool
	}{
		{
			name:    "valid_short_hex",
			hex:     "#123",
			want:    color.RGBA{R: 17, G: 34, B: 51, A: 255},
			wantErr: false,
		},
		{
			name:    "valid_short_hex_with_alpha",
			hex:     "#1234",
			want:    color.RGBA{R: 17, G: 34, B: 51, A: 68},
			wantErr: false,
		},
		{
			name:    "valid_long_hex",
			hex:     "#112233",
			want:    color.RGBA{R: 17, G: 34, B: 51, A: 255},
			wantErr: false,
		},
		{
			name:    "valid_long_hex_with_alpha",
			hex:     "#11223344",
			want:    color.RGBA{R: 17, G: 34, B: 51, A: 68},
			wantErr: false,
		},
		{
			name:    "invalid_characters",
			hex:     "#12G",
			want:    color.RGBA{},
			wantErr: true,
		},
		{
			name:    "invalid_length",
			hex:     "#12345",
			want:    color.RGBA{},
			wantErr: true,
		},
		{
			name:    "missing_prefix",
			hex:     "123456",
			want:    color.RGBA{R: 18, G: 52, B: 86, A: 255},
			wantErr: false,
		},
		{
			name:    "empty_input",
			hex:     "",
			want:    color.RGBA{},
			wantErr: true,
		},
		{
			name:    "only_prefix",
			hex:     "#",
			want:    color.RGBA{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HexToRGBA(tt.hex)
			if (err != nil) != tt.wantErr {
				t.Errorf("HexToRGBA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HexToRGBA() got = %v, want %v", got, tt.want)
			}
		})
	}
}
