package main

import "math"

type Translate struct {
	Offset    Vec3
	Object    *Hittable
	BBOXField AABB
}

func NewTranslateY(offset Vec3, object Hittable) *Translate {
	// bboxfield = object.boundingbox.shift(offset)
	return &Translate{Offset: offset, Object: &object, BBOXField: object.BBOX().ShiftAABB(&offset)}
}

func (t *Translate) Hit(r *Ray, i Interval, rec *HitRecord) bool {
	offsetRay := NewRay(r.Origin.Sub(t.Offset), r.Direction, r.Time)
	if !(*t.Object).Hit(&offsetRay, i, rec) {
		return false
	}
	rec.P.PlusEq(t.Offset)
	return true
}

func (t *Translate) BBOX() *AABB {
	return &t.BBOXField
}

type RotateY struct {
	Object             *Hittable
	CosTheta, SinTheta float64
	BBOXField          AABB
}

func NewRotateY(object Hittable, angle float64) *RotateY {
	sinTheta := math.Sin(DegreesToRadians(angle))
	cosTheta := math.Cos(DegreesToRadians(angle))
	bbox := object.BBOX()

	minPoint, maxPoint := NewVec3(math.Inf(-1), math.Inf(-1), math.Inf(-1)), NewVec3(math.Inf(1), math.Inf(1), math.Inf(1))

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				x := float64(i)*bbox.X.Max + float64(1-i)*bbox.X.Min
				y := float64(j)*bbox.Y.Max + float64(1-j)*bbox.Y.Min
				z := float64(k)*bbox.Z.Max + float64(1-k)*bbox.Z.Min

				newX := cosTheta*x + sinTheta*z
				newZ := -sinTheta*x + cosTheta*z

				tester := NewVec3(newX, y, newZ)

				for c := 0; c < 3; c++ {
					minPoint.SetDim(c, min(minPoint.GetDim(c), tester.GetDim(c)))
					maxPoint.SetDim(c, max(maxPoint.GetDim(c), maxPoint.GetDim(c)))
				}

			}
		}

	}

	return &RotateY{Object: &object, CosTheta: cosTheta, SinTheta: sinTheta, BBOXField: NewAABBFromPoints(minPoint, maxPoint)}

}

func (ro *RotateY) Hit(r *Ray, i Interval, rec *HitRecord) bool {
	origin := NewVec3(ro.CosTheta*r.Origin.X-(ro.SinTheta*r.Origin.Z), r.Origin.Y, ro.SinTheta*r.Origin.X+ro.CosTheta*r.Origin.Z)
	direction := NewVec3(ro.CosTheta*r.Direction.X-(ro.SinTheta*r.Direction.Z), r.Direction.Y, ro.SinTheta*r.Direction.X+ro.CosTheta*r.Direction.Z)

	rotatedR := NewRay(origin, direction, r.Time)
	if !(*ro.Object).Hit(&rotatedR, i, rec) {
		return false
	}

	rec.P = NewVec3(ro.CosTheta*rec.P.X+ro.SinTheta*rec.P.Z, rec.P.Y, -ro.SinTheta*rec.P.X+ro.CosTheta*rec.P.Z)
	rec.Normal = NewVec3(ro.CosTheta*rec.Normal.X+ro.SinTheta*rec.Normal.Z, rec.Normal.Y, -ro.SinTheta*rec.Normal.X+ro.CosTheta*rec.Normal.Z)
	return true

}

func (ro *RotateY) BBOX() *AABB {
	return &ro.BBOXField
}
