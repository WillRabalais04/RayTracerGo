package main

import "math"

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
