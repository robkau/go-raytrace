package patterns

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"math"
)

type TextureMapPattern struct {
	basePattern
	uv     UvPattern
	mapper UvMappingF
}

type UvPattern interface {
	ColorAt(u, v float64) colors.Color
}

// 3d point to UV-mapped 2d point
type UvMappingF func(p geom.Tuple) (u, v float64)

func NewTextureMapPattern(uv UvPattern, mapper UvMappingF) *TextureMapPattern {
	return &TextureMapPattern{
		basePattern: newBasePattern(),
		uv:          uv,
		mapper:      mapper,
	}
}

func (p *TextureMapPattern) ColorAt(t geom.Tuple) colors.Color {
	u, v := p.mapper(t)
	return UvPatternAt(p.uv, u, v)
}

func (p *TextureMapPattern) ColorAtShape(wtof WorldToObjectF, t geom.Tuple) colors.Color {
	return ColorAtShape(p, wtof, t)
}

// CheckerPatternUV is similar to CheckerPattern but only considers 2 dimensions.
type CheckerPatternUV struct {
	w int
	h int
	a colors.Color
	b colors.Color
}

func NewCheckerPatternUV(width, height int, a, b colors.Color) *CheckerPatternUV {
	return &CheckerPatternUV{
		w: width,
		h: height,
		a: a,
		b: b,
	}
}

func (p *CheckerPatternUV) ColorAt(u, v float64) colors.Color {
	u2 := int(math.Floor(u * float64(p.w)))
	v2 := int(math.Floor(v * float64(p.h)))

	if (u2+v2)%2 == 0 {
		return p.a
	} else {
		return p.b
	}
}

func UvPatternAt(p UvPattern, u, v float64) colors.Color {
	return p.ColorAt(u, v)
}

type UVAlignCheck struct {
	main colors.Color
	ul   colors.Color
	ur   colors.Color
	bl   colors.Color
	br   colors.Color
}

type CubeMapPattern struct {
	basePattern
	l UvPattern
	f UvPattern
	r UvPattern
	b UvPattern
	u UvPattern
	d UvPattern
}

func NewPrismaticCube() *CubeMapPattern {
	left := NewUVAlignCheck(colors.Yellow(), colors.Cyan(), colors.Red(), colors.Blue(), colors.Brown())
	front := NewUVAlignCheck(colors.Cyan(), colors.Red(), colors.Yellow(), colors.Brown(), colors.Green())
	right := NewUVAlignCheck(colors.Red(), colors.Yellow(), colors.Purple(), colors.Green(), colors.White())
	back := NewUVAlignCheck(colors.Green(), colors.Purple(), colors.Cyan(), colors.White(), colors.Blue())
	up := NewUVAlignCheck(colors.Brown(), colors.Cyan(), colors.Purple(), colors.Red(), colors.Yellow())
	down := NewUVAlignCheck(colors.Purple(), colors.Brown(), colors.Green(), colors.Blue(), colors.White())
	pattern := NewCubeMapPattern(left, front, right, back, up, down)
	return pattern
}

func NewCubeMapPattern(l, f, r, b, u, d UvPattern) *CubeMapPattern {
	return &CubeMapPattern{
		basePattern: newBasePattern(),
		l:           l,
		f:           f,
		r:           r,
		b:           b,
		u:           u,
		d:           d,
	}
}

func (p *CubeMapPattern) ColorAt(t geom.Tuple) colors.Color {
	face := CubeFaceFromPoint(t)

	var u, v float64
	var cubeFace UvPattern

	switch face {
	case Left:
		u, v = CubeUvLeft(t)
		cubeFace = p.l
	case Right:
		u, v = CubeUvRight(t)
		cubeFace = p.r
	case Front:
		u, v = CubeUvFront(t)
		cubeFace = p.f
	case Back:
		u, v = CubeUvBack(t)
		cubeFace = p.b
	case Up:
		u, v = CubeUvUp(t)
		cubeFace = p.u
	default: // down
		u, v = CubeUvDown(t)
		cubeFace = p.d
	}

	return UvPatternAt(cubeFace, u, v)
}

func (p *CubeMapPattern) ColorAtShape(wtof WorldToObjectF, t geom.Tuple) colors.Color {
	return ColorAtShape(p, wtof, t)
}

func NewUVAlignCheck(main, ul, ur, bl, br colors.Color) *UVAlignCheck {
	return &UVAlignCheck{
		main: main,
		ul:   ul,
		ur:   ur,
		bl:   bl,
		br:   br,
	}
}

func (p *UVAlignCheck) ColorAt(u, v float64) colors.Color {
	if v > 0.8 {
		if u < 0.2 {
			return p.ul
		}
		if u > 0.8 {
			return p.ur
		}
	} else if v < 0.2 {
		if u < 0.2 {
			return p.bl
		}
		if u > 0.8 {
			return p.br
		}
	}
	return p.main
}

func SphericalMap(p geom.Tuple) (u, v float64) {
	// compute the azimuthal angle
	// -π < theta <= π
	// angle increases clockwise as viewed from above,
	// which is opposite of what we want, but we'll fix it later.
	theta := math.Atan2(p.X, p.Z)

	// vec is the vector pointing from the sphere's origin (the world origin)
	// to the point, which will also happen to be exactly equal to the sphere's
	// radius.
	vec := geom.NewVector(p.X, p.Y, p.Z)
	radius := vec.Mag()

	// compute the polar angle
	// 0 <= phi <= π
	phi := math.Acos(p.Y / radius)

	// -0.5 < raw_u <= 0.5
	rawU := theta / (2.0 * math.Pi)

	// 0 <= u < 1
	// here's also where we fix the direction of u. Subtract it from 1,
	// so that it increases counterclockwise as viewed from above.
	u = 1 - (rawU + 0.5)

	// we want v to be 0 at the south pole of the sphere,
	// and 1 at the north pole, so we have to "flip it over"
	// by subtracting it from 1.
	v = 1 - (phi / math.Pi)

	return u, v
}

func PlanarMap(p geom.Tuple) (u, v float64) {
	u = remainderOfOneCloserToZero(p.X)
	v = remainderOfOneCloserToZero(p.Z)
	return
}

func CylindricalMap(p geom.Tuple) (u, v float64) {
	// compute the azimuthal angle, same as with spherical_map()
	theta := math.Atan2(p.X, p.Z)
	rawU := theta / (2 * math.Pi)
	u = 1 - (rawU + 0.5)

	// let v go from 0 to 1 between whole units of y
	v = remainderOfOneCloserToZero(p.Y)

	return
}

func remainderOfOneCloserToZero(v float64) float64 {
	var flipped bool
	if v < 0 {
		flipped = true
	}

	vi := math.Abs(math.Mod(v, 1))

	if vi == 0 {
		return vi
	}

	if flipped {
		vi = 1 - vi
	}

	return vi

}

type CubeFace uint8

const (
	Left CubeFace = iota
	Right
	Front
	Back
	Up
	Down
)

// For a unit cube.
func CubeFaceFromPoint(p geom.Tuple) CubeFace {
	coord := math.Max(math.Max(math.Abs(p.X), math.Abs(p.Y)), math.Abs(p.Z))

	if coord == p.X {
		return Right
	}
	if coord == -p.X {
		return Left
	}
	if coord == p.Y {
		return Up
	}
	if coord == -p.Y {
		return Down
	}
	if coord == p.Z {
		return Front
	}
	return Back
}

func CubeUvFront(p geom.Tuple) (u, v float64) {
	u = math.Mod(p.X+1., 2) / 2.0
	v = math.Mod(p.Y+1., 2) / 2.0
	return
}

func CubeUvBack(p geom.Tuple) (u, v float64) {
	u = math.Mod(1-p.X, 2) / 2.0
	v = math.Mod(p.Y+1., 2) / 2.0
	return
}

func CubeUvRight(p geom.Tuple) (u, v float64) {
	u = math.Mod(1-p.Z, 2) / 2.0
	v = math.Mod(p.Y+1., 2) / 2.0
	return
}

func CubeUvLeft(p geom.Tuple) (u, v float64) {
	u = math.Mod(p.Z+1., 2) / 2.0
	v = math.Mod(p.Y+1., 2) / 2.0
	return
}

func CubeUvUp(p geom.Tuple) (u, v float64) {
	u = math.Mod(p.X+1, 2) / 2.0
	v = math.Mod(1-p.Z, 2) / 2.0
	return
}

func CubeUvDown(p geom.Tuple) (u, v float64) {
	u = math.Mod(p.X+1, 2) / 2.0
	v = math.Mod(p.Z+1., 2) / 2.0
	return
}
