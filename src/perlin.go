package main

import (
	"math"
	"math/rand/v2"
)

type PerlinNoise struct {
	PointCount          int
	PermX, PermY, PermZ []int
	RandVecs            []Vec3
}

func NewPerlinNoise() *PerlinNoise {
	count := 256
	n := PerlinNoise{
		PointCount: count,
		RandVecs:   make([]Vec3, count),
		PermX:      make([]int, count),
		PermY:      make([]int, count),
		PermZ:      make([]int, count),
	}
	for i := range count {
		n.RandVecs[i] = (NewBoundedRandomVec(-1, 1)).GetUnitVec()
	}
	n.GeneratePerm(n.PermX)
	n.GeneratePerm(n.PermY)
	n.GeneratePerm(n.PermZ)

	return &n
}

func (n *PerlinNoise) Noise(p Vec3) float64 {
	u, v, w := p.X-math.Floor(p.X), p.Y-math.Floor(p.Y), p.Z-math.Floor(p.Z)
	i, j, k := int(math.Floor(p.X)), int(math.Floor(p.Y)), int(math.Floor(p.Z))
	var c [2][2][2]Vec3

	for di := range 2 {
		for dj := range 2 {
			for dk := range 2 {
				c[di][dj][dk] = n.RandVecs[n.PermX[(i+di)&255]^n.PermY[(j+dj)&255]^n.PermZ[(k+dk)&255]]
			}
		}
	}
	return PerlinInterpolation(c, u, v, w)
}

func (n *PerlinNoise) GeneratePerm(p []int) {
	for i := range n.PointCount {
		p[i] = i
	}
	Permute(p, n.PointCount)
}

func Permute(p []int, n int) {
	for i := n - 1; i > 0; i-- {
		target := rand.IntN(i + 1)
		p[i], p[target] = p[target], p[i]
	}
}

func PerlinInterpolation(c [2][2][2]Vec3, u, v, w float64) float64 {
	uu, vv, ww := u*u*(3-2*u), v*v*(3-2*v), w*w*(3-2*w)
	accum := 0.0
	for i := 0.0; i < 2; i++ {
		for j := 0.0; j < 2; j++ {
			for k := 0.0; k < 2; k++ {
				weightV := NewVec3(u-i, v-j, w-k)
				accum += ((i * uu) + (1-i)*(1-uu)) * ((j * vv) + (1-j)*(1-vv)) * ((k * ww) + (1-k)*(1-ww)) * Dot(&c[int(i)][int(j)][int(k)], &weightV)
			}
		}
	}
	return accum
}

func (n *PerlinNoise) Turbulence(p Vec3, depth int) float64 {
	accum := 0.0
	tempP := NewVec3(p.X, p.Y, p.Z)
	weight := 1.0

	for range depth {
		accum += weight * n.Noise(tempP)
		weight *= 0.5
		tempP.ScaleAssign(2.0)
	}
	return math.Abs(accum)
}
