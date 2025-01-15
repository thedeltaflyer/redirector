package helpers

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

// HexToRGBA converts a hex color string to an RGBA color representation, supporting 3, 4, 6, or 8-character formats.
func HexToRGBA(hex string) (color.RGBA, error) {
	// If the string starts with "#", drop it
	if strings.HasPrefix(hex, "#") {
		hex = hex[1:]
	}

	c := color.RGBA{}

	// Convert the string to a hex int
	value, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return c, fmt.Errorf("invalid hex value: #%s", hex)
	}

	// Different offsets for different hex color formats...
	switch len(hex) {
	case 3: // "#RGB", multiply values by 17 to get the full 0-255 range.
		c.R = uint8(value>>8) * 17
		c.G = uint8(value>>4) & 0xF * 17
		c.B = uint8(value) & 0xF * 17
		c.A = 255
	case 4: // "#RGBA", multiply values by 17 to get the full 0-255 range.
		c.R = uint8(value>>12) * 17
		c.G = uint8(value>>8) & 0xF * 17
		c.B = uint8(value>>4) & 0xF * 17
		c.A = uint8(value) & 0xF * 17
	case 6: // "#RRGGBB", your standard hex color.
		c.R = uint8(value >> 16)
		c.G = uint8(value>>8) & 0xFF
		c.B = uint8(value) & 0xFF
		c.A = 255
	case 8: // "#RRGGBBAA", your standard hex color with alpha.
		c.R = uint8(value >> 24)
		c.G = uint8(value>>16) & 0xFF
		c.B = uint8(value>>8) & 0xFF
		c.A = uint8(value) & 0xFF
	default:
		// We only support 3, 4, 6, and 8 character hex colors
		return c, fmt.Errorf("invalid length for a hex color: #%s (%d)", hex, len(hex))
	}

	return c, nil
}
