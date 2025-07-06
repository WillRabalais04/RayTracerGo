package main

import "fmt"

type AABB struct { // axis-aligned bounding box
	X, Y, Z Interval
}

func NewAABB(x, y, z Interval) AABB {
	box := AABB{X: x, Y: y, Z: z}
	box.PadToMinimums()
	return box
}
func NewAABBFromPoints(a, b Vec3) AABB {
	var x, y, z Interval
	if a.X <= b.X {
		x = NewInterval(a.X, b.X)
	} else {
		x = NewInterval(b.X, a.X)
	}
	if a.Y <= b.Y {
		y = NewInterval(a.Y, b.Y)
	} else {
		y = NewInterval(b.Y, a.Y)
	}
	if a.Z <= b.Z {
		z = NewInterval(a.Z, b.Z)
	} else {
		z = NewInterval(b.Z, a.Z)
	}
	return AABB{X: x, Y: y, Z: z}
}
func MergeAABBs(b0, b1 *AABB) AABB {
	return NewAABB(NewEnclosingInterval(&b0.X, &b1.X), NewEnclosingInterval(&b0.Y, &b1.Y), NewEnclosingInterval(&b0.Z, &b1.Z))
}
func (aabb *AABB) axisInterval(n int) *Interval {
	if n == 1 {
		return &aabb.Y
	}
	if n == 2 {
		return &aabb.Z
	}
	return &aabb.X
}
func (aabb *AABB) Hit(r *Ray, i Interval) bool {

	for axis := 0; axis < 3; axis++ {
		ax := aabb.axisInterval(axis)
		adinv := 1.0 / r.Direction.GetDim(axis)
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
		fmt.Printf("boxcomparefailed")
		return false
	}
	// compares min of chosen axis of a and b
	return a.BBOX().axisInterval(axisIndex).Min < b.BBOX().axisInterval(axisIndex).Min
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

func (aabb *AABB) ShiftAABB(offset *Vec3) AABB {
	return NewAABB(aabb.X.ShiftInterval(offset.X), aabb.Y.ShiftInterval(offset.Y), aabb.Z.ShiftInterval(offset.Z))
}
