package main

import "math"

type Texture interface {
	Value(u, v float64, p *Vec3) Vec3
}

type SolidColor struct {
	Albedo Vec3
}

func NewSolidColor(albedo Vec3) SolidColor {
	return SolidColor{Albedo: albedo}
}
func NewSolidColorFromRGB(r, g, b float64) SolidColor {
	return NewSolidColor(NewVec3(r, g, b))
}
func (c *SolidColor) Value(u, v float64, p *Vec3) Vec3 {
	return c.Albedo
}

type CheckeredTexture struct {
	InvScale float64
	Even     Texture
	Odd      Texture
}

func NewCheckeredTexture(scale float64, even, odd Texture) CheckeredTexture {
	return CheckeredTexture{InvScale: 1.0 / scale, Even: even, Odd: odd}
}

func NewCheckeredTextureFromColors(scale float64, c1, c2 Vec3) CheckeredTexture {
	t1, t2 := NewSolidColor(c1), NewSolidColor(c2)
	return NewCheckeredTexture(scale, &t1, &t2)
}
func (t *CheckeredTexture) Value(u, v float64, p *Vec3) Vec3 {
	isEven := int((math.Floor(t.InvScale*p.X)+math.Floor(t.InvScale*p.Y)+math.Floor(t.InvScale*p.Z)))%2 == 0
	// even & odd map to the two different colors in the checkered pattern
	if isEven {
		return t.Even.Value(u, v, p)
	} else {
		return t.Odd.Value(u, v, p)
	}
}

type NoiseTexture struct {
	Noise PerlinNoise
	Scale float64
}

func NewNoiseTexture(scale float64) NoiseTexture {
	return NoiseTexture{Noise: NewPerlinNoise(), Scale: scale}
}

func (t *NoiseTexture) Value(u, v float64, p *Vec3) Vec3 {
	// return NewVec3(1.0, 1.0, 1.0).Scale(t.Noise.Turbulence(p, 7))
	return NewVec3(0.5, 0.5, 0.5).Scale(1 + math.Sin(t.Scale*p.Z+10*t.Noise.Turbulence(p, 7)))
}
