package main

type Ray struct {
	Origin    Vec3
	Direction Vec3
	Time      float64
}

func (r Ray) at(t float64) Vec3 {
	return r.Origin.Add(r.Direction.Scale(t))
}
func NewRay(o, d Vec3, t float64) Ray {
	return Ray{Origin: o, Direction: d, Time: t}
}
