package main

import (
	"fmt"
	"math"
	"math/rand/v2"
)

type Camera struct {
	ImageWidth        int
	ImageHeight       int
	SamplesPerPixel   int
	MaxDepth          int
	AspectRatio       float64
	PixelSamplesScale float64
	VFov              float64
	defocusAngle      float64
	focusDistance     float64
	Center            Vec3
	PixelDeltaU       Vec3
	PixelDeltaV       Vec3
	Pixel00Loc        Vec3
	lookFrom          Vec3
	lookAt            Vec3
	vup               Vec3
	u                 Vec3
	v                 Vec3
	w                 Vec3
	defocusDiskU      Vec3
	defocusDiskV      Vec3
}

func NewCamera() Camera {
	return Camera{
		ImageWidth:      100,
		SamplesPerPixel: 10,
		MaxDepth:        50,
		AspectRatio:     1.0,
		VFov:            90,
		defocusAngle:    0,
		focusDistance:   10,
		lookFrom:        NewVec3(0.0, 0.0, 0.0),
		lookAt:          NewVec3(0.0, 0.0, -1.0),
		vup:             NewVec3(0.0, 1.0, 0.0),
	}
}
func (c *Camera) InitCamera() {
	c.ImageHeight = max(int(float64(c.ImageWidth)/c.AspectRatio), 1)
	c.Center = c.lookFrom
	c.PixelSamplesScale = 1.0 / float64(c.SamplesPerPixel)

	theta := DegreesToRadians(c.VFov)
	h := math.Tan(theta / 2)

	viewPortHeight := 2.0 * h * c.focusDistance
	viewPortWidth := viewPortHeight * (float64(c.ImageWidth) / float64(c.ImageHeight))

	c.w = (c.lookFrom.Sub(c.lookAt)).GetUnitVec()
	c.u = (Cross(&c.vup, &c.w)).GetUnitVec()
	c.v = Cross(&c.w, &c.u)

	viewPortU := c.u.Scale(viewPortWidth)
	viewPortV := ((c.v).Negate()).Scale(viewPortHeight)

	c.PixelDeltaU, c.PixelDeltaV = viewPortU.Scale(1.0/float64(c.ImageWidth)), viewPortV.Scale(1.0/float64(c.ImageHeight))
	viewPortUpperLeft := c.Center.Sub(c.w.Scale(c.focusDistance)).Sub(viewPortU.Scale(0.5)).Sub(viewPortV.Scale(0.5)) // center - <0,0,focal length> - (viewportU / 2) - (viewportV / 2)
	c.Pixel00Loc = viewPortUpperLeft.Add((c.PixelDeltaU.Add(c.PixelDeltaV)).Scale(0.5))
	// viewPortUpperLeft + 0.5*(PixelDeltaU + PixelDeltaV)
	defocusRadius := c.focusDistance * math.Tan(DegreesToRadians(c.defocusAngle/2))
	c.defocusDiskU = c.u.Scale(defocusRadius)
	c.defocusDiskV = c.v.Scale(defocusRadius)

}
func (c *Camera) GetRay(i, j float64) Ray {
	offset := NewBoundedRandomVec(-0.5, 0.5)
	pixelSample := c.Pixel00Loc.Add(c.PixelDeltaU.Scale(i + offset.X).Add(c.PixelDeltaV.Scale(j + offset.Y)))

	var rayOrigin Vec3
	if c.defocusAngle <= 0 {
		rayOrigin = c.Center
	} else {
		rayOrigin = c.defocusDiskSample()
	}
	rayDirection := pixelSample.Sub(rayOrigin)
	return NewRay(rayOrigin, rayDirection, rand.Float64())
}
func (c *Camera) RayColor(r Ray, depth int, world Hittable) Vec3 {
	if depth <= 0 {
		return NewVec3(0.0, 0.0, 0.0)
	}

	var rec HitRecord
	if world.Hit(&r, NewInterval(0.001, math.MaxFloat64), &rec) {
		var scattered Ray
		var attenuation Vec3

		if rec.MaterialPointer != nil && (*rec.MaterialPointer).Scatter(&r, &rec, &attenuation, &scattered) {
			return attenuation.Mul(c.RayColor(scattered, depth-1, world))
		}
		return NewVec3(0, 0, 0)
	}

	// if not hit it renders the background color
	unitDirection := r.Direction.GetUnitVec()
	t := 0.5 * (unitDirection.Y + 1.0)
	// ((1-t) * <1,1,1>) + (t * <0.5,0.7,1>)
	return (NewVec3(1.0, 1.0, 1.0).Scale(1.0 - t)).Add((NewVec3(0.5, 0.7, 1.0).Scale(t)))
}
func (c *Camera) Render(world Hittable) {
	c.InitCamera()
	InitImage(c.ImageWidth, c.ImageHeight)
	lastPercent := -1
	for i := 0; i < c.ImageHeight; i++ {
		for j := 0; j < c.ImageWidth; j++ {
			pixelColor := NewVec3(0.0, 0.0, 0.0)
			for sample := 0; sample < c.SamplesPerPixel; sample++ {
				r := c.GetRay(float64(j), float64(i))
				pixelColor.PlusEq(c.RayColor(r, c.MaxDepth, world))
			}
			WriteColor(pixelColor.Scale(c.PixelSamplesScale))
		}
		percent := (i*100)/c.ImageHeight + 1
		if percent%5 == 0 && percent != lastPercent {
			fmt.Printf("%d percent done.\n", percent)
			lastPercent = percent
		}
	}
}
func (c *Camera) defocusDiskSample() Vec3 {
	p := RandomInUnitDisk()
	return c.Center.Add(c.defocusDiskU.Scale(p.X).Add(c.defocusDiskV.Scale(p.Y)))
}
