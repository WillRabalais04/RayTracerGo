package main

import "sort"

type BVHNode struct {
	Left, Right Hittable
	BBOXField   AABB
} // pointer to BVHNode implements hittable but BVHNode itself does not

func NewBVHNode(objects []Hittable, start, end int) *BVHNode {
	bbox := NewAABB(EmptyInterval, EmptyInterval, EmptyInterval)
	subList := make([]Hittable, end-start)
	copy(subList, objects[start:end])

	for objectIndex := start; objectIndex < end; objectIndex++ {
		bbox = MergeAABBs(&bbox, objects[objectIndex].BBOX())
	}

	var leftNode, rightNode Hittable
	span := end - start
	switch span {
	case 1:
		leftNode = objects[start]
		rightNode = objects[start]
	case 2:
		leftNode = objects[start]
		rightNode = objects[start+1]
	default:
		sort.Slice(subList, func(i, j int) bool {
			// pick a longest axis to compare
			return BoxCompare(bbox.LongestAxis(), subList[i], subList[j])
		})
		mid := start + span/2
		leftNode = NewBVHNode(objects, start, mid)
		rightNode = NewBVHNode(objects, mid, end)
	}
	return &BVHNode{
		Left:      leftNode,
		Right:     rightNode,
		BBOXField: bbox,
	}
}
func NewBVHNodeFromList(list HittableList) *BVHNode {
	return NewBVHNode(list.Objects, 0, len(list.Objects))
}
func (n *BVHNode) Hit(r *Ray, i Interval, rec *HitRecord) bool {
	if !n.BBOXField.Hit(r, i) {
		return false
	}
	hitLeft := n.Left.Hit(r, i, rec)
	var hitT float64
	if hitLeft {
		hitT = rec.T
	} else {
		hitT = i.Max
	}
	hitRight := n.Right.Hit(r, NewInterval(i.Min, hitT), rec)

	return hitLeft || hitRight
}
func (n *BVHNode) BBOX() *AABB {
	return &n.BBOXField
}
