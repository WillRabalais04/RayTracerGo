package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"strings"
)

type Vec3 struct {
	X, Y, Z float32
}

// Vec3 Ops
func NewVec3(x, y, z float32) Vec3 {
	return Vec3{X: x, Y: y, Z: z}
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
func (v Vec3) Scale(t float32) Vec3 {
	return Vec3{X: t * v.X, Y: t * v.Y, Z: t * v.Z}
}
func (v *Vec3) ScaleAssign(t float32) {
	v.X *= t
	v.Y *= t
	v.Z *= t
}
func (v1 Vec3) Dot(v2 Vec3) float32 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}
func (v1 Vec3) Cross(v2 Vec3) Vec3 { // check this
	return Vec3{
		X: (v1.Y*v2.Z - v1.Z*v2.Y),
		Y: (v1.Z*v2.X - v1.X*v2.Z),
		Z: (v1.X*v2.Y - v1.Y*v2.X)}
}
func (v *Vec3) MakeUnit() {
	v.ScaleAssign((1.0 / v.Length()))
}
func (v Vec3) GetUnit() Vec3 {
	return v.Scale((1.0 / v.Length()))
}
func (v Vec3) Negate() Vec3 {
	return Vec3{X: -v.X, Y: -v.Y, Z: -v.Z}
}
func (v Vec3) SquaredLength() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}
func (v Vec3) Length() float32 {
	return float32(math.Sqrt(float64(v.SquaredLength())))
}
func RandomUnitVector() Vec3 {
	for {
		p := NewBoundedRandomVec(-1, 1)
		lensq := p.SquaredLength()
		if 1e-160 < lensq && lensq <= 1 {
			return p.Scale(1.0 / float32(math.Sqrt(float64(lensq))))
		}
	}

}
func (Normal Vec3) RandomOnHemisphere() Vec3 {
	OnUnitSphere := RandomUnitVector()
	if OnUnitSphere.Dot(Normal) < 0.0 {
		OnUnitSphere.ScaleAssign(-1.0)
	}
	return OnUnitSphere
}
func NewRandomVec() Vec3 {
	return NewVec3(rand.Float32(), rand.Float32(), rand.Float32())
}
func NewBoundedRandomVec(min, max float32) Vec3 {
	r := max - min
	return NewVec3(rand.Float32()*r+min, rand.Float32()*r+min, rand.Float32()*r+min)
}

type Ray struct {
	Origin    Vec3
	Direction Vec3
}

// Ray ops
func (r Ray) PointAtParameter(t float32) Vec3 {
	return r.Origin.Add(r.Direction.Scale(t))
}
func NewRay(o, d Vec3) Ray {
	return Ray{Origin: o, Direction: d}
}
func (r Ray) Color(world Hittable) Vec3 {
	var rec HitRecord
	if world.Hit(r, 1e-3, math.MaxFloat32, &rec) {
		direction := rec.Normal.RandomOnHemisphere()
		return (NewRay(rec.P, direction).Color(world)).Scale(0.5)
	} else {
		unitDirection := r.Direction.GetUnit()
		t := 0.5 * (unitDirection.Y + 1.0)
		// 	// ((1-t) * <1,1,1>) + (t * <0.5,0.7,1>)
		return (NewVec3(1.0, 1.0, 1.0).Scale(1.0 - t)).Add((NewVec3(0.5, 0.7, 1.0).Scale(t)))
	}
}
func (r Ray) HitSphere(center Vec3, radius float32) float32 {
	oc := r.Origin.Sub(center)
	// <dot(ray_direction, ray_direction), 2*dot(oc, ray_direction), dot(oc,oc) - radius^2>
	a, b, c := r.Direction.Dot(r.Direction), 2.0*oc.Dot(r.Direction), oc.Dot(oc)-radius*radius
	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return -1.0
	} else {
		return ((-1.0 * b) - float32(math.Sqrt(float64(discriminant)))) / (2.0 * a)
	}
}

type HitRecord struct {
	T      float32
	P      Vec3
	Normal Vec3
}

type Hittable interface {
	Hit(r Ray, t_min, t_max float32, rec *HitRecord) bool
}
type HittableList struct {
	List []Hittable
}

func NewHittableList(h ...Hittable) *HittableList {
	return &HittableList{
		List: h,
	}
}
func (hl *HittableList) Add(h Hittable) {
	hl.List = append(hl.List, h)
}
func (hl *HittableList) Hit(r Ray, t_min, t_max float32, rec *HitRecord) bool {
	var temp HitRecord
	hitAnything := false
	closestSoFar := t_max
	for _, hittableObject := range hl.List {
		if hittableObject.Hit(r, t_min, closestSoFar, &temp) {
			hitAnything = true
			closestSoFar = temp.T
			*rec = temp
		}
	}
	return hitAnything
}

type Sphere struct {
	Center Vec3
	Radius float32
}

func NewSphere(c Vec3, r float32) Sphere {
	return Sphere{Center: c, Radius: r}
}
func (s *Sphere) Hit(r Ray, t_min, t_max float32, rec *HitRecord) bool {
	oc := r.Origin.Sub(s.Center)
	a, b, c := r.Direction.Dot(r.Direction), oc.Dot(r.Direction), oc.Dot(oc)-(s.Radius*s.Radius)
	discriminant := b*b - a*c

	if discriminant > 0 {
		sqrtDiscriminant := float32(math.Sqrt(float64(discriminant)))
		posRoot, negRoot := (-b-sqrtDiscriminant)/a, (-b+sqrtDiscriminant)/a
		hitT := float32(-1.0)
		if posRoot < t_max && posRoot > t_min {
			hitT = posRoot
		} else if negRoot < t_max && negRoot > t_min {
			hitT = negRoot
		}
		if hitT != -1.0 {
			rec.T = hitT
			rec.P = r.PointAtParameter(rec.T)
			rec.Normal = rec.P.Sub(s.Center).Scale(1.0 / s.Radius)
			return true
		}
	}
	return false
}

type Camera struct {
	origin          Vec3
	lowerLeftCorner Vec3
	horizontal      Vec3
	vertical        Vec3
}

func NewCamera() Camera {
	return Camera{
		lowerLeftCorner: NewVec3(-2.0, -1.0, -1.0),
		horizontal:      NewVec3(4.0, 0.0, 0.0),
		vertical:        NewVec3(0.0, 2.0, 0.0),
		origin:          NewVec3(0.0, 0.0, 0.0),
	}
}
func (c Camera) GetRay(u, v float32) Ray {
	// llc + u*horizontal + v*vertical - origin
	direction := c.lowerLeftCorner.Add(c.horizontal.Scale(u).Add(c.vertical.Scale(v))).Sub(c.origin)
	return NewRay(c.origin, direction)
}

func main() {
	nx, ny, ns := 400, 200, 200
	var out strings.Builder
	out.WriteString(fmt.Sprintf("P3\n%d %d\n255\n", nx, ny))
	s1, s2 := &Sphere{Center: NewVec3(0.0, 0.0, -1.0), Radius: 0.5}, &Sphere{Center: NewVec3(0.0, -100.5, -1.0), Radius: 100}
	world := NewHittableList(s1, s2)
	cam := NewCamera()
	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			col := NewVec3(0.0, 0.0, 0.0)
			for s := 0; s < ns; s++ {
				u, v := float32(float64(i)+rand.Float64())/float32(nx), float32(float64(j)+rand.Float64())/float32(ny)
				r := cam.GetRay(u, v)
				// p := r.PointAtParameter(2.0)
				col.PlusEq(r.Color(world))
			}
			col.ScaleAssign(1.0 / float32(ns))
			// col = NewVec3(float32(math.Sqrt(float64(col.X))), float32(math.Sqrt(float64(col.X))), float32(math.Sqrt(float64(col.Z)))) // optional gamma corection
			ir, ig, ib := int(255.99*col.X), int(255.99*col.Y), int(255.99*col.Z)
			out.WriteString(fmt.Sprintf("%d %d %d\n", ir, ig, ib))
		}
	}
	err := os.WriteFile("out.ppm", []byte(out.String()), 0644)
	if err != nil {
		fmt.Println("Error writing file: ", err)
	}

}
