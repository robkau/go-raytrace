package patterns

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewStripePattern(t *testing.T) {
	pattern := NewStripePattern(NewSolidColorPattern(colors.White()), NewSolidColorPattern(colors.Black()))

	assert.Equal(t, geom.NewIdentityMatrixX4(), pattern.GetTransform())
	assert.Equal(t, NewSolidColorPattern(colors.White()), pattern.a)
	assert.Equal(t, NewSolidColorPattern(colors.Black()), pattern.b)
}

func Test_NewStripePattern_AssignTransformation(t *testing.T) {
	pattern := NewStripePattern(NewSolidColorPattern(colors.White()), NewSolidColorPattern(colors.Black()))
	pattern.SetTransform(geom.Translate(0.5, 0, 0))

	assert.Equal(t, geom.Translate(0.5, 0, 0), pattern.GetTransform())
	assert.Equal(t, NewSolidColorPattern(colors.White()), pattern.a)
	assert.Equal(t, NewSolidColorPattern(colors.Black()), pattern.b)
}

func Test_StripeConstantY(t *testing.T) {
	pattern := NewStripePattern(NewSolidColorPattern(colors.White()), NewSolidColorPattern(colors.Black()))

	assert.Equal(t, colors.White(), pattern.ColorAt(geom.NewPoint(0, 0, 0)))
	assert.Equal(t, colors.White(), pattern.ColorAt(geom.NewPoint(0, 1, 0)))
	assert.Equal(t, colors.White(), pattern.ColorAt(geom.NewPoint(0, 2, 0)))
}

func Test_StripeConstantZ(t *testing.T) {
	pattern := NewStripePattern(NewSolidColorPattern(colors.White()), NewSolidColorPattern(colors.Black()))

	assert.Equal(t, colors.White(), pattern.ColorAt(geom.NewPoint(0, 0, 0)))
	assert.Equal(t, colors.White(), pattern.ColorAt(geom.NewPoint(0, 0, 1)))
	assert.Equal(t, colors.White(), pattern.ColorAt(geom.NewPoint(0, 0, 2)))
}

func Test_StripeChangesX(t *testing.T) {
	pattern := NewStripePattern(NewSolidColorPattern(colors.White()), NewSolidColorPattern(colors.Black()))

	assert.Equal(t, colors.White(), pattern.ColorAt(geom.NewPoint(0, 0, 0)))
	assert.Equal(t, colors.White(), pattern.ColorAt(geom.NewPoint(0.9, 0, 0)))
	assert.Equal(t, colors.Black(), pattern.ColorAt(geom.NewPoint(1, 0, 0)))
	assert.Equal(t, colors.Black(), pattern.ColorAt(geom.NewPoint(-0.1, 0, 0)))
	assert.Equal(t, colors.Black(), pattern.ColorAt(geom.NewPoint(-1, 0, 0)))
	assert.Equal(t, colors.White(), pattern.ColorAt(geom.NewPoint(-1.1, 0, 0)))
}

func Test_Stripe_WithObjectTransformation(t *testing.T) {
	pattern := NewStripePattern(NewSolidColorPattern(colors.White()), NewSolidColorPattern(colors.Black()))

	assert.Equal(t, colors.White(), pattern.ColorAtShape(geom.DoubleScale, geom.NewPoint(1.5, 0, 0)))
}

func Test_Stripe_WithPatternTransformation(t *testing.T) {
	pattern := NewStripePattern(NewSolidColorPattern(colors.White()), NewSolidColorPattern(colors.Black()))
	pattern.SetTransform(geom.Scale(2, 2, 2))

	assert.Equal(t, colors.White(), pattern.ColorAtShape(geom.IdentityTransform, geom.NewPoint(1.5, 0, 0)))
}

func Test_Stripe_WithObjectAndPatternTransformation(t *testing.T) {
	pattern := NewStripePattern(NewSolidColorPattern(colors.White()), NewSolidColorPattern(colors.Black()))
	pattern.SetTransform(geom.Translate(0.5, 0, 0))

	assert.Equal(t, colors.White(), pattern.ColorAtShape(geom.DoubleScale, geom.NewPoint(2.5, 0, 0)))
}
