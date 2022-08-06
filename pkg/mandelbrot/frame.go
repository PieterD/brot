package mandelbrot

import (
	"fmt"
	"image"
	"image/color"
	"math/cmplx"
)

type RenderConfig struct {
	MaxIterations int
	Zoom          float64
	TargetX       float64
	TargetY       float64
}

func Render(cfg RenderConfig, image *image.Paletted, palette color.Palette) error {
	s := newScaler(image.Bounds().Dx(), image.Bounds().Dy())
	if err := s.Zoom(cfg.Zoom); err != nil {
		return fmt.Errorf("setting zoom (%f): %w", cfg.Zoom, err)
	}
	if err := s.Target(cfg.TargetX, cfg.TargetY); err != nil {
		return fmt.Errorf("targeting: %w", err)
	}
	pg := newPixelGenerator(cfg.MaxIterations)
	for x := 0; x < image.Bounds().Dx(); x++ {
		for y := 0; y < image.Bounds().Dy(); y++ {
			mx, my, err := s.Transform(x, y)
			if err != nil {
				return fmt.Errorf("scaling pixel: %w", err)
			}
			severity := pg.Render(mx, my)
			//fmt.Printf("%3d,%3d %f\n", x, y, severity)
			if severity < 0 || severity > 1.0 {
				return fmt.Errorf("severity (%f) out of bounds", severity)
			}
			idx := int(float64(len(palette)-1) * severity)
			//fmt.Printf("%3d,%3d %d\n", x, y, idx)
			if idx < 0 || idx >= len(palette) {
				return fmt.Errorf("palette index (%d) out of bounds", idx)
			}
			c := palette[idx]
			image.Set(x, y, c)
		}
	}
	return nil
}

type pixelGenerator struct {
	maxIterations int
}

func newPixelGenerator(maxIterations int) *pixelGenerator {
	return &pixelGenerator{
		maxIterations: maxIterations,
	}
}

// Render calculates the "color" of the provided pixel.
// Pixel coordinates are in mandelbrot space.
// The "color" returned is a value between 0 and 1 inclusive,
// scaled to the amount of iterations required to escape.
func (g *pixelGenerator) Render(x, y float64) float64 {
	var z complex128
	c := complex(x, y)
	iteration := 0
	for cmplx.Abs(z) <= 2 && iteration < g.maxIterations {
		z = z*z + c
		iteration++
	}
	return float64(iteration) / float64(g.maxIterations)
}
