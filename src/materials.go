package main

import (
	"math"
	"math/rand/v2"
)

type NoEmittable struct{}

// any struct with this type will promote this method to be called eg. lambertian.emitted
func (ne *NoEmittable) Emitted(u, v float64, p *Vec3) Vec3 {
	return NewVec3(0, 0, 0)
}

type NoScatter struct{}

func (ns *NoScatter) Scatter(rIn *Ray, rec *HitRecord, attenuation *Vec3, scattered *Ray) bool {
	return false
}

type Material interface {
	Scatter(rIn *Ray, rec *HitRecord, attenuation *Vec3, scattered *Ray) bool
	Emitted(u, v float64, p *Vec3) Vec3
}

type Lambertian struct {
	NoEmittable
	Tex Texture
}

func NewLambertian(albedo Vec3) Lambertian {
	t := NewSolidColor(albedo)
	return Lambertian{Tex: &t}
}
func NewLambertianFromTexture(t Texture) Lambertian {
	return Lambertian{Tex: t}
}
func (l *Lambertian) Scatter(rIn *Ray, rec *HitRecord, attenuation *Vec3, scattered *Ray) bool {
	scatterDirection := rec.Normal.Add(RandomUnitVector())
	if scatterDirection.NearZero() {
		scatterDirection = rec.Normal
	}
	*scattered = NewRay(rec.P, scatterDirection, rIn.Time)
	*attenuation = l.Tex.Value(rec.U, rec.V, &rec.P)
	return true
}

type Metal struct {
	NoEmittable
	Albedo Vec3
	Fuzz   float64
}

func NewMetal(albedo Vec3, fuzz float64) Metal {
	return Metal{Albedo: albedo, Fuzz: min(1, fuzz)}
}
func (m *Metal) Scatter(rIn *Ray, rec *HitRecord, attenuation *Vec3, scattered *Ray) bool {
	reflected := Reflect(&rIn.Direction, &rec.Normal)
	reflected = reflected.GetUnitVec().Add(RandomUnitVector().Scale(m.Fuzz))
	*scattered = NewRay(rec.P, reflected, rIn.Time)
	*attenuation = m.Albedo
	return Dot(&scattered.Direction, &rec.Normal) > 0
}

type Dielectric struct {
	NoEmittable
	RefractionIndex float64
}

func NewDielectric(ri float64) Dielectric {
	return Dielectric{RefractionIndex: ri}
}
func (d *Dielectric) Reflectance(cosine, refractionIndex float64) float64 {
	r0 := math.Pow(((1 - refractionIndex) / (1 + refractionIndex)), 2)
	return r0 + (1-r0)*math.Pow((1-cosine), 5)
}
func (d *Dielectric) Scatter(rIn *Ray, rec *HitRecord, attenuation *Vec3, scattered *Ray) bool {
	*attenuation = NewVec3(1.0, 1.0, 1.0)
	ri := d.RefractionIndex
	if rec.FrontFace {
		ri = (1.0 / ri)
	}

	unitDirection := rIn.Direction.GetUnitVec()
	negatedUnitDirection := unitDirection.Negate()
	cosTheta := min(Dot(&negatedUnitDirection, &rec.Normal), 1.0)
	sinTheta := math.Sqrt(1.0 - cosTheta*cosTheta)
	cannotRefract := ri*sinTheta > 1.0
	var direction Vec3

	if cannotRefract || d.Reflectance(cosTheta, ri) > rand.Float64() {
		direction = Reflect(&unitDirection, &rec.Normal)
	} else {
		direction = Refract(&unitDirection, &rec.Normal, ri)
	}
	*scattered = NewRay(rec.P, direction, rIn.Time)

	return true
}

type DiffuseLight struct {
	NoScatter
	Tex Texture
}

func NewDiffuseLight(emit Vec3) DiffuseLight {
	t := NewSolidColor(emit)
	return DiffuseLight{Tex: &t}
}

func NewDiffuseLightFromTexture(t Texture) DiffuseLight {
	return DiffuseLight{Tex: t}
}

func (d *DiffuseLight) Emitted(u, v float64, p *Vec3) Vec3 {
	return d.Tex.Value(u, v, p)
}
