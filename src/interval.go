package main

import "math"

type Interval struct {
	Min float64
	Max float64
}

var (
	EmptyInterval    = Interval{Min: math.Inf(1), Max: math.Inf(-1)}
	UniverseInterval = Interval{Min: math.Inf(-1), Max: math.Inf(1)}
)

func NewInterval(min, max float64) Interval {
	return Interval{Min: min, Max: max}
}
func (i *Interval) Size() float64 {
	return i.Max - i.Min
}
func (i *Interval) Contains(x float64) bool {
	return i.Min <= x && x <= i.Max
}
func (i *Interval) Surrounds(x float64) bool {
	return i.Min < x && x < i.Max
}
func (i *Interval) Clamp(x float64) float64 {
	if x < i.Min {
		return i.Min
	}
	if x > i.Max {
		return i.Max
	}
	return x
}
func (i *Interval) Expand(delta float64) {
	i.Min -= delta / 2
	i.Max += delta / 2
}
func NewEnclosingInterval(a, b *Interval) Interval {
	return Interval{Min: min(a.Min, b.Min), Max: max(a.Max, b.Max)}
}

func (i *Interval) ShiftInterval(shiftAmt float64) Interval {
	return NewInterval(i.Min+shiftAmt, i.Max+shiftAmt)
}
