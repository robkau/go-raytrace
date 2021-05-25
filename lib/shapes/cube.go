package shapes

import (
	"go-raytrace/lib/geom"
	"math"
)

type Cube struct {
	t          geom.X4Matrix
	M          Material
	id         string
	shadowless bool
	unshaded   bool
}

func NewCube() Cube {
	return Cube{
		t:  geom.NewIdentityMatrixX4(),
		M:  NewMaterial(),
		id: newId(),
	}
}

func NewCubeWith(t geom.X4Matrix) Cube {
	c := NewCube()
	c.t = t
	return c
}

func NewGlassCube() Cube {
	c := NewCube()

	m := NewGlassMaterial()
	c.M = m
	return c
}

func (c Cube) Id() string {
	return c.id
}

func (c Cube) NormalAt(p geom.Tuple) geom.Tuple {
	return NormalAt(p, c.t, c.LocalNormalAt)
}

func (c Cube) LocalNormalAt(p geom.Tuple) geom.Tuple {
	maxC := math.Max(math.Max(math.Abs(p.X), math.Abs(p.Y)), math.Abs(p.Z))

	if maxC == math.Abs(p.X) {
		return geom.NewVector(p.X, 0, 0)
	} else if maxC == math.Abs(p.Y) {
		return geom.NewVector(0, p.Y, 0)
	} else {
		return geom.NewVector(0, 0, p.Z)
	}
}

func (c Cube) Intersect(r geom.Ray) Intersections {
	return Intersect(r, c.t, c.LocalIntersect)
}

func (c Cube) LocalIntersect(r geom.Ray) Intersections {
	// todo optimization possible
	xtMin, xtMax := checkAxis(r.Origin.X, r.Direction.X)
	ytMin, ytMax := checkAxis(r.Origin.Y, r.Direction.Y)
	ztMin, ztMax := checkAxis(r.Origin.Z, r.Direction.Z)

	tMin := math.Max(math.Max(xtMin, ytMin), ztMin)
	tMax := math.Min(math.Min(xtMax, ytMax), ztMax)

	if tMin > tMax {
		return NewIntersections()
	}

	return NewIntersections(
		Intersection{
			T: tMin,
			O: c,
		},
		Intersection{
			T: tMax,
			O: c,
		},
	)
}

func (c Cube) GetTransform() geom.X4Matrix {
	return c.t
}

func (c Cube) SetTransform(m geom.X4Matrix) Shape {
	c.t = m
	return c
}

func (c Cube) GetMaterial() Material {
	return c.M
}

func (c Cube) SetMaterial(m Material) Shape {
	c.M = m
	return c
}

func (c Cube) GetShadowless() bool {
	return c.shadowless
}

func (c Cube) SetShadowless(ss bool) Shape {
	c.shadowless = ss
	return c
}

func (c Cube) GetShaded() bool {
	return !c.unshaded
}

func (c Cube) SetShaded(ss bool) Shape {
	c.unshaded = !ss
	return c
}

func checkAxis(origin float64, direction float64) (tMin float64, tMax float64) {
	tMinNumerator := -1 - origin
	tMaxNumerator := 1 - origin

	if math.Abs(direction) >= geom.FloatComparisonEpsilon {
		tMin = tMinNumerator / direction
		tMax = tMaxNumerator / direction
	} else {
		tMin = tMinNumerator * math.Inf(1)
		tMax = tMaxNumerator * math.Inf(1)
	}

	if tMin > tMax {
		// swap
		return tMax, tMin
	}
	// no swap
	return

}
