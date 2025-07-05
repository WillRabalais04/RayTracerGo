package main

import (
	"fmt"
	"math"
	"math/rand/v2"
)

type Vec3 struct {
	X, Y, Z float64
}

// Vec3 Ops
func NewVec3(x, y, z float64) Vec3 {
	return Vec3{X: x, Y: y, Z: z}
}
func NewVec3Pointer(x, y, z float64) *Vec3 {
	vec := Vec3{X: x, Y: y, Z: z}
	return &vec
}
func (v1 Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{X: v1.X + v2.X, Y: v1.Y + v2.Y, Z: v1.Z + v2.Z}
}
func (v1 Vec3) Sub(v2 Vec3) Vec3 {
	return Vec3{X: v1.X - v2.X, Y: v1.Y - v2.Y, Z: v1.Z - v2.Z}
}
func (v1 Vec3) Mul(v2 Vec3) Vec3 {
	return Vec3{X: v1.X * v2.X, Y: v1.Y * v2.Y, Z: v1.Z * v2.Z}
}
func (v1 Vec3) Div(v2 Vec3) Vec3 {
	return Vec3{X: v1.X / v2.X, Y: v1.Y / v2.Y, Z: v1.Z / v2.Z}
}
func (v1 *Vec3) PlusEq(v2 Vec3) {
	v1.X += v2.X
	v1.Y += v2.Y
	v1.Z += v2.Z
}
func (v1 *Vec3) MinEq(v2 Vec3) {
	v1.X -= v2.X
	v1.Y -= v2.Y
	v1.Z -= v2.Z
}
func (v1 *Vec3) TimesEq(v2 Vec3) {
	v1.X *= v2.X
	v1.Y *= v2.Y
	v1.Z *= v2.Z
}
func (v1 *Vec3) SlashEq(v2 Vec3) {
	v1.X /= v2.X
	v1.Y /= v2.Y
	v1.Z /= v2.Z
}
func (v Vec3) Scale(t float64) Vec3 {
	return Vec3{X: t * v.X, Y: t * v.Y, Z: t * v.Z}
}
func (v *Vec3) ScaleAssign(t float64) {
	v.X *= t
	v.Y *= t
	v.Z *= t
}
func Dot(v1, v2 *Vec3) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}
func Cross(v1, v2 *Vec3) Vec3 {
	return Vec3{
		X: (v1.Y*v2.Z - v1.Z*v2.Y),
		Y: (v1.Z*v2.X - v1.X*v2.Z),
		Z: (v1.X*v2.Y - v1.Y*v2.X)}
}
func (v *Vec3) MakeUnitVec() {
	v.ScaleAssign((1.0 / v.Length()))
}
func (v Vec3) GetUnitVec() Vec3 {
	return v.Scale((1.0 / v.Length()))
}
func (v Vec3) Negate() Vec3 {
	return Vec3{X: -v.X, Y: -v.Y, Z: -v.Z}
}
func (v *Vec3) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}
func (v Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}
func (Normal *Vec3) RandomOnHemisphere() Vec3 {
	OnUnitSphere := RandomUnitVector()
	if Dot(&OnUnitSphere, Normal) < 0.0 {
		OnUnitSphere.ScaleAssign(-1.0)
	}
	return OnUnitSphere
}
func (v *Vec3) NearZero() bool {
	s := 1e-8
	return (math.Abs(v.X) < s) && math.Abs(v.Y) < s && math.Abs(v.Z) < s
}
func (v *Vec3) GetDim(i int) float64 {
	if i == 0 {
		return v.X
	}
	if i == 1 {
		return v.Y
	}
	if i == 2 {
		return v.Z
	}

	return -1.0
}
func (v *Vec3) SetDim(i int, val float64) {
	if i == 0 {
		v.X = val
	}
	if i == 1 {
		v.Y = val
	}
	if i == 2 {
		v.Z = val
	}
}

func Reflect(v, n *Vec3) Vec3 {
	return v.Sub(n.Scale(2 * Dot(v, n)))
}
func Refract(uv, n *Vec3, etaiOverEtat float64) Vec3 {
	negatedUV := uv.Negate()
	cosTheta := min(Dot(&negatedUV, n), 1.0)
	rOutPerp := (uv.Add(n.Scale(cosTheta))).Scale(etaiOverEtat)
	rOutParallel := n.Scale(-math.Sqrt(math.Abs(1.0 - rOutPerp.LengthSquared())))
	return rOutParallel.Add(rOutPerp)
}

func RandomUnitVector() Vec3 {
	for {
		p := NewBoundedRandomVec(-1, 1)
		lensq := p.LengthSquared()
		if 1e-160 < lensq && lensq <= 1 {
			return p.Scale(1.0 / math.Sqrt(lensq))
		}
	}
}
func RandomInUnitDisk() Vec3 {
	for {
		p := NewVec3(rand.Float64()*2-1, rand.Float64()*2-1, rand.Float64()*2-1)
		if p.LengthSquared() < 1 {
			return p
		}
	}
}
func NewRandomVec() Vec3 {
	return NewVec3(rand.Float64(), rand.Float64(), rand.Float64())
}
func NewBoundedRandomVec(min, max float64) Vec3 {
	r := max - min
	return NewVec3(rand.Float64()*r+min, rand.Float64()*r+min, rand.Float64()*r+min)
}

func (v *Vec3) PrintVec() {
	fmt.Printf("<%f,%f,%f>", v.X, v.Y, v.Z)
}
