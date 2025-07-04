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
