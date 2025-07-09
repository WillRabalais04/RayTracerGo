package main

import "math"

type Translate struct {
	Offset    Vec3
	Object    *Hittable
	BBOXField *AABB
}

func NewTranslateY(object *Hittable, offset Vec3) *Hittable {
	(*object).BBOX().ShiftAABB(offset)
	t := Hittable(&Translate{Offset: offset, Object: object, BBOXField: (*object).BBOX()})
	return &t
}

func (t *Translate) Hit(r Ray, i *Interval, rec *HitRecord) bool {
	offsetRay := NewRay(r.Origin.Sub(t.Offset), r.Direction, r.Time)
	if !(*t.Object).Hit(offsetRay, i, rec) {
		return false
	}
	rec.P.PlusEq(t.Offset)
	return true
}

func (t *Translate) BBOX() *AABB {
	return t.BBOXField
}

type RotateY struct {
	Angle     float64
	Object    *Hittable
	BBOXField *AABB
}

func NewRotateY(object *Hittable, angle float64) *Hittable {
	sinTheta := math.Sin(DegreesToRadians(angle))
	cosTheta := math.Cos(DegreesToRadians(angle))
	bbox := (*object).BBOX()

	minPoint, maxPoint := NewVec3(math.Inf(-1), math.Inf(-1), math.Inf(-1)), NewVec3(math.Inf(1), math.Inf(1), math.Inf(1))

	for i := range 2 {
		for j := range 2 {
			for k := range 2 {
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
	r := Hittable(&RotateY{Object: object, Angle: angle, BBOXField: NewAABBFromPoints(minPoint, maxPoint)})
	return &r
}

func (ro *RotateY) Hit(r Ray, i *Interval, rec *HitRecord) bool {
	sinTheta, cosTheta := math.Sin(DegreesToRadians(ro.Angle)), math.Cos(DegreesToRadians(ro.Angle))
	origin := NewVec3(cosTheta*r.Origin.X-(sinTheta*r.Origin.Z), r.Origin.Y, sinTheta*r.Origin.X+cosTheta*r.Origin.Z)
	direction := NewVec3(cosTheta*r.Direction.X-(sinTheta*r.Direction.Z), r.Direction.Y, sinTheta*r.Direction.X+cosTheta*r.Direction.Z)

	rotatedR := NewRay(origin, direction, r.Time)
	if !(*ro.Object).Hit(rotatedR, i, rec) {
		return false
	}

	rec.P = NewVec3(cosTheta*rec.P.X+sinTheta*rec.P.Z, rec.P.Y, -sinTheta*rec.P.X+cosTheta*rec.P.Z)
	rec.Normal = NewVec3(cosTheta*rec.Normal.X+sinTheta*rec.Normal.Z, rec.Normal.Y, -sinTheta*rec.Normal.X+cosTheta*rec.Normal.Z)
	return true
}

func (ro *RotateY) BBOX() *AABB {
	return ro.BBOXField
}
