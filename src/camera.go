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
	DefocusAngle      float64
	FocusDistance     float64
	Center            Vec3
	PixelDeltaU       Vec3
	PixelDeltaV       Vec3
	Pixel00Loc        Vec3
	LookFrom          Vec3
	LookAt            Vec3
	VUP               Vec3
	U                 Vec3
	V                 Vec3
	W                 Vec3
	DefocusDiskU      Vec3
	DefocusDiskV      Vec3
	Background        Vec3
}

func NewCamera() Camera {
	return Camera{
		ImageWidth:      100,
		SamplesPerPixel: 10,
		MaxDepth:        50,
		AspectRatio:     1.0,
		VFov:            90,
		DefocusAngle:    0,
		FocusDistance:   10,
		LookFrom:        NewVec3(0.0, 0.0, 0.0),
		LookAt:          NewVec3(0.0, 0.0, -1.0),
		VUP:             NewVec3(0.0, 1.0, 0.0),
	}
}
func (c *Camera) InitCamera() {
	c.ImageHeight = max(int(float64(c.ImageWidth)/c.AspectRatio), 1)
	c.Center = c.LookFrom
	c.PixelSamplesScale = 1.0 / float64(c.SamplesPerPixel)

	theta := DegreesToRadians(c.VFov)
	h := math.Tan(theta / 2)

	viewPortHeight := 2.0 * h * c.FocusDistance
	viewPortWidth := viewPortHeight * (float64(c.ImageWidth) / float64(c.ImageHeight))

	c.W = (c.LookFrom.Sub(c.LookAt)).GetUnitVec()
	c.U = (Cross(&c.VUP, &c.W)).GetUnitVec()
	c.V = Cross(&c.W, &c.U)

	viewPortU := c.U.Scale(viewPortWidth)
	viewPortV := ((c.V).Negate()).Scale(viewPortHeight)

	c.PixelDeltaU, c.PixelDeltaV = viewPortU.Scale(1.0/float64(c.ImageWidth)), viewPortV.Scale(1.0/float64(c.ImageHeight))
	viewPortUpperLeft := c.Center.Sub(c.W.Scale(c.FocusDistance)).Sub(viewPortU.Scale(0.5)).Sub(viewPortV.Scale(0.5)) // center - <0,0,focal length> - (viewportU / 2) - (viewportV / 2)
	c.Pixel00Loc = viewPortUpperLeft.Add((c.PixelDeltaU.Add(c.PixelDeltaV)).Scale(0.5))
	// viewPortUpperLeft + 0.5*(PixelDeltaU + PixelDeltaV)
	defocusRadius := c.FocusDistance * math.Tan(DegreesToRadians(c.DefocusAngle/2))
	c.DefocusDiskU = c.U.Scale(defocusRadius)
	c.DefocusDiskV = c.V.Scale(defocusRadius)

}
func (c *Camera) GetRay(i, j float64) Ray {
	offset := NewBoundedRandomVec(-0.5, 0.5)
	pixelSample := c.Pixel00Loc.Add(c.PixelDeltaU.Scale(i + offset.X).Add(c.PixelDeltaV.Scale(j + offset.Y)))

	var rayOrigin Vec3
	if c.DefocusAngle <= 0 {
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
	if !world.Hit(&r, NewInterval(0.001, math.Inf(1)), &rec) {
		return c.Background
	}

	var scattered Ray
	var attenuation Vec3
	colorFromEmission := (*rec.MaterialPointer).Emitted(rec.U, rec.V, &rec.P)

	if !(*rec.MaterialPointer).Scatter(&r, &rec, &attenuation, &scattered) {
		return colorFromEmission
	}
	colorFromScatter := attenuation.Mul(c.RayColor(scattered, depth-1, world))

	return colorFromEmission.Add(colorFromScatter)
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
	return c.Center.Add(c.DefocusDiskU.Scale(p.X).Add(c.DefocusDiskV.Scale(p.Y)))
}
