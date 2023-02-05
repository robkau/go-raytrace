package patterns

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"math"
)

type TextureMapPattern struct {
	basePattern
	uv     *CheckerPatternUV
	mapper UvMappingF
}

func NewTextureMapPattern(uv *CheckerPatternUV, mapper UvMappingF) *TextureMapPattern {
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

func UvPatternAt(c *CheckerPatternUV, u, v float64) colors.Color {
	u2 := int(math.Floor(u * float64(c.w)))
	v2 := int(math.Floor(v * float64(c.h)))

	if (u2+v2)%2 == 0 {
		return c.a
	} else {
		return c.b
	}
}

// 3d point to UV-mapped 2d point
type UvMappingF func(p geom.Tuple) (u, v float64)

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

// todo top/bottom same as cube mapping later
func CylindricalMap(p geom.Tuple) (u, v float64) {
	// compute the azimuthal angle, same as with spherical_map()
	theta := math.Atan2(p.X, p.Z)
	rawU := theta / (2 * math.Pi)
	u = 1 - (rawU + 0.5)

	// let v go from 0 to 1 between whole units of y
	v = remainderOfOneCloserToZero(p.Y)

	return
}

// todo revisit me.
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
