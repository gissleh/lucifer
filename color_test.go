package lucifer_test

import (
	"fmt"
	"github.com/gissleh/lucifer"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestColor_FromString(t *testing.T) {
	table := map[string]struct {
		r string
		e bool
	}{
		"#FF00FF":                {r: "ff00ff"},
		"#00ff00":                {r: "00ff00"},
		"00FF00":                 {r: "00ff00"},
		"#deaded":                {r: "deaded"},
		"#123456":                {r: "123456"},
		"78c28a":                 {r: "78c28a"},
		"rgb(255, 255, 0)":       {r: "ffff00"},
		"rgb(0, 255, 0)":         {r: "00ff00"},
		"rgb(0, 0, 255)":         {r: "0000ff"},
		"rgb(184, 255, 160)":     {r: "b8ffa0"},
		"rgb(255, 160, 184)":     {r: "ffa0b8"},
		"hsv(0, 1, 1)":           {r: "ff0000"},
		"hsv(344,0.221,0.871)":   {r: "deadba"},
		"rgb(255, 160, 184, 0)":  {e: true},
		"rgb(256, 255, 255)":     {e: true},
		"rgb(255, 256, 255)":     {e: true},
		"rgb(255, 255, 256)":     {e: true},
		"rgba(255, 160, 184, 0)": {e: true},
		"rgb(255, 160)":          {e: true},
		"rgb(255, 160,)":         {e: true},
		"rgb(two, one, one)":     {e: true},
		"hsv(360,0.221,0.871)":   {e: true},
		"hsv(359,1.221,0.871)":   {e: true},
		"hsv(359,0.221,2.871)":   {e: true},
		"abc":                    {r: "aabbcc"},
		"abcd":                   {e: true},
		"6000":                   {e: true},
		"#aaccff00":              {e: true},
		"#aabbss":                {e: true},
		"fuck":                   {e: true},
		"500k":                   {r: "ff3800"},
		"1000k":                  {r: "ff3800"},
		"6000K":                  {r: "fff3ef"},
		"6034k":                  {r: "fff3f0"},
		"6067k":                  {r: "fff4f1"},
		"12000k":                 {r: "c3d1ff"},
		"10000000k":              {r: "c3d1ff"},
	}

	for input, output := range table {
		t.Run(input, func(t *testing.T) {
			c := lucifer.Color{}
			err := c.SetFromString(input)
			if err != nil {
				if output.e {
					return
				}

				t.Error(input, err)
				return
			} else if output.e {
				t.Error(input, "error expected, but got none")
				return
			}

			if c.Hex() != output.r {
				t.Error(input, "mismatch")
				t.Error(input, "expected:", output)
				t.Error(input, "actual  :", c.Hex())
			}
		})
	}
}

func TestColor_HSV(t *testing.T) {
	table := map[string]struct {
		h, s, v float64
	}{
		"#ff0000": {h: 0, s: 1, v: 1},
		"#000000": {h: 0, s: 0, v: 0},
		"#ffcc11": {h: 47, s: 0.933, v: 1},
		"#11ffcc": {h: 167, s: 0.933, v: 1},
		"#abcdef": {h: 210, s: 0.285, v: 0.937},
		"#deadba": {h: 344, s: 0.221, v: 0.871},
	}

	for input, output := range table {
		t.Run(input, func(t *testing.T) {
			c := lucifer.Color{}
			assert.NoError(t, c.SetFromString(input))

			h, s, v := c.HSV()
			assert.Equal(t, output.h, math.Round(h), "hue")
			assert.Equal(t, output.s, math.Round(s*1000)/1000, "sat")
			assert.Equal(t, output.v, math.Round(v*1000)/1000, "value")
		})
	}
}

func TestColor_SetHSV(t *testing.T) {
	table := []struct {
		h, s, v float64
	}{
		{h: 0, s: 1, v: 1},
		{h: 0, s: 0, v: 0},
		{h: 47, s: 0.933, v: 1},
		{h: 167, s: 0.933, v: 1},
		{h: 210, s: 0.285, v: 0.937},
		{h: 344, s: 0.221, v: 0.871},
	}

	for _, color := range table {
		t.Run(fmt.Sprintf("%.0f,%.3f,%.3f", color.h, color.s, color.v), func(t *testing.T) {
			c := lucifer.Color{}
			c.SetHSV(color.h, color.s, color.v)

			h, s, v := c.HSV()
			assert.Equal(t, color.h, math.Round(h), "hue")
			assert.Equal(t, color.s, math.Round(s*1000)/1000, "sat")
			assert.Equal(t, color.v, math.Round(v*1000)/1000, "value")
		})
	}
}
