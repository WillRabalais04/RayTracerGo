package main

import "math"

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
