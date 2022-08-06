package mandelbrot

import (
	"image/color"
)

func Gradient(from color.Color, to color.Color, size int) color.Palette {
	fr, fg, fb, fa := from.RGBA()
	tr, tg, tb, ta := to.RGBA()
	interpolators := []*Interpolator{
		NewInterpolator(float64(fr)/256.0, float64(tr)/256.0, size),
		NewInterpolator(float64(fg)/256.0, float64(tg)/256.0, size),
		NewInterpolator(float64(fb)/256.0, float64(tb)/256.0, size),
		NewInterpolator(float64(fa)/256.0, float64(ta)/256.0, size),
	}
	p := make(color.Palette, size)
	for i := range p {
		p[i] = color.RGBA{
			R: byte(interpolators[0].At(i)),
			G: byte(interpolators[1].At(i)),
			B: byte(interpolators[2].At(i)),
			A: byte(interpolators[3].At(i)),
		}
	}
	return p
}

type Interpolator struct {
	from  float64
	to    float64
	steps int
}

func NewInterpolator(from, to float64, steps int) *Interpolator {
	return &Interpolator{
		from:  from,
		to:    to,
		steps: steps,
	}
}

func (in *Interpolator) At(index int) float64 {
	if in.steps <= 1 {
		return in.to
	}
	diff := in.to - in.from
	fraction := float64(index) / float64(in.steps-1)
	change := diff * fraction
	return in.from + change
}
