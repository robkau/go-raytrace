package materials

import (
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DefaultMaterial(t *testing.T) {
	m := NewMaterial()

	assert.Equal(t, colors.White(), m.Color)
	assert.Equal(t, 0.1, m.Ambient)
	assert.Equal(t, 0.9, m.Diffuse)
	assert.Equal(t, 0.9, m.Specular)
	assert.Equal(t, 200.0, m.Shininess)
	assert.Equal(t, 0.0, m.Reflective)
	assert.Equal(t, 0.0, m.Transparency)
	assert.Equal(t, 1.0, m.RefractiveIndex)
}
