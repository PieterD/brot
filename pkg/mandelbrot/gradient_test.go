package mandelbrot

import (
	"github.com/stretchr/testify/require"
	"image/color"
	"testing"
)

func TestGradient(t *testing.T) {
	black := color.RGBA{
		A: 255,
	}
	gray := color.RGBA{
		R: 127,
		G: 127,
		B: 127,
		A: 255,
	}
	white := color.RGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}
	var tests = []struct {
		desc   string
		from   color.RGBA
		to     color.RGBA
		steps  int
		result color.Palette
	}{
		{
			desc:  "single",
			from:  black,
			to:    white,
			steps: 1,
			result: color.Palette{
				white,
			},
		},
		{
			desc:  "double",
			from:  black,
			to:    white,
			steps: 2,
			result: color.Palette{
				black,
				white,
			},
		},
		{
			desc:  "triple",
			from:  black,
			to:    white,
			steps: 3,
			result: color.Palette{
				black,
				gray,
				white,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			palette := Gradient(test.from, test.to, test.steps)
			require.Equal(t, test.result, palette)
		})
	}
}
