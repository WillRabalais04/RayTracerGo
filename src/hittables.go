package main

import (
	"math"
	"math/rand/v2"
)

type HitRecord struct {
	T               float64
	U               float64
	V               float64
	P               Vec3
	Normal          Vec3
	FrontFace       bool
	MaterialPointer *Material
}

func (h *HitRecord) SetFaceNormal(r *Ray, outwardNormal Vec3) {
	h.FrontFace = Dot(&r.Direction, &outwardNormal) < 0
	if h.FrontFace {
		h.Normal = outwardNormal
	} else {
		h.Normal = outwardNormal.Negate()
	}
}

type Hittable interface {
	Hit(r *Ray, i Interval, rec *HitRecord) bool
	BBOX() *AABB
}
type HittableList struct {
	Objects   []Hittable
	BBOXField AABB
} // pointer to hittable list implements hittable but hittable list itself does not

func NewHittableList(h ...Hittable) *HittableList {
	return &HittableList{
		Objects: h, BBOXField: NewAABB(EmptyInterval, EmptyInterval, EmptyInterval), // check this
	}
}
func (hl *HittableList) Add(h Hittable) {
	hl.Objects = append(hl.Objects, h)
	hl.BBOXField = MergeAABBs(&hl.BBOXField, h.BBOX())
}
func (hl *HittableList) Hit(r *Ray, i Interval, rec *HitRecord) bool {
	var temp HitRecord
	hitAnything := false
	closestSoFar := i.Max
	for _, hittableObject := range hl.Objects {
		if hittableObject.Hit(r, NewInterval(i.Min, closestSoFar), &temp) {
			hitAnything = true
			closestSoFar = temp.T
			*rec = temp
		}
	}
	return hitAnything
}
func (hl *HittableList) BBOX() *AABB {
	return &hl.BBOXField
}

type Sphere struct {
	Center    Ray
	Radius    float64
	Mat       Material
	BBOXField AABB
}

func NewSphere(c Vec3, r float64, mat Material) Sphere {
	rvec := NewVec3(r, r, r)
	bbox := NewAABBFromPoints(c.Sub(rvec), c.Add(rvec))
	return Sphere{Center: NewRay(c, NewVec3(0.0, 0.0, 0.0), 0), Radius: max(0, r), Mat: mat, BBOXField: bbox}
}
func NewMovingSphere(c1, c2 Vec3, r float64, mat Material) Sphere {
	center := NewRay(c1, c2.Sub(c1), 0)
	rvec := NewVec3(r, r, r)
	box1 := NewAABBFromPoints(center.at(0).Sub(rvec), center.at(0).Add(rvec))
	box2 := NewAABBFromPoints(center.at(1).Sub(rvec), center.at(1).Add(rvec))
	bbox := MergeAABBs(&box1, &box2)
	return Sphere{Center: center, Radius: max(0, r), Mat: mat, BBOXField: bbox}
}
func (s *Sphere) Hit(r *Ray, i Interval, rec *HitRecord) bool {
	currentCenter := s.Center.at(r.Time)
	oc := currentCenter.Sub(r.Origin)

	a, h, c := r.Direction.LengthSquared(), Dot(&r.Direction, &oc), oc.LengthSquared()-(s.Radius*s.Radius)
	discriminant := h*h - a*c

	if discriminant < 0 {
		return false
	}
	sqrtDiscriminant := math.Sqrt(discriminant)

	root := (h - sqrtDiscriminant) / a

	if !i.Surrounds(root) {
		root = (h + sqrtDiscriminant) / a
		if !i.Surrounds(root) {
			return false
		}
	}

	rec.T = root
	rec.P = r.at(rec.T)
	outwardNormal := rec.P.Sub(currentCenter).Scale(1.0 / s.Radius)
	rec.SetFaceNormal(r, outwardNormal)
	GetSphereUV(&outwardNormal, &rec.U, &rec.V)
	rec.MaterialPointer = &s.Mat

	return true

}
func (s *Sphere) BBOX() *AABB {
	return &s.BBOXField
}

func GetSphereUV(p *Vec3, u, v *float64) {
	theta := math.Acos(-p.Y)
	phi := math.Atan2(-p.Z, p.X) + math.Pi

	*u = phi / (2 * math.Pi)
	*v = theta / math.Pi
}

type Quad struct {
	Q, U, V, W, Normal Vec3
	D                  float64
	Mat                Material
	BBOXField          AABB
}

func NewQuad(q, u, v Vec3, m Material) Quad {
	quad := Quad{Q: q, U: u, V: v, Mat: m}
	n := Cross(&u, &v)
	quad.Normal = n.GetUnitVec()
	quad.D = Dot(&quad.Normal, &q)
	quad.W = n.Scale(1.0 / Dot(&n, &n))
	quad.SetBoundingBox()
	return quad

}
func (q *Quad) SetBoundingBox() {
	bboxDiag1 := NewAABBFromPoints(q.Q, q.Q.Add(q.U.Add(q.V)))
	bboxDiag2 := NewAABBFromPoints(q.Q.Add(q.U), q.Q.Add(q.V))
	q.BBOXField = MergeAABBs(&bboxDiag1, &bboxDiag2)
}

func (q *Quad) Hit(r *Ray, i Interval, rec *HitRecord) bool {
	denominator := Dot(&q.Normal, &r.Direction)
	if math.Abs(denominator) < 1e-8 {
		return false
	}
	t := (q.D - Dot(&q.Normal, &r.Origin)) / denominator
	if !i.Contains(t) {
		return false
	}

	intersection := r.at(t)
	PlanarHitPointVector := intersection.Sub(q.Q)
	phpXv := Cross(&PlanarHitPointVector, &q.V)
	uXphp := Cross(&q.U, &PlanarHitPointVector)
	alpha := Dot(&q.W, &phpXv)
	beta := Dot(&q.W, &uXphp)
	if !q.IsInterior(alpha, beta, rec) {
		return false
	}
	rec.T = t
	rec.P = intersection
	rec.MaterialPointer = &q.Mat
	rec.SetFaceNormal(r, q.Normal)
	return true
}

func (q *Quad) IsInterior(a, b float64, rec *HitRecord) bool {
	// could explore different shapes
	unitInterval := NewInterval(0, 1)
	if !unitInterval.Contains(a) || !unitInterval.Contains(b) {
		return false
	}

	rec.U = a
	rec.V = b
	return true
}

func (q *Quad) BBOX() *AABB {
	return &q.BBOXField
}

func NewBox(a, b Vec3, m Material) *HittableList {

	min := NewVec3(min(a.X, b.X), min(a.Y, b.Y), min(a.Z, b.Z))
	max := NewVec3(max(a.X, b.X), max(a.Y, b.Y), max(a.Z, b.Z))

	dx := NewVec3(max.X-min.X, 0, 0)
	dy := NewVec3(0, max.Y-min.Y, 0)
	dz := NewVec3(0, 0, max.Z-min.Z)

	q1 := NewQuad(NewVec3(min.X, min.Y, max.Z), dx, dy, m)           // front
	q2 := NewQuad(NewVec3(max.X, min.Y, max.Z), dz.Scale(-1), dy, m) // right
	q3 := NewQuad(NewVec3(max.X, min.Y, min.Z), dx.Scale(-1), dy, m) // back
	q4 := NewQuad(NewVec3(min.X, min.Y, min.Z), dz, dy, m)           // left
	q5 := NewQuad(NewVec3(min.X, max.Y, max.Z), dx, dz.Scale(-1), m) // top
	q6 := NewQuad(NewVec3(min.X, min.Y, min.Z), dx, dz, m)           // bottom
	return NewHittableList(&q1, &q2, &q3, &q4, &q5, &q6)

}

func NewCornellBox(left, back, right, light Material) *HittableList {

	q1 := NewQuad(NewVec3(555, 0, 0), NewVec3(0, 555, 0), NewVec3(0, 0, 555), left)
	q2 := NewQuad(NewVec3(0, 0, 0), NewVec3(0, 555, 0), NewVec3(0, 0, 555), right)
	q3 := NewQuad(NewVec3(113, 554, 127), NewVec3(330, 0, 0), NewVec3(0, 0, 305), light)
	q4 := NewQuad(NewVec3(0, 0, 0), NewVec3(555, 0, 0), NewVec3(0, 0, 555), back)
	q5 := NewQuad(NewVec3(555, 555, 555), NewVec3(-555, 0, 0), NewVec3(0, 0, -555), back)
	q6 := NewQuad(NewVec3(0, 0, 555), NewVec3(555, 0, 0), NewVec3(0, 555, 0), back)

	return NewHittableList(&q1, &q2, &q3, &q4, &q5, &q6)
}

type ConstantMedium struct {
	Boundary      *Hittable
	NegInvDensity float64
	PhaseFunction Material
	BBOXField     AABB
}

func NewConstantMedium(boundary Hittable, density float64, tex *Texture) ConstantMedium {
	t := NewIsotropicFromTexture(tex)
	return ConstantMedium{Boundary: &boundary, NegInvDensity: -1 / density, PhaseFunction: &t}
}

func NewConstantMediumFromColor(boundary Hittable, density float64, c *Vec3) ConstantMedium {
	t := NewIsotropic(c)
	return ConstantMedium{Boundary: &boundary, NegInvDensity: -1 / density, PhaseFunction: &t}
}

func (c *ConstantMedium) Hit(r *Ray, i Interval, rec *HitRecord) bool {
	var rec1, rec2 HitRecord

	if !(*c.Boundary).Hit(r, UniverseInterval, &rec1) {
		return false
	}

	if !(*c.Boundary).Hit(r, NewInterval(rec1.T+0.0001, math.Inf(1)), &rec2) {
		return false
	}

	rec1.T = max(rec1.T, i.Min)
	rec2.T = min(rec2.T, i.Max)

	if rec1.T >= rec2.T {
		return false
	}

	rec1.T = max(0, rec1.T)

	rayLength := r.Direction.Length()
	distanceInsideBoundary := (rec2.T - rec1.T) * rayLength
	hitDistance := c.NegInvDensity * math.Log(rand.Float64())

	if hitDistance > distanceInsideBoundary {
		return false
	}

	rec.T = rec1.T + hitDistance/rayLength
	rec.P = r.at(rec.T)

	rec.Normal = NewVec3(1, 0, 0)
	rec.FrontFace = true
	rec.MaterialPointer = &c.PhaseFunction

	return true

}
func (c *ConstantMedium) BBOX() *AABB {
	return (*c.Boundary).BBOX()
}
