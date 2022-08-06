package mandelbrot

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"os"
	"time"
)

type AnimationConfig struct {
	Width         int
	Height        int
	FPS           int
	MaxIterations int
	Path          []AnimationConfigPathElement
}

type AnimationConfigPathElement struct {
	Zoom        float64
	TargetX     float64
	TargetY     float64
	RawDuration string        `json:"Duration"`
	Duration    time.Duration `json:"-"`
}

func (cfg AnimationConfig) Validate() error {
	if cfg.FPS <= 0 {
		return fmt.Errorf("invalid FPS (%d)", cfg.FPS)
	}
	if cfg.MaxIterations <= 0 {
		return fmt.Errorf("invalid MaxIterations (%d)", cfg.MaxIterations)
	}
	if cfg.Width <= 0 {
		return fmt.Errorf("invalid Width (%d)", cfg.Width)
	}
	if cfg.Height <= 0 {
		return fmt.Errorf("invalid Height (%d)", cfg.Height)
	}
	if len(cfg.Path) == 0 {
		return fmt.Errorf("at least one Path element is required")
	}
	for i, pe := range cfg.Path {
		if err := pe.Validate(i == 0); err != nil {
			return fmt.Errorf("path element (%d): %w", i, err)
		}
	}
	return nil
}

func (cfg AnimationConfigPathElement) Validate(first bool) error {
	switch {
	case first && cfg.Duration == 0:
	case first && cfg.Duration != 0:
		return fmt.Errorf("first Path element cannot have a Duration")
	case !first && cfg.Duration == 0:
		return fmt.Errorf("non-first Path element must have a Duration")
	case !first && cfg.Duration != 0:
	}
	if cfg.TargetX < 0 || cfg.TargetX > 1 || cfg.TargetY < 0 || cfg.TargetY > 1 {
		return fmt.Errorf("target (%f,%f) is out of bounds", cfg.TargetX, cfg.TargetY)
	}
	if cfg.Zoom < 1 {
		return fmt.Errorf("zoom (%f) cannot be below 1", cfg.Zoom)
	}
	return nil
}

func NewAnimateConfigFromFile(filePath string) (AnimationConfig, error) {
	var v AnimationConfig
	f, err := os.Open(filePath)
	if err != nil {
		return AnimationConfig{}, fmt.Errorf("opening file: %w", err)
	}
	defer func() { _ = f.Close() }()
	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&v); err != nil {
		return AnimationConfig{}, fmt.Errorf("decoding JSON: %w", err)
	}
	for i := range v.Path {
		pathElement := &v.Path[i]
		if pathElement.RawDuration == "" {
			continue
		}
		dur, err := time.ParseDuration(pathElement.RawDuration)
		if err != nil {
			return AnimationConfig{}, fmt.Errorf("invalid duration (%s): %w", pathElement.RawDuration, err)
		}
		pathElement.Duration = dur
	}
	return v, nil
}

func Animate(cfg AnimationConfig, palette color.Palette) (*gif.GIF, error) {
	g := &gif.GIF{
		Image:     nil,
		Delay:     nil,
		LoopCount: 0,
		Disposal:  nil,
		Config: image.Config{
			ColorModel: palette,
			Width:      cfg.Width,
			Height:     cfg.Height,
		},
		BackgroundIndex: 0,
	}
	if len(cfg.Path) < 1 {
		return nil, fmt.Errorf("at least one path elements are required")
	}
	if len(cfg.Path) == 1 {
		renderCfg := RenderConfig{
			MaxIterations: cfg.MaxIterations,
			Zoom:          cfg.Path[0].Zoom,
			TargetX:       cfg.Path[0].TargetX,
			TargetY:       cfg.Path[0].TargetY,
		}
		imgRect := image.Rect(0, 0, cfg.Width, cfg.Height)
		img := image.NewPaletted(imgRect, palette)
		if err := Render(renderCfg, img, palette); err != nil {
			return nil, fmt.Errorf("rendering single frame: %w", err)
		}

	}
	for i := 1; i < len(cfg.Path); i++ {
		from := cfg.Path[i-1]
		to := cfg.Path[i]
		if err := animatePathLink(cfg, i == 1, palette, g, from, to); err != nil {
			return nil, fmt.Errorf("animating path link (%d): %w", i, err)
		}
	}
	return g, nil
}

func animatePathLink(cfg AnimationConfig, includeFirst bool, palette color.Palette, g *gif.GIF, from, to AnimationConfigPathElement) error {
	fmt.Printf("link: from %+v to %+v\n", from, to)
	frameDuration := time.Second / time.Duration(cfg.FPS)
	frameDelay := int(frameDuration.Seconds() * 100)
	frameCount := int(to.Duration.Seconds() * float64(cfg.FPS))
	zoomInterp := NewInterpolator(from.Zoom, to.Zoom, frameCount)
	xInterp := NewInterpolator(from.TargetX, to.TargetX, frameCount)
	yInterp := NewInterpolator(from.TargetY, to.TargetY, frameCount)
	currentFrame := 0
	if !includeFirst {
		currentFrame++
	}
	for ; currentFrame < frameCount; currentFrame++ {
		fmt.Printf("rendering %d/%d\n", currentFrame+1, frameCount)
		imgRect := image.Rect(0, 0, cfg.Width, cfg.Height)
		img := image.NewPaletted(imgRect, palette)
		renderCfg := RenderConfig{
			MaxIterations: cfg.MaxIterations,
			Zoom:          zoomInterp.At(currentFrame),
			TargetX:       xInterp.At(currentFrame),
			TargetY:       yInterp.At(currentFrame),
		}
		if err := Render(renderCfg, img, palette); err != nil {
			return fmt.Errorf("rendering (frame %d/%d): %w", currentFrame, frameCount, err)
		}
		g.Image = append(g.Image, img)
		g.Delay = append(g.Delay, frameDelay)
	}
	return nil
}
