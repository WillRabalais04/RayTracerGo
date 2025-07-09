package main

type AABB struct { // axis-aligned bounding box
	X, Y, Z *Interval
}

func NewAABB(x, y, z *Interval) *AABB {
	box := AABB{X: x, Y: y, Z: z}
	box.PadToMinimums()
	return &box
}
func NewAABBFromPoints(a, b Vec3) *AABB {
	return &AABB{
		X: NewInterval(min(a.X, b.X), max(a.X, b.X)),
		Y: NewInterval(min(a.Y, b.Y), max(a.Y, b.Y)),
		Z: NewInterval(min(a.Z, b.Z), max(a.Z, b.Z)),
	}
}
func NewEmptyAABB() *AABB {
	return &AABB{X: EmptyInterval, Y: EmptyInterval, Z: EmptyInterval}

}
func MergedAABBs(b0, b1 *AABB) *AABB {
	if b0 == nil {
		b0 = NewEmptyAABB()
	}
	if b1 == nil {
		b1 = NewEmptyAABB()
	}
	return NewAABB(NewEnclosingInterval(b0.X, b1.X), NewEnclosingInterval(b0.Y, b1.Y), NewEnclosingInterval(b0.Z, b1.Z))
}
func (b0 *AABB) MergeAABB(b1 *AABB) {
	if b0 == nil {
		b0 = NewEmptyAABB()
	}
	if b1 == nil {
		b1 = NewEmptyAABB()
	}
	b0.X = NewEnclosingInterval(b0.X, b1.X)
	b0.Y = NewEnclosingInterval(b0.Y, b1.Y)
	b0.Z = NewEnclosingInterval(b0.Z, b1.Z)
}

func (aabb *AABB) AxisInterval(n int) *Interval {
	if n == 1 {
		return aabb.Y
	}
	if n == 2 {
		return aabb.Z
	}
	return aabb.X
}
func (aabb *AABB) Hit(r Ray, i *Interval) bool {

	for axis := range 3 {
		ax := aabb.AxisInterval(axis)
		adinv := 1 / r.Direction.GetDim(axis)
		t0 := (ax.Min - r.Origin.GetDim(axis)) * adinv
		t1 := (ax.Max - r.Origin.GetDim(axis)) * adinv

		if t0 < t1 {
			i.Min = max(t0, i.Min)
			i.Max = min(t1, i.Max)
		} else {
			i.Min = max(t1, i.Min)
			i.Max = min(t0, i.Max)
		}

		if i.Max <= i.Min {
			return false
		}
	}
	return true
}
func (aabb *AABB) LongestAxis() int {
	return ArgMax(aabb.X.Size(), aabb.Y.Size(), aabb.Z.Size())
}

func BoxCompare(axisIndex int, a, b Hittable) bool {

	if a.BBOX() == nil || b.BBOX() == nil {
		panic("boxcomparefailed: nil BBOX")
	}
	// compares min of chosen axis of a and b
	return a.BBOX().AxisInterval(axisIndex).Min < b.BBOX().AxisInterval(axisIndex).Min
}
func (aabb *AABB) PadToMinimums() {
	delta := 0.0001
	if aabb.X.Size() < delta {
		aabb.X.Expand(delta)
	}
	if aabb.Y.Size() < delta {
		aabb.Y.Expand(delta)
	}
	if aabb.Z.Size() < delta {
		aabb.Z.Expand(delta)
	}
}

func (aabb *AABB) ShiftAABB(offset Vec3) {
	aabb.X.ShiftInterval(offset.X)
	aabb.Y.ShiftInterval(offset.Y)
	aabb.Z.ShiftInterval(offset.Z)
}
