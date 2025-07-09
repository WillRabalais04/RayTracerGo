package main

import (
	"sort"
)

type BVHNode struct {
	Left, Right *Hittable
	BBOXField   *AABB
}

func NewBVHNode(objects []*Hittable, start, end int) *Hittable {
	bbox := NewEmptyAABB()
	subList := make([]*Hittable, end-start)
	copy(subList, objects[start:end])

	for _, obj := range subList {
		bbox.MergeAABB((*obj).BBOX())
	}
	var leftNode, rightNode Hittable
	span := end - start
	switch span {
	case 1:
		leftNode = *subList[0]
		rightNode = *subList[0]
	case 2:
		leftNode = *subList[0]
		rightNode = *subList[1]
	default:
		sort.Slice(subList, func(i, j int) bool {
			// pick a longest axis to compare
			return BoxCompare(bbox.LongestAxis(), *subList[i], *subList[j])
		})
		mid := span / 2
		leftNode = *NewBVHNode(subList, 0, mid)
		rightNode = *NewBVHNode(subList, mid, span)
	}

	h := Hittable(&BVHNode{
		Left:      &leftNode,
		Right:     &rightNode,
		BBOXField: bbox,
	})
	return &h
}
func NewBVHNodeFromList(list *HittableList) *Hittable {
	return NewBVHNode(list.Objects, 0, len(list.Objects))
}
func (n *BVHNode) Hit(r Ray, i *Interval, rec *HitRecord) bool {
	rayInterval := *i
	if !n.BBOXField.Hit(r, &rayInterval) {
		return false
	}
	hitLeft := (*n.Left).Hit(r, i, rec)
	var hitT float64
	if hitLeft {
		hitT = rec.T
	} else {
		hitT = i.Max
	}
	hitRight := (*n.Right).Hit(r, NewInterval(i.Min, hitT), rec)
	return hitLeft || hitRight
}
func (n *BVHNode) BBOX() *AABB {
	return n.BBOXField
}
