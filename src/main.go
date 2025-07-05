package main

import (
	"math"
	"math/rand/v2"
)

func main() {

	cam := NewCamera()
	cam.CamConfig9()
	world := World9()
	cam.Render(&world)

}

func World1() HittableList { // 3 spheres
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
func World2() HittableList { // red & blue spheres
	R := math.Cos(math.Pi / 4)
	materialLeft := NewLambertian(NewVec3(0.0, 0.0, 1.0))
	materialRight := NewLambertian(NewVec3(1.0, 0.0, 0.0))

	s1 := NewSphere(NewVec3(-R, 0.0, -1.0), R, &materialLeft)
	s2 := NewSphere(NewVec3(R, 0.0, -1.0), R, &materialRight)

	return *NewHittableList(&s1, &s2)
}
func World3() HittableList { // bouncing spheres
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

func World4() HittableList { // checkered spheres
	checker := NewCheckeredTextureFromColors(0.32, NewVec3(0.2, 0.3, 0.1), NewVec3(0.9, 0.9, 0.9))
	groundMaterial := NewLambertianFromTexture(&checker)
	s1 := NewSphere(NewVec3(0, -10, 0), 10, &groundMaterial)
	s2 := NewSphere(NewVec3(0, 10, 0), 10, &groundMaterial)

	return *NewHittableList(&s1, &s2)
}

func World5() HittableList { // texture
	earthTexture := NewImageTexture("../textures/earthmap.jpg")
	earthSurface := NewLambertianFromTexture(earthTexture)
	globe := NewSphere(NewVec3(0, 0, 0), 2, &earthSurface)
	return *NewHittableList(&globe)
}

func World6() HittableList { // perlin noise marble
	perlinTexture := NewNoiseTexture(4)
	perlinMaterial := NewLambertianFromTexture(&perlinTexture)
	s1 := NewSphere(NewVec3(0, -1000, 0), 1000, &perlinMaterial)
	s2 := NewSphere(NewVec3(0, 2, 0), 2, &perlinMaterial)

	return *NewHittableList(&s1, &s2)

}
func World7() HittableList { // quads
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

func World8() HittableList { // purple marble

	perlinTexture := NewNoiseTexture(3)
	perlinMaterial := NewLambertianFromTexture(&perlinTexture)
	diffuseLightRed := NewDiffuseLight(NewVec3(0, 0, 255))
	diffuseLightBlue := NewDiffuseLight(NewVec3(255, 0, 0))
	// diffuseLightWhite := NewDiffuseLight(NewVec3(4, 4, 4))

	s1 := NewSphere(NewVec3(0, -1000, 0), 1000, &perlinMaterial)
	s7 := NewSphere(NewVec3(0, 1012, 0), 1000, &perlinMaterial)
	s2 := NewSphere(NewVec3(0, 6, 0), 3, &perlinMaterial)
	s3 := NewSphere(NewVec3(0, 1, 0), 1, &diffuseLightRed)
	s4 := NewSphere(NewVec3(0, 11, 0), 1, &diffuseLightBlue)
	// s5 := NewSphere(NewVec3(500, 6, 500), 300, &perlinMaterial)

	// q1 := NewQuad(NewVec3(2, 0, 2), NewVec3(-2, 0, -2), NewVec3(0, 2, 0), &diffuseLightRed)
	// q2 := NewQuad(NewVec3(-2, 0, -2), NewVec3(2, 0, 2), NewVec3(0, 2, 0), &diffuseLightBlue)
	return *NewHittableList(&s2, &s3, &s4, &s7, &s1)

}

func World9() HittableList { // cornell box
	red := NewLambertian(NewVec3(0.65, 0.05, 0.05))
	white := NewLambertian(NewVec3(0.73, 0.73, 0.73))
	green := NewLambertian(NewVec3(0.12, 0.45, 0.15))
	light := NewDiffuseLight(NewVec3(15, 15, 15))

	q1 := NewQuad(NewVec3(555, 0, 0), NewVec3(0, 555, 0), NewVec3(0, 0, 555), &green)
	q2 := NewQuad(NewVec3(0, 0, 0), NewVec3(0, 555, 0), NewVec3(0, 0, 555), &red)
	q3 := NewQuad(NewVec3(343, 554, 332), NewVec3(-130, 0, 0), NewVec3(0, 0, -105), &light)
	q4 := NewQuad(NewVec3(0, 0, 0), NewVec3(555, 0, 0), NewVec3(0, 0, 555), &white)
	q5 := NewQuad(NewVec3(555, 555, 555), NewVec3(-555, 0, 0), NewVec3(0, 0, -555), &white)
	q6 := NewQuad(NewVec3(0, 0, 555), NewVec3(555, 0, 0), NewVec3(0, 555, 0), &white)

	world := NewHittableList(&q1, &q2, &q3, &q4, &q5, &q6)

	// world.Add(NewBox(NewVec3(130, 0, 65), NewVec3(295, 165, 230), &white))
	// world.Add(NewBox(NewVec3(265, 0, 295), NewVec3(430, 330, 460), &white))
	translateVector1 := NewVec3(265, 0, 295)
	ogBox1 := NewBox(NewVec3(0, 0, 0), NewVec3(165, 330, 165), &white)
	rotatedBox1 := NewRotateY(ogBox1, 15)
	box1 := NewTranslateY(translateVector1, rotatedBox1)
	world.Add(box1)

	translateVector2 := NewVec3(130, 0, 65)
	ogBox2 := NewBox(NewVec3(0, 0, 0), NewVec3(165, 165, 165), &white)
	rotatedBox2 := NewRotateY(ogBox2, -18)
	box2 := NewTranslateY(translateVector2, rotatedBox2)
	world.Add(box2)

	return *world

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
	c.MaxDepth = 50
	c.SamplesPerPixel = 100

	c.VFov = 80
	c.LookFrom = NewVec3(0, 0, 9)
	c.LookAt = NewVec3(0, 0, 0)
	c.VUP = NewVec3(0, 1, 0)

	c.DefocusAngle = 0.0
	c.Background = NewVec3(0.7, 0.8, 1)
	// c.focusDistance = 10.0
}

func (c *Camera) CamConfig8() {
	c.AspectRatio = 16.0 / 9.0
	c.ImageWidth = 1000
	c.MaxDepth = 50
	c.SamplesPerPixel = 500

	c.VFov = 25
	c.LookFrom = NewVec3(10, 6, 10)
	c.LookAt = NewVec3(0, 6, 0)
	c.VUP = NewVec3(0, 1, 0)

	c.DefocusAngle = 0.0
	c.Background = NewVec3(0, 0, 0)
	// c.focusDistance = 10.0
}

func (c *Camera) CamConfig9() {
	c.AspectRatio = 1.0
	c.ImageWidth = 600
	c.MaxDepth = 50
	c.SamplesPerPixel = 200

	c.VFov = 40
	c.LookFrom = NewVec3(278, 278, -800)
	c.LookAt = NewVec3(278, 278, 0)
	c.VUP = NewVec3(0, 1, 0)

	c.DefocusAngle = 0.0
	c.Background = NewVec3(0, 0, 0)
	// c.focusDistance = 10.0
}
