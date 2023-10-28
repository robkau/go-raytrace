package scenes

import (
	"fmt"
	"github.com/robkau/go-raytrace/lib/colors"
	"github.com/robkau/go-raytrace/lib/geom"
	"github.com/robkau/go-raytrace/lib/shapes"
	"github.com/robkau/go-raytrace/lib/view"
	"image/jpeg"
	"log"
	"os"
	"time"
)

func NewBeadTapestry() (*view.World, []CameraLocation) {

	w := view.NewWorld()

	f, err := os.Open("data/png/spooky.jpeg")
	if err != nil {
		log.Println(err)
		panic("opening img")
	}
	p, err := jpeg.Decode(f)
	if err != nil {
		log.Println(err)
		panic("decoding img")
	}
	f.Close()

	b := p.Bounds()
	fmt.Println(b)

	xSpan := b.Max.X - b.Min.X
	ySpan := b.Max.Y - b.Min.Y
	stride := 10
	objectCount := (xSpan / stride) * (ySpan / stride)

	g := shapes.NewGroupWithCapacity(objectCount)

	cameraDistance := 4.5 * float64(xSpan)
	cameraPos := geom.NewPoint(0, -cameraDistance/3, cameraDistance)
	cameraLookingAt := geom.NewPoint(0, 0, 0)

	tStart := time.Now()
	count := 0
	for i := b.Min.X; i < b.Max.X; i += stride {
		for j := b.Min.Y; j < b.Max.Y; j += stride {
			s := shapes.NewSphere()
			m := s.GetMaterial()
			m.Specular = 0
			m.Color = colors.NewColorFromStdlibColor(p.At(i, j))
			s.SetMaterial(m)

			height := float64(ySpan/100) * (m.Color.Mag())

			s.SetTransform(geom.Translate(float64(xSpan/2-i)*2, float64(ySpan/2-j)*2, height).MulX4Matrix(geom.Scale(float64(stride), float64(stride), float64(stride)+height)))

			g.AddChild(s)
			count++
		}
	}
	fmt.Println("expected count", objectCount, "actual count", count)
	fmt.Println("it took to build", time.Since(tStart))

	// light above
	w.AddAreaLight(shapes.NewAreaLight(cameraPos, geom.NewVector(float64(xSpan*2), 0, 0), 6, geom.NewVector(0, float64(ySpan*2), 0), 6, colors.NewColor(1.9, 1.4, 1.4), nil))
	w.AddObject(g)

	tStart = time.Now()
	divideFactor := 128
	w.Divide(objectCount / divideFactor)
	fmt.Println("it took to divide", time.Since(tStart))
	return w, []CameraLocation{CameraLocation{cameraPos, cameraLookingAt}}
}
