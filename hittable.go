package main

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
