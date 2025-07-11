package main

import (
	"math"
	"math/rand/v2"
)

type Material interface {
	Scatter(rIn Ray, rec *HitRecord, attenuation *Vec3, scattered *Ray) bool
	Emitted(u, v float64, p Vec3) Vec3
}

type NoEmittable struct{} // any struct with this type will promote this method to be called eg. lambertian.emitted
func (ne *NoEmittable) Emitted(u, v float64, p Vec3) Vec3 {
	return NewVec3(0, 0, 0)
}

type NoScatter struct{}

func (ns *NoScatter) Scatter(rIn Ray, rec *HitRecord, attenuation *Vec3, scattered *Ray) bool {
	return false
}

type Lambertian struct {
	Tex *Texture
	NoEmittable
}

func NewLambertian(albedo Vec3) *Material {
	m := Material(&Lambertian{Tex: NewSolidColor(albedo)})
	return &m
}
func NewLambertianFromTexture(t *Texture) *Material {
	m := Material(&Lambertian{Tex: t})
	return &m
}
func (l *Lambertian) Scatter(rIn Ray, rec *HitRecord, attenuation *Vec3, scattered *Ray) bool {
	scatterDirection := rec.Normal.Add(RandomUnitVector())
	if scatterDirection.NearZero() {
		scatterDirection = rec.Normal
	}
	*scattered = NewRay(rec.P, scatterDirection, rIn.Time)
	*attenuation = (*l.Tex).Value(rec.U, rec.V, rec.P)
	return true
}

type Metal struct {
	Albedo Vec3
	Fuzz   float64
	NoEmittable
}

func NewMetal(albedo Vec3, fuzz float64) *Material {
	m := Material(&Metal{Albedo: albedo, Fuzz: min(1, fuzz)})
	return &m
}
func (m *Metal) Scatter(rIn Ray, rec *HitRecord, attenuation *Vec3, scattered *Ray) bool {
	reflected := Reflect(&rIn.Direction, &rec.Normal)
	reflected = reflected.GetUnitVec().Add(RandomUnitVector().Scale(m.Fuzz))
	*scattered = NewRay(rec.P, reflected, rIn.Time)
	*attenuation = m.Albedo
	return Dot(&scattered.Direction, &rec.Normal) > 0
}

type Dielectric struct {
	RefractionIndex float64
	NoEmittable
}

func NewDielectric(ri float64) *Material {
	m := Material(&Dielectric{RefractionIndex: ri})
	return &m
}
func (d *Dielectric) Reflectance(cosine, refractionIndex float64) float64 {
	r0 := math.Pow(((1 - refractionIndex) / (1 + refractionIndex)), 2)
	return r0 + (1-r0)*math.Pow((1-cosine), 5)
}
func (d *Dielectric) Scatter(rIn Ray, rec *HitRecord, attenuation *Vec3, scattered *Ray) bool {
	*attenuation = NewVec3(1, 1, 1)
	ri := d.RefractionIndex
	if rec.FrontFace {
		ri = (1 / ri)
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
	Tex *Texture
	NoScatter
}

func NewDiffuseLight(t *Texture) *Material {
	m := Material(&DiffuseLight{Tex: t})
	return &m
}
func NewColoredDiffuseLight(emit Vec3) *Material {
	m := Material(&DiffuseLight{Tex: NewSolidColor(emit)})
	return &m
}

func (d *DiffuseLight) Emitted(u, v float64, p Vec3) Vec3 {
	return (*d.Tex).Value(u, v, p)
}

type Isotropic struct {
	Tex *Texture
	NoEmittable
}

func NewIsotropic(albedo Vec3) *Material {
	m := Material(&Isotropic{Tex: NewSolidColor(albedo)})
	return &m
}
func NewIsotropicFromTexture(t *Texture) *Material {
	m := Material(&Isotropic{Tex: t})
	return &m
}

func (i Isotropic) Scatter(rIn Ray, rec *HitRecord, attenuation *Vec3, scattered *Ray) bool {
	*scattered = NewRay(rec.P, RandomUnitVector(), rIn.Time)
	*attenuation = (*i.Tex).Value(rec.U, rec.V, rec.P)
	return true
}
