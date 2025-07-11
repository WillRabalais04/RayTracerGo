package main

import (
	"math"
	"math/rand/v2"
)

func main() {
	cam := NewCamera()
	cam.CamConfig10()
	cam.Render(World10())
}

func World1() *HittableList { // 3 spheres
	materialGround := NewLambertian(NewVec3(0.8, 0.8, 0))
	materialCenter := NewLambertian(NewVec3(0.1, 0.2, 0.5))
	materialLeft := NewDielectric(1.50)
	materialBubble := NewDielectric(1 / 1.50)
	materialRight := NewMetal(NewVec3(0.8, 0.6, 0.2), 1)

	s1 := NewSphere(NewVec3(0, -100.5, -1), 100, materialGround)
	s2 := NewSphere(NewVec3(0, 0, -1.2), 0.5, materialCenter)
	s3 := NewSphere(NewVec3(-1, 0, -1), 0.5, materialLeft)
	s4 := NewSphere(NewVec3(-1, 0, -1), 0.4, materialBubble)
	s5 := NewSphere(NewVec3(1, 0, -1), 0.5, materialRight)

	return NewHittableList(s1, s2, s3, s4, s5)
}
func World2() *HittableList { // red & blue spheres
	R := math.Cos(math.Pi / 4)
	materialLeft := NewLambertian(NewVec3(0, 0, 1))
	materialRight := NewLambertian(NewVec3(1, 0, 0))

	s1 := NewSphere(NewVec3(-R, 0, -1), R, materialLeft)
	s2 := NewSphere(NewVec3(R, 0, -1), R, materialRight)

	return NewHittableList(s1, s2)
}

func World3() *HittableList { // bouncing spheres
	groundMaterial := NewLambertianFromTexture(NewCheckeredTexture(0.32, NewVec3(0.2, 0.3, 0.1), NewVec3(0.9, 0.9, 0.9)))
	s1 := NewSphere(NewVec3(0, -1000, -0), 1000, groundMaterial)

	world := NewHittableList(s1)

	for i := -11; i < 11; i++ {
		for j := -11; j < 11; j++ {
			chooseMat := rand.Float64()
			center := NewVec3(float64(i)+0.9*rand.Float64(), 0.2, float64(j)+0.9*rand.Float64())

			if center.Sub(NewVec3(4, 0.2, 0)).Length() > 0.9 {
				if chooseMat < 0.8 {
					albedo := NewVec3(rand.Float64(), rand.Float64(), rand.Float64()).Mul(NewVec3(rand.Float64(), rand.Float64(), rand.Float64()))
					center2 := center.Add(NewVec3(0, RandFloatInRange(0, 0.5), 0))
					sM := NewLambertian(albedo)
					sN := NewMovingSphere(center, center2, 0.2, sM)
					world.Add(sN)
				} else if chooseMat < 0.95 {
					albedo := NewBoundedRandomVec(0.5, 1)
					fuzz := RandFloatInRange(0, 0.5)
					sM := NewMetal(albedo, fuzz)
					sN := NewSphere(center, 0.2, sM)
					world.Add(sN)
				} else {
					sM := NewDielectric(1.5)
					sN := NewSphere(center, 0.2, sM)
					world.Add(sN)
				}
			}
		}
	}

	mat1 := NewDielectric(1.5)
	s2 := NewSphere(NewVec3(0, 1, 0), 1, mat1)
	world.Add(s2)

	mat2 := NewLambertian(NewVec3(0.4, 0.2, 0.1))
	s3 := NewSphere(NewVec3(-4, 1, 0), 1, mat2)
	world.Add(s3)

	mat3 := NewMetal(NewVec3(0.7, 0.6, 0.5), 0)
	s4 := NewSphere(NewVec3(4, 1, 0), 1, mat3)
	world.Add(s4)
	return NewHittableList(NewBVHNodeFromList(world))
}

func World4() *HittableList { // checkered spheres
	groundMaterial := NewLambertianFromTexture(NewCheckeredTexture(0.32, NewVec3(0.2, 0.3, 0.1), NewVec3(0.9, 0.9, 0.9)))
	s1 := NewSphere(NewVec3(0, -10, 0), 10, groundMaterial)
	s2 := NewSphere(NewVec3(0, 10, 0), 10, groundMaterial)

	return NewHittableList(s1, s2)
}

func World5() *HittableList { // texture
	earthSurface := NewLambertianFromTexture(NewImageTexture("../textures/earthmap.jpg"))
	globe := NewSphere(NewVec3(0, 0, 0), 2, earthSurface)
	return NewHittableList(globe)
}

func World6() *HittableList { // perlin noise marble
	perlinMaterial := NewLambertianFromTexture(NewNoiseTexture(4))
	s1 := NewSphere(NewVec3(0, -1000, 0), 1000, perlinMaterial)
	s2 := NewSphere(NewVec3(0, 2, 0), 2, perlinMaterial)

	return NewHittableList(s1, s2)
}

func World7() *HittableList { // quads
	leftRed := NewLambertian(NewVec3(1.0, 0.2, 0.2))
	backGreen := NewLambertian(NewVec3(0.2, 1.0, 0.2))
	rightBlue := NewLambertian(NewVec3(0.2, 0.2, 1.0))
	upperOrange := NewLambertian(NewVec3(1.0, 0.5, 0.0))
	lowerTeal := NewLambertian(NewVec3(0.2, 0.8, 0.8))

	q1 := NewQuad(NewVec3(-3, -2, 5), NewVec3(0, 0, -4), NewVec3(0, 4, 0), leftRed)
	q2 := NewQuad(NewVec3(-2, -2, 0), NewVec3(4, 0, 0), NewVec3(0, 4, 0), backGreen)
	q3 := NewQuad(NewVec3(3, -2, 1), NewVec3(0, 0, 4), NewVec3(0, 4, 0), rightBlue)
	q4 := NewQuad(NewVec3(-2, 3, 1), NewVec3(4, 0, 0), NewVec3(0, 0, 4), upperOrange)
	q5 := NewQuad(NewVec3(-2, -3, 5), NewVec3(4, 0, 0), NewVec3(0, 0, -4), lowerTeal)

	return NewHittableList(q1, q2, q3, q4, q5)
}

func World8() *HittableList { // purple marble

	perlinMaterial := NewLambertianFromTexture(NewNoiseTexture(3))
	diffuseLightRed := NewColoredDiffuseLight(NewVec3(0, 0, 255))
	diffuseLightBlue := NewColoredDiffuseLight(NewVec3(255, 0, 0))

	s1 := NewSphere(NewVec3(0, -1000, 0), 1000, perlinMaterial)
	s2 := NewSphere(NewVec3(0, 1012, 0), 1000, perlinMaterial)
	s3 := NewSphere(NewVec3(0, 6, 0), 3, perlinMaterial)
	s4 := NewSphere(NewVec3(0, 1, 0), 1, diffuseLightRed)
	s5 := NewSphere(NewVec3(0, 11, 0), 1, diffuseLightBlue)

	return NewHittableList(s1, s2, s3, s4, s5)
}

func World9() *HittableList { // cornell box
	red := NewLambertian(NewVec3(0.65, 0.05, 0.05))
	white := NewLambertian(NewVec3(0.73, 0.73, 0.73))
	green := NewLambertian(NewVec3(0.12, 0.45, 0.15))
	light := NewColoredDiffuseLight(NewVec3(7, 7, 7))

	world := NewHittableList()
	world.Add(NewCornellBox(green, white, red, light))

	box1 := NewTranslateY(NewRotateY(NewBox(NewVec3(0, 0, 0), NewVec3(165, 330, 165), white), 15), NewVec3(265, 0, 295))
	box2 := NewTranslateY(NewRotateY(NewBox(NewVec3(0, 0, 0), NewVec3(165, 165, 165), white), -18), NewVec3(130, 0, 65))

	cm1 := NewConstantMediumFromColor(box1, 0.01, NewVec3(0, 0, 0))
	world.Add(cm1)
	cm2 := NewConstantMediumFromColor(box2, 0.01, NewVec3(1, 1, 1))
	world.Add(cm2)

	return world
}

func World10() *HittableList { // in a weekend final output
	var boxes1 HittableList

	ground := NewLambertian(NewVec3(0.48, 0.83, 0.53))
	boxesPerSide := 20
	for i := range boxesPerSide {
		for j := range boxesPerSide {
			w := 100.0
			x0, y0, z0 := -1000.0+float64(i)*w, 0.0, -1000.0+float64(j)*w
			x1, y1, z1 := x0+w, RandFloatInRange(1, 101), z0+w

			boxes1.Add(NewBox(NewVec3(x0, y0, z0), NewVec3(x1, y1, z1), ground))
		}
	}
	world := NewHittableList(NewBVHNodeFromList(&boxes1))
	light := NewColoredDiffuseLight(NewVec3(7, 7, 7))
	q1 := NewQuad(NewVec3(123, 554, 147), NewVec3(300, 0, 0), NewVec3(0, 0, 265), light)
	world.Add(q1)

	center1 := NewVec3(400, 400, 200)
	center2 := center1.Add(NewVec3(30, 0, 0))

	sphereMaterial1 := NewLambertian(NewVec3(0.7, 0.3, 0.1))
	s1 := NewMovingSphere(center1, center2, 50, sphereMaterial1)
	world.Add(s1)

	sphereMaterial2 := NewDielectric(1.5)
	s2 := NewSphere(NewVec3(260, 150, 45), 50, sphereMaterial2)
	world.Add(s2)

	sphereMaterial3 := NewMetal(NewVec3(0.8, 0.8, 0.9), 1.0)
	s3 := NewSphere(NewVec3(0, 150, 145), 50, sphereMaterial3)
	world.Add(s3)

	boundary := NewSphere(NewVec3(360, 150, 145), 70, sphereMaterial2)
	world.Add(boundary)
	m1 := NewConstantMediumFromColor(boundary, 0.2, (NewVec3(0.2, 0.4, 0.9)))
	world.Add(m1)

	boundary2 := NewSphere(NewVec3(0, 0, 0), 5000, sphereMaterial2)
	m2 := NewConstantMediumFromColor(boundary2, 0.0001, (NewVec3(1, 1, 1)))
	world.Add(m2)

	earthSurface := NewLambertianFromTexture(NewImageTexture("../textures/earthmap.jpg"))
	s4 := NewSphere(NewVec3(400, 200, 400), 100, earthSurface)
	world.Add(s4)

	perlinMaterial := NewLambertianFromTexture(NewNoiseTexture(0.2))
	s5 := NewSphere(NewVec3(220, 280, 300), 80, perlinMaterial)
	world.Add(s5)

	var boxes2 HittableList

	white := NewLambertian(NewVec3(0.73, 0.73, 0.73))
	for range 1000 {
		s := NewSphere(NewBoundedRandomVec(0, 165), 10, white)
		boxes2.Add(s)
	}

	world.Add(NewTranslateY(NewRotateY(NewBVHNodeFromList(&boxes2), 15), NewVec3(-100, 270, 395)))
	return world
}

func (c *Camera) CamConfig1() {
	c.AspectRatio = 16.0 / 9.0
	c.ImageWidth = 400
	c.MaxDepth = 50
	c.SamplesPerPixel = 50

	c.VFov = 40
	c.LookFrom = NewVec3(2.0, 2.0, 1.0)
	c.LookAt = NewVec3(0.0, 0.0, -1.0)
	c.VUP = NewVec3(0.0, 1.0, 0.0)

	c.DefocusAngle = 0.0
	c.FocusDistance = 5.0
	c.Background = NewVec3(0.7, 0.8, 1)
}
func (c *Camera) CamConfig2() {
	c.AspectRatio = 16.0 / 9.0
	c.ImageWidth = 400
	c.MaxDepth = 50
	c.SamplesPerPixel = 50

	c.VFov = 20
	c.LookFrom = NewVec3(0.0, 0.0, 5.0)
	c.LookAt = NewVec3(0.0, 0.0, 0.0)
	c.VUP = NewVec3(0.0, 1.0, 0.0)

	c.DefocusAngle = 0.6
	c.FocusDistance = 10.0
	c.Background = NewVec3(0.7, 0.8, 1)
}

func (c *Camera) CamConfig3() {
	c.AspectRatio = 16.0 / 9.0
	c.ImageWidth = 400
	c.MaxDepth = 50
	c.SamplesPerPixel = 50

	c.VFov = 20
	c.LookFrom = NewVec3(13.0, 2.0, 3.0)
	c.LookAt = NewVec3(0.0, 0.0, 0.0)
	c.VUP = NewVec3(0.0, 1.0, 0.0)

	c.DefocusAngle = 0.6
	c.FocusDistance = 10.0
	c.Background = NewVec3(0.7, 0.8, 1)
}

func (c *Camera) CamConfig4() {
	c.AspectRatio = 16.0 / 9.0
	c.ImageWidth = 400
	c.MaxDepth = 50
	c.SamplesPerPixel = 100

	c.VFov = 20
	c.LookFrom = NewVec3(13.0, 2.0, 3.0)
	c.LookAt = NewVec3(0.0, 0.0, 0.0)
	c.VUP = NewVec3(0.0, 1.0, 0.0)

	c.DefocusAngle = 0.0
	c.Background = NewVec3(0.7, 0.8, 1)
}

func (c *Camera) CamConfig5() {
	c.AspectRatio = 16.0 / 9.0
	c.ImageWidth = 400
	c.MaxDepth = 50
	c.SamplesPerPixel = 100

	c.VFov = 15
	c.LookFrom = NewVec3(12, 0, -12)
	c.LookAt = NewVec3(0, 0, 0)
	c.VUP = NewVec3(0, 1, 0)

	c.DefocusAngle = 0.0
	c.Background = NewVec3(0.7, 0.8, 1)

}

func (c *Camera) CamConfig6() {
	c.AspectRatio = 16.0 / 9.0
	c.ImageWidth = 400
	c.MaxDepth = 50
	c.SamplesPerPixel = 100

	c.VFov = 20
	c.LookFrom = NewVec3(13, 2, 3)
	c.LookAt = NewVec3(0, 2, 0)
	c.VUP = NewVec3(0, 1, 0)

	c.DefocusAngle = 0.0
	c.Background = NewVec3(0.7, 0.8, 1)
}

func (c *Camera) CamConfig7() {
	c.AspectRatio = 1.0
	c.ImageWidth = 400
	c.MaxDepth = 10
	c.SamplesPerPixel = 50

	c.VFov = 80
	c.LookFrom = NewVec3(0, 0, 9)
	c.LookAt = NewVec3(0, 0, 0)
	c.VUP = NewVec3(0, 1, 0)

	c.DefocusAngle = 0.0
	c.Background = NewVec3(0.7, 0.8, 1)
}

func (c *Camera) CamConfig8() {
	c.AspectRatio = 16.0 / 9.0
	c.ImageWidth = 600
	c.MaxDepth = 50
	c.SamplesPerPixel = 100

	c.VFov = 25
	c.LookFrom = NewVec3(10, 6, 10)
	c.LookAt = NewVec3(0, 6, 0)
	c.VUP = NewVec3(0, 1, 0)

	c.DefocusAngle = 0.0
	c.Background = NewVec3(0, 0, 0)
}

func (c *Camera) CamConfig9() {
	c.AspectRatio = 1.0
	c.ImageWidth = 600
	c.MaxDepth = 50
	c.SamplesPerPixel = 50

	c.VFov = 40
	c.LookFrom = NewVec3(278, 278, -800)
	c.LookAt = NewVec3(278, 278, 0)
	c.VUP = NewVec3(0, 1, 0)

	c.DefocusAngle = 0.0
	c.Background = NewVec3(0, 0, 0)
}

func (c *Camera) CamConfig10() {
	c.AspectRatio = 1.0
	c.ImageWidth = 200
	c.MaxDepth = 50
	c.SamplesPerPixel = 1000

	c.VFov = 40
	c.LookFrom = NewVec3(478, 278, -600)
	c.LookAt = NewVec3(278, 278, 0)
	c.VUP = NewVec3(0, 1, 0)

	c.DefocusAngle = 0.0
	c.Background = NewVec3(0, 0, 0)
}
