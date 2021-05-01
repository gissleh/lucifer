package lucifer

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gissleh/lucifer/internal/constants"
	"github.com/lucasb-eyer/go-colorful"
	"math"
	"strconv"
	"strings"
)

var ErrInvalidInput = errors.New("invalid input")

type Color struct {
	R float64
	G float64
	B float64
	K int
}

// ParseColor parses the color.
func ParseColor(str string) (Color, error) {
	color := Color{}
	err := color.SetFromString(str)

	return color, err
}

// MustParseColor parses a color.
func MustParseColor(str string) Color {
	color, err := ParseColor(str)
	if err != nil {
		panic(err)
	}

	return color
}

// ColorFromHSV creates a color from HSV values.
func ColorFromHSV(h, s, v float64) Color {
	color := Color{}
	color.SetHSV(h, s, v)

	return color
}

// FullBright gets the color with value 1.
func (c *Color) FullBright() {
	h, s, _ := c.HSV()
	c.SetHSV(h, s, 1)
}

// HSV gets the hue/sat/value.
func (c *Color) HSV() (h float64, s float64, v float64) {
	return (colorful.Color{R: c.R, G: c.G, B: c.B}).Hsv()
}

func (c *Color) SetHSV(h, s, v float64) {
	cc := colorful.Hsv(h, s, v)
	c.R = cc.R
	c.G = cc.G
	c.B = cc.B
	c.K = 0
}

func (c *Color) Hex() string {
	var data [3]byte
	data[0] = byte(c.R * 255)
	data[1] = byte(c.G * 255)
	data[2] = byte(c.B * 255)

	return hex.EncodeToString(data[:])
}

func (c *Color) SetHex(hexStr string) error {
	if hexStr[0] == '#' {
		hexStr = hexStr[1:]
	}

	// ABC -> AABBCC
	if len(hexStr) == 3 {
		sb := strings.Builder{}
		sb.Grow(6)
		for _, r := range hexStr {
			sb.WriteRune(r)
			sb.WriteRune(r)
		}
		hexStr = sb.String()
	}

	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return err
	}
	if len(data) != 3 {
		return ErrInvalidInput
	}

	c.R = float64(data[0]) / 255
	c.G = float64(data[1]) / 255
	c.B = float64(data[2]) / 255
	c.K = 0

	return nil
}

func (c *Color) SetKelvin(kelvin int) {
	if kelvin < 1000 {
		kelvin = 1000
	} else if kelvin > 12000 {
		kelvin = 12000
	}

	var values [3]int

	for i, k := range constants.KelvinTable {
		if kelvin > k.Temp {
			continue
		}
		if k.Temp == kelvin || i == 0 {
			values = k.RGB
			break
		}

		k2 := constants.KelvinTable[i-1]
		diff := kelvin - k.Temp
		fac := float64(diff) / float64(k2.Temp-k.Temp)

		for i := 0; i < 3; i++ {
			values[i] = int(math.Floor(float64(k2.RGB[i])*fac)) + int(math.Ceil(float64(k.RGB[i])*(1-fac)))
		}

		break
	}

	c.R = float64(values[0]) / 255
	c.G = float64(values[1]) / 255
	c.B = float64(values[2]) / 255
	c.K = kelvin
}

func (c *Color) SetFromString(s string) error {
	isRgbFunc := strings.HasPrefix(s, "rgb(")
	isHsvFunc := strings.HasPrefix(s, "hsv(")
	if (isRgbFunc || isHsvFunc) && strings.HasSuffix(s, ")") {
		argsStr := strings.SplitN(s[4:len(s)-1], ",", 3)
		args := make([]float64, len(argsStr))
		if len(args) < 3 {
			return errors.New("invalid rgb(...)")
		}

		for i, argStr := range argsStr {
			arg, err := strconv.ParseFloat(strings.Trim(argStr, " \t\r\n "), 64)
			if err != nil {
				return err
			}

			args[i] = arg
		}

		if isRgbFunc {
			if args[0] > 255 {
				return fmt.Errorf("invalid red: %f (0..255)", args[0])
			}
			if args[1] > 255 {
				return fmt.Errorf("invalid green: %f (0..255)", args[0])
			}
			if args[2] > 255 {
				return fmt.Errorf("invalid blue: %f (0..255)", args[0])
			}

			c.R = args[0] / 255
			c.G = args[1] / 255
			c.B = args[2] / 255

			return nil
		} else if isHsvFunc {
			if args[0] >= 360 {
				return fmt.Errorf("invalid hue: %f (0..360)", args[0])
			}
			if args[1] > 1 {
				return fmt.Errorf("invalid sat: %f (0..1)", args[0])
			}
			if args[2] > 1 {
				return fmt.Errorf("invalid value: %f (0..1)", args[0])
			}

			c.SetHSV(args[0], args[1], args[2])

			return nil
		}
	}

	if strings.HasSuffix(s, "k") || strings.HasSuffix(s, "K") {
		n, err := strconv.Atoi(s[:len(s)-1])
		if err != nil {
			return err
		}

		c.SetKelvin(n)
		return nil
	}

	return c.SetHex(s)
}

func (c *Color) MarshalJSON() ([]byte, error) {
	return []byte(c.Hex()), nil
}

func (c *Color) UnmarshalJSON(v []byte) error {
	var s string

	err := json.Unmarshal(v, &s)
	if err != nil {
		return err
	}

	return c.SetFromString(s)
}
