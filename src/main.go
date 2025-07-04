package main

import (
	"math"
	"math/rand/v2"
)

func main() {

	cam := NewCamera()
	cam.CamConfig7()
	world := World7()
	cam.Render(&world)

}

func World1() HittableList {
	materialGround := NewLambertian(NewVec3(0.8, 0.8, 0.0))
	materialCenter := NewLambertian(NewVec3(0.1, 0.2, 0.5))
	materialLeft := NewDielectric(1.50)
	materialBubble := NewDielectric(1.0 / 1.50)
	materialRight := NewMetal(NewVec3(0.8, 0.6, 0.2), 1.0)

	s1 := NewSphere(NewVec3(0.0, -100.5, -1.0), 100, &materialGround)
	s2 := NewSphere(NewVec3(0.0, 0.0, -1.2), 0.5, &materialCenter)
	s3 := NewSphere(NewVec3(-1.0, 0.0, -1.0), 0.5, &materialLeft)
	s4 := NewSphere(NewVec3(-1.0, 0.0, -1.0), 0.4, &materialBubble)
	s5 := NewSphere(NewVec3(1.0, 0.0, -1.0), 0.5, &materialRight)

	return *NewHittableList(&s1, &s2, &s3, &s4, &s5)
}
func World2() HittableList {
	R := math.Cos(math.Pi / 4)
	materialLeft := NewLambertian(NewVec3(0.0, 0.0, 1.0))
	materialRight := NewLambertian(NewVec3(1.0, 0.0, 0.0))

	s1 := NewSphere(NewVec3(-R, 0.0, -1.0), R, &materialLeft)
	s2 := NewSphere(NewVec3(R, 0.0, -1.0), R, &materialRight)

	return *NewHittableList(&s1, &s2)
}
func World3() HittableList {
	checker := NewCheckeredTextureFromColors(0.32, NewVec3(0.2, 0.3, 0.1), NewVec3(0.9, 0.9, 0.9))
	groundMaterial := NewLambertianFromTexture(&checker)
	s1 := NewSphere(NewVec3(0, -1000, -0), 1000, &groundMaterial)

	world := NewHittableList(&s1)

	for i := -11; i < 11; i++ {
		for j := -11; j < 11; j++ {
			chooseMat := rand.Float64()
			center := NewVec3(float64(i)+0.9*rand.Float64(), 0.2, float64(j)+0.9*rand.Float64())

			if center.Sub(NewVec3(4, 0.2, 0.0)).Length() > 0.9 {
				if chooseMat < 0.8 {
					albedo := NewVec3(rand.Float64(), rand.Float64(), rand.Float64()).Mul(NewVec3(rand.Float64(), rand.Float64(), rand.Float64()))
					center2 := center.Add(NewVec3(0.0, rand.Float64()*0.5, 0.0))
					sM := NewLambertian(albedo)
					sN := NewMovingSphere(center, center2, 0.2, &sM)
					world.Add(&sN)
				} else if chooseMat < 0.95 {
					albedo := NewVec3(rand.Float64()*0.5+0.5, rand.Float64()*0.5+0.5, rand.Float64()*0.5+0.5)
					fuzz := rand.Float64() * 0.5
					sM := NewMetal(albedo, fuzz)
					sN := NewSphere(center, 0.2, &sM)
					world.Add(&sN)
				} else {
					sM := NewDielectric(1.5)
					sN := NewSphere(center, 0.2, &sM)
					world.Add(&sN)
				}
			}
		}
	}

	mat1 := NewDielectric(1.5)
	s2 := NewSphere(NewVec3(0.0, 1.0, 0.0), 1.0, &mat1)
	world.Add(&s2)

	mat2 := NewLambertian(NewVec3(0.4, 0.2, 0.1))
	s3 := NewSphere(NewVec3(-4.0, 1.0, 0.0), 1.0, &mat2)
	world.Add(&s3)

	mat3 := NewMetal(NewVec3(0.7, 0.6, 0.5), 0.0)
	s4 := NewSphere(NewVec3(4.0, 1.0, 0.0), 1.0, &mat3)
	world.Add(&s4)
	return *NewHittableList(NewBVHNodeFromList(*world))
}

func World4() HittableList {
	checker := NewCheckeredTextureFromColors(0.32, NewVec3(0.2, 0.3, 0.1), NewVec3(0.9, 0.9, 0.9))
	groundMaterial := NewLambertianFromTexture(&checker)
	s1 := NewSphere(NewVec3(0, -10, 0), 10, &groundMaterial)
	s2 := NewSphere(NewVec3(0, 10, 0), 10, &groundMaterial)

	return *NewHittableList(&s1, &s2)
}

func World5() HittableList {
	earthTexture := NewImageTexture("../textures/earthmap.jpg")
	earthSurface := NewLambertianFromTexture(earthTexture)
	globe := NewSphere(NewVec3(0, 0, 0), 2, &earthSurface)
	return *NewHittableList(&globe)
}

func World6() HittableList { // perlin noise
	perlinTexture := NewNoiseTexture(4)
	perlinMaterial := NewLambertianFromTexture(&perlinTexture)
	s1 := NewSphere(NewVec3(0, -1000, 0), 1000, &perlinMaterial)
	s2 := NewSphere(NewVec3(0, 2, 0), 2, &perlinMaterial)

	return *NewHittableList(&s1, &s2)

}
func World7() HittableList { // perlin noise
	leftRed := NewLambertian(NewVec3(1.0, 0.2, 0.2))
	backGreen := NewLambertian(NewVec3(0.2, 1.0, 0.2))
	rightBlue := NewLambertian(NewVec3(0.2, 0.2, 1.0))
	upperOrange := NewLambertian(NewVec3(1.0, 0.5, 0.0))
	lowerTeal := NewLambertian(NewVec3(0.2, 0.8, 0.8))

	q1 := NewQuad(NewVec3(-3, -2, 5), NewVec3(0, 0, -4), NewVec3(0, 4, 0), &leftRed)
	q2 := NewQuad(NewVec3(-2, -2, 0), NewVec3(4, 0, 0), NewVec3(0, 4, 0), &backGreen)
	q3 := NewQuad(NewVec3(3, -2, 1), NewVec3(0, 0, 4), NewVec3(0, 4, 0), &rightBlue)
	q4 := NewQuad(NewVec3(-2, 3, 1), NewVec3(4, 0, 0), NewVec3(0, 0, 4), &upperOrange)
	q5 := NewQuad(NewVec3(-2, -3, 5), NewVec3(4, 0, 0), NewVec3(0, 0, -4), &lowerTeal)

	return *NewHittableList(&q1, &q2, &q3, &q4, &q5)

}
func (c *Camera) CamConfig1() {
	c.AspectRatio = 16.0 / 9.0
	c.ImageWidth = 400
	c.MaxDepth = 50
	c.SamplesPerPixel = 50

	c.VFov = 40
	c.lookFrom = NewVec3(2.0, 2.0, 1.0)
	c.lookAt = NewVec3(0.0, 0.0, -1.0)
	c.vup = NewVec3(0.0, 1.0, 0.0)

	c.defocusAngle = 0.0
	c.focusDistance = 5.0
}
func (c *Camera) CamConfig2() {
	c.AspectRatio = 16.0 / 9.0
	c.ImageWidth = 400
	c.MaxDepth = 50
	c.SamplesPerPixel = 50

	c.VFov = 20
	c.lookFrom = NewVec3(0.0, 0.0, 5.0)
	c.lookAt = NewVec3(0.0, 0.0, 0.0)
	c.vup = NewVec3(0.0, 1.0, 0.0)

	c.defocusAngle = 0.6
	c.focusDistance = 10.0
}

func (c *Camera) CamConfig3() {
	c.AspectRatio = 16.0 / 9.0
	c.ImageWidth = 400
	c.MaxDepth = 50
	c.SamplesPerPixel = 50

	c.VFov = 20
	c.lookFrom = NewVec3(13.0, 2.0, 3.0)
	c.lookAt = NewVec3(0.0, 0.0, 0.0)
	c.vup = NewVec3(0.0, 1.0, 0.0)

	c.defocusAngle = 0.6
	c.focusDistance = 10.0
}

func (c *Camera) CamConfig4() {
	c.AspectRatio = 16.0 / 9.0
	c.ImageWidth = 400
	c.MaxDepth = 50
	c.SamplesPerPixel = 100

	c.VFov = 20
	c.lookFrom = NewVec3(13.0, 2.0, 3.0)
	c.lookAt = NewVec3(0.0, 0.0, 0.0)
	c.vup = NewVec3(0.0, 1.0, 0.0)

	c.defocusAngle = 0.0
	// c.focusDistance = 10.0
}

func (c *Camera) CamConfig5() {
	c.AspectRatio = 16.0 / 9.0
	c.ImageWidth = 400
	c.MaxDepth = 50
	c.SamplesPerPixel = 100

	c.VFov = 15
	c.lookFrom = NewVec3(12, 0, -12)
	c.lookAt = NewVec3(0, 0, 0)
	c.vup = NewVec3(0, 1, 0)

	c.defocusAngle = 0.0
}

func (c *Camera) CamConfig6() {
	c.AspectRatio = 16.0 / 9.0
	c.ImageWidth = 400
	c.MaxDepth = 50
	c.SamplesPerPixel = 100

	c.VFov = 20
	c.lookFrom = NewVec3(13, 2, 3)
	c.lookAt = NewVec3(0, 2, 0)
	c.vup = NewVec3(0, 1, 0)

	c.defocusAngle = 0.0
}

func (c *Camera) CamConfig7() {
	c.AspectRatio = 1.0
	c.ImageWidth = 400
	c.MaxDepth = 50
	c.SamplesPerPixel = 100

	c.VFov = 80
	c.lookFrom = NewVec3(0, 0, 9)
	c.lookAt = NewVec3(0, 0, 0)
	c.vup = NewVec3(0, 1, 0)

	c.defocusAngle = 0.0
	// c.focusDistance = 10.0
}
