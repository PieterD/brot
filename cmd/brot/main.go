package main

import (
	"fmt"
	"github.com/PieterD/brot/pkg/mandelbrot"
	"image/color"
	"image/gif"
	"os"
)

func main() {
	cfg, ok := NewConfigFromFlags()
	if !ok {
		os.Exit(1)
	}
	if err := run(cfg); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "running brot: %v", err)
		os.Exit(1)
	}
}

func run(cfg Config) error {
	const paletteSize = 256
	animationConfig, err := mandelbrot.NewAnimateConfigFromFile(cfg.ConfigFile)
	if err != nil {
		return fmt.Errorf("extracting animation config: %w", err)
	}
	blue := color.RGBA{
		B: 255,
		A: 255,
	}
	palette := mandelbrot.Gradient(blue, color.Black, 256)
	g, err := mandelbrot.Animate(animationConfig, palette)
	if err != nil {
		return fmt.Errorf("animating: %w", err)
	}
	out, err := os.Create(cfg.OutputFile)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer func() { _ = out.Close() }()
	if err := gif.EncodeAll(out, g); err != nil {
		return fmt.Errorf("encoding GIF: %w", err)
	}
	return nil
}
