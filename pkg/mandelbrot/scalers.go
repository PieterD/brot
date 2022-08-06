package mandelbrot

import "fmt"

type normalizingScaler struct {
	width  int
	height int
}

func newNormalizingScaler(width, height int) *normalizingScaler {
	return &normalizingScaler{
		width:  width,
		height: height,
	}
}

func (s *normalizingScaler) Transform(x, y int) (float64, float64, error) {
	if x < 0 || x >= s.width || y < 0 || y >= s.height {
		return 0, 0, fmt.Errorf("coordinates (%d,%d) out of bounds (%d,%d)", x, y, s.width, s.height)
	}
	fracX := float64(x) / float64(s.width-1)
	fracY := float64(y) / float64(s.height-1)
	if fracX < 0.0 || fracX > 1.0 || fracY < 0.0 || fracY > 1.0 {
		return 0, 0, fmt.Errorf("normalized coordinates out of bound (%f,%f)", fracX, fracY)
	}
	return fracX, fracY, nil
}

type zoomingScaler struct {
	source  *normalizingScaler
	zoom    float64
	targetX float64
	targetY float64
}

func newZoomingScaler() *zoomingScaler {
	return &zoomingScaler{
		zoom:    1.0,
		targetX: 0.5,
		targetY: 0.5,
	}
}

func (s *zoomingScaler) Zoom(z float64) error {
	if z < 1.0 {
		return fmt.Errorf("zoom level less than 1 (%f) not allowed", z)
	}
	s.zoom = z
	return nil
}

func (s *zoomingScaler) Target(x, y float64) error {
	if x < 0 || x > 1 || y < 0 || y > 1 {
		return fmt.Errorf("target (%f,%f) out of bounds", x, y)
	}
	s.targetX = x
	s.targetY = y
	return nil
}

func (s *zoomingScaler) Transform(x, y float64) (zx float64, zy float64, err error) {
	zoomedX := (x-s.targetX)/s.zoom + s.targetX
	zoomedY := (y-s.targetY)/s.zoom + s.targetY
	return zoomedX, zoomedY, nil
}

type scaler struct {
	source *normalizingScaler
	zoom   *zoomingScaler
}

func newScaler(width, height int) *scaler {
	return &scaler{
		source: newNormalizingScaler(width, height),
		zoom:   newZoomingScaler(),
	}
}

func (s *scaler) Zoom(z float64) error {
	return s.zoom.Zoom(z)
}

func (s *scaler) Target(x, y float64) error {
	return s.zoom.Target(x, y)
}

func (s *scaler) Transform(x, y int) (sx float64, sy float64, err error) {
	normX, normY, err := s.source.Transform(x, y)
	if err != nil {
		return 0, 0, fmt.Errorf("normalizing scale: %w", err)
	}
	zoomedX, zoomedY, err := s.zoom.Transform(normX, normY)
	if err != nil {
		return 0, 0, fmt.Errorf("zooming scale: %w", err)
	}
	sx = zoomedX*(maxX-minX) + minX
	sy = zoomedY*(maxY-minY) + minY
	return sx, sy, nil
}
