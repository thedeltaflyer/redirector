package helpers

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

// QRConfig defines the configuration for generating a QR code, including size, error recovery level, colors, and border.
type QRConfig struct {
	Size    int
	Level   qrcode.RecoveryLevel
	BgColor color.Color
	FgColor color.Color
	Border  bool
}

// QRParams defines query parameters for configuring a QR code, including size, error correction level, colors, and border.
type QRParams struct {
	Size    int    `form:"size"`
	Level   string `form:"level"`
	BgColor string `form:"bg_color"`
	FgColor string `form:"fg_color"`
	Border  bool   `form:"border"`
}

// GetQRParamsFromContext extracts QR code configuration from the provided gin.Context query parameters.
// It returns a QRConfig struct with parsed values and an error if any parsing issue occurs.
func GetQRParamsFromContext(c *gin.Context) (QRConfig, error) {
	// Create an empty QRConfig.
	conf := QRConfig{}

	// Bind the query parameters, return an error if there's any trouble parsing.
	var params QRParams
	err := c.ShouldBindQuery(&params)
	if err != nil {
		return conf, err
	}

	// Check for sane Size values.
	if params.Size == 0 {
		params.Size = 256
	} else if params.Size > 4096 {
		params.Size = 4096
	} else if params.Size < -164 {
		params.Size = -164 // This is ~4100px with no border
	}
	conf.Size = params.Size

	// Map the QR Levels to L, M, H, or B
	switch strings.ToUpper(params.Level) {
	case "":
		conf.Level = qrcode.Medium
	case "L":
		conf.Level = qrcode.Low
	case "M":
		conf.Level = qrcode.Medium
	case "H":
		conf.Level = qrcode.High
	case "B": // "B" is for "Best" since "H" is already used for "High"
		conf.Level = qrcode.Highest
	default:
		return conf, fmt.Errorf("invalid QR level (must be one of L,M,H,B): %s", params.Level)
	}

	// Parse the background color, fall back to White if not specified
	if params.BgColor == "" {
		conf.BgColor = color.White
	} else {
		if bgRGBA, err := HexToRGBA(params.BgColor); err == nil {
			conf.BgColor = bgRGBA
		} else {
			return conf, err
		}
	}

	// Parse the foreground color, fall back to Black if not specified
	if params.FgColor == "" {
		conf.FgColor = color.Black
	} else {
		if fgRGBA, err := HexToRGBA(params.FgColor); err == nil {
			conf.FgColor = fgRGBA
		} else {
			return conf, err
		}
	}

	// Set if a border is needed
	conf.Border = params.Border

	return conf, nil
}
