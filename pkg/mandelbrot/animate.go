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
	Width  int
	Height int
	FPS    int
	Path   []AnimationConfigPathElement
}

type AnimationConfigPathElement struct {
	Zoom        float64
	TargetX     float64
	TargetY     float64
	RawDuration string        `json:"Duration"`
	Duration    time.Duration `json:"-"`
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
	if len(cfg.Path) < 2 {
		return nil, fmt.Errorf("at least two path elements are required")
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
			Zoom:    zoomInterp.At(currentFrame),
			TargetX: xInterp.At(currentFrame),
			TargetY: yInterp.At(currentFrame),
		}
		if err := Render(renderCfg, img, palette); err != nil {
			return fmt.Errorf("rendering (frame %d/%d): %w", currentFrame, frameCount, err)
		}
		g.Image = append(g.Image, img)
		g.Delay = append(g.Delay, frameDelay)
	}
	return nil
}
