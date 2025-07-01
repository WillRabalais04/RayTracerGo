package main

import (
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"os"
)

type Vec3 struct {
	X, Y, Z float64
}

// Vec3 Ops
func NewVec3(x, y, z float64) Vec3 {
	return Vec3{X: x, Y: y, Z: z}
}
func (v1 Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{X: v1.X + v2.X, Y: v1.Y + v2.Y, Z: v1.Z + v2.Z}
}
func (v1 Vec3) Sub(v2 Vec3) Vec3 {
	return Vec3{X: v1.X - v2.X, Y: v1.Y - v2.Y, Z: v1.Z - v2.Z}
}
func (v1 Vec3) Mul(v2 Vec3) Vec3 {
	return Vec3{X: v1.X * v2.X, Y: v1.Y * v2.Y, Z: v1.Z * v2.Z}
}
func (v1 Vec3) Div(v2 Vec3) Vec3 {
	return Vec3{X: v1.X / v2.X, Y: v1.Y / v2.Y, Z: v1.Z / v2.Z}
}
func (v1 *Vec3) PlusEq(v2 Vec3) {
	v1.X += v2.X
	v1.Y += v2.Y
	v1.Z += v2.Z
}
func (v1 *Vec3) MinEq(v2 Vec3) {
	v1.X -= v2.X
	v1.Y -= v2.Y
	v1.Z -= v2.Z
}
func (v1 *Vec3) TimesEq(v2 Vec3) {
	v1.X *= v2.X
	v1.Y *= v2.Y
	v1.Z *= v2.Z
}
func (v1 *Vec3) SlashEq(v2 Vec3) {
	v1.X /= v2.X
	v1.Y /= v2.Y
	v1.Z /= v2.Z
}
func (v Vec3) Scale(t float64) Vec3 {
	return Vec3{X: t * v.X, Y: t * v.Y, Z: t * v.Z}
}
func (v *Vec3) ScaleAssign(t float64) {
	v.X *= t
	v.Y *= t
	v.Z *= t
}
func Dot(v1, v2 *Vec3) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}
func Cross(v1, v2 Vec3) Vec3 {
	return Vec3{
		X: (v1.Y*v2.Z - v1.Z*v2.Y),
		Y: (v1.Z*v2.X - v1.X*v2.Z),
		Z: (v1.X*v2.Y - v1.Y*v2.X)}
}
func (v *Vec3) MakeUnitVec() {
	v.ScaleAssign((1.0 / v.Length()))
}
func (v Vec3) GetUnitVec() Vec3 {
	return v.Scale((1.0 / v.Length()))
}
func (v Vec3) Negate() Vec3 {
	return Vec3{X: -v.X, Y: -v.Y, Z: -v.Z}
}
func (v Vec3) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}
func (v Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}
func RandomUnitVector() Vec3 {
	for {
		p := NewBoundedRandomVec(-1, 1)
		lensq := p.LengthSquared()
		if 1e-160 < lensq && lensq <= 1 {
			return p.Scale(1.0 / math.Sqrt(lensq))
		}
	}
}
func (Normal *Vec3) RandomOnHemisphere() Vec3 {
	OnUnitSphere := RandomUnitVector()
	if Dot(&OnUnitSphere, Normal) < 0.0 {
		OnUnitSphere.ScaleAssign(-1.0)
	}
	return OnUnitSphere
}
func RandomInUnitDisk() Vec3 {
	for {
		p := NewVec3(rand.Float64()*2-1, rand.Float64()*2-1, rand.Float64()*2-1)
		if p.LengthSquared() < 1 {
			return p
		}
	}
}
func (v *Vec3) NearZero() bool {
	s := 1e-8
	return (math.Abs(v.X) < s) && math.Abs(v.Y) < s && math.Abs(v.Z) < s
}

func Reflect(v, n *Vec3) Vec3 {
	return v.Sub(n.Scale(2 * Dot(v, n)))
}
func Refract(uv, n *Vec3, etaiOverEtat float64) Vec3 {
	negatedUV := uv.Negate()
	cosTheta := min(Dot(&negatedUV, n), 1.0)
	rOutPerp := (uv.Add(n.Scale(cosTheta))).Scale(etaiOverEtat)
	rOutParallel := n.Scale(-math.Sqrt(math.Abs(1.0 - rOutPerp.LengthSquared())))
	return rOutParallel.Add(rOutPerp)
}
func NewRandomVec() Vec3 {
	return NewVec3(rand.Float64(), rand.Float64(), rand.Float64())
}
func NewBoundedRandomVec(min, max float64) Vec3 {
	r := max - min
	return NewVec3(rand.Float64()*r+min, rand.Float64()*r+min, rand.Float64()*r+min)
}

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
func NewEmptyInterval() Interval {
	return EmptyInterval
}
func NewUniverseInterval(min float64, max float64) Interval {
	return UniverseInterval
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

func WriteColor(color Vec3) {

	intensity := NewInterval(0.0, 0.999)

	r := int(256 * intensity.Clamp(LinearToGamma(color.X)))
	g := int(256 * intensity.Clamp(LinearToGamma(color.Y)))
	b := int(256 * intensity.Clamp(LinearToGamma(color.Z)))

	file, err := os.OpenFile("out.ppm", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Fprintf(file, "%d %d %d\n", r, g, b)
}

type Ray struct {
	Origin    Vec3
	Direction Vec3
}

// Ray ops
func (r Ray) PointAtParameter(t float64) Vec3 {
	return r.Origin.Add(r.Direction.Scale(t))
}
func NewRay(o, d Vec3) Ray {
	return Ray{Origin: o, Direction: d}
}
func (c *Camera) RayColor(r Ray, depth int, world Hittable) Vec3 {
	if depth <= 0 {
		return NewVec3(0.0, 0.0, 0.0)
	}

	var rec HitRecord
	if world.Hit(r, NewInterval(0.001, math.MaxFloat64), &rec) {
		var scattered Ray
		var attenuation Vec3

		if rec.MaterialPointer != nil && (*rec.MaterialPointer).Scatter(&r, &rec, &attenuation, &scattered) {
			return attenuation.Mul(c.RayColor(scattered, depth-1, world))
		}
		return NewVec3(0, 0, 0)
	}

	// if not hit it renders the background color
	unitDirection := r.Direction.GetUnitVec()
	t := 0.5 * (unitDirection.Y + 1.0)
	// ((1-t) * <1,1,1>) + (t * <0.5,0.7,1>)
	return (NewVec3(1.0, 1.0, 1.0).Scale(1.0 - t)).Add((NewVec3(0.5, 0.7, 1.0).Scale(t)))
}

type HitRecord struct {
	T               float64
	P               Vec3
	Normal          Vec3
	FrontFace       bool
	MaterialPointer *Material
}

func (h *HitRecord) SetFaceNormal(r *Ray, outwardNormal Vec3) {
	h.FrontFace = Dot(&r.Direction, &outwardNormal) < 0
	if h.FrontFace {
		h.Normal = outwardNormal
	} else {
		h.Normal = outwardNormal.Negate()
	}
}

type Hittable interface {
	Hit(r Ray, i Interval, rec *HitRecord) bool
}
type HittableList struct {
	List []Hittable
}

func NewHittableList(h ...Hittable) *HittableList {
	return &HittableList{
		List: h,
	}
}
func (hl *HittableList) Add(h Hittable) {
	hl.List = append(hl.List, h)
}
func (hl *HittableList) Hit(r Ray, i Interval, rec *HitRecord) bool {
	var temp HitRecord
	hitAnything := false
	closestSoFar := i.Max
	for _, hittableObject := range hl.List {
		if hittableObject.Hit(r, NewInterval(i.Min, closestSoFar), &temp) {
			hitAnything = true
			closestSoFar = temp.T
			*rec = temp
		}
	}
	return hitAnything
}

type Sphere struct {
	Center Vec3
	Radius float64
	Mat    Material
}

func NewSphere(c Vec3, r float64, mat Material) Sphere {
	return Sphere{Center: c, Radius: max(0, r), Mat: mat}
}

func (s *Sphere) Hit(r Ray, i Interval, rec *HitRecord) bool {
	oc := r.Origin.Sub(s.Center)
	a, b, c := r.Direction.LengthSquared(), Dot(&oc, &r.Direction), oc.LengthSquared()-(s.Radius*s.Radius)
	discriminant := b*b - a*c
	rec.MaterialPointer = &s.Mat

	if discriminant > 0 {
		sqrtDiscriminant := math.Sqrt(discriminant)
		posRoot, negRoot := (-b-sqrtDiscriminant)/a, (-b+sqrtDiscriminant)/a
		hitT := -1.0
		if i.Surrounds(posRoot) {
			hitT = posRoot
		} else if i.Surrounds(negRoot) {
			hitT = negRoot
		}

		if hitT != -1.0 {
			rec.T = hitT
			rec.P = r.PointAtParameter(rec.T)
			outwardNormal := rec.P.Sub(s.Center).Scale(1.0 / s.Radius)
			rec.SetFaceNormal(&r, outwardNormal)
			return true
		}
	}
	return false
}

type Camera struct {
	ImageWidth        int
	ImageHeight       int
	SamplesPerPixel   int
	MaxDepth          int
	AspectRatio       float64
	PixelSamplesScale float64
	VFov              float64
	defocusAngle      float64
	focusDistance     float64
	Center            Vec3
	PixelDeltaU       Vec3
	PixelDeltaV       Vec3
	Pixel00Loc        Vec3
	lookFrom          Vec3
	lookAt            Vec3
	vup               Vec3
	u                 Vec3
	v                 Vec3
	w                 Vec3
	defocusDiskU      Vec3
	defocusDiskV      Vec3
}

func NewCamera() Camera {
	return Camera{
		ImageWidth:      100,
		SamplesPerPixel: 10,
		MaxDepth:        50,
		AspectRatio:     1.0,
		VFov:            90,
		defocusAngle:    0,
		focusDistance:   10,
		lookFrom:        NewVec3(0.0, 0.0, 0.0),
		lookAt:          NewVec3(0.0, 0.0, -1.0),
		vup:             NewVec3(0.0, 1.0, 0.0),
	}
}

func DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
func (c *Camera) InitCamera() {
	c.ImageHeight = max(int(float64(c.ImageWidth)/c.AspectRatio), 1)
	c.Center = c.lookFrom
	c.PixelSamplesScale = 1.0 / float64(c.SamplesPerPixel)

	theta := DegreesToRadians(c.VFov)
	h := math.Tan(theta / 2)

	viewPortHeight := 2.0 * h * c.focusDistance
	viewPortWidth := viewPortHeight * (float64(c.ImageWidth) / float64(c.ImageHeight))

	c.w = (c.lookFrom.Sub(c.lookAt)).GetUnitVec()
	c.u = (Cross(c.vup, c.w)).GetUnitVec()
	c.v = Cross(c.w, c.u)

	viewPortU := c.u.Scale(viewPortWidth)
	viewPortV := ((c.v).Negate()).Scale(viewPortHeight)

	c.PixelDeltaU, c.PixelDeltaV = viewPortU.Scale(1.0/float64(c.ImageWidth)), viewPortV.Scale(1.0/float64(c.ImageHeight))
	viewPortUpperLeft := c.Center.Sub(c.w.Scale(c.focusDistance)).Sub(viewPortU.Scale(0.5)).Sub(viewPortV.Scale(0.5)) // center - <0,0,focal length> - (viewportU / 2) - (viewportV / 2)
	c.Pixel00Loc = viewPortUpperLeft.Add((c.PixelDeltaU.Add(c.PixelDeltaV)).Scale(0.5))
	// viewPortUpperLeft + 0.5*(PixelDeltaU + PixelDeltaV)
	defocusRadius := c.focusDistance * math.Tan(DegreesToRadians(c.defocusAngle/2))
	c.defocusDiskU = c.u.Scale(defocusRadius)
	c.defocusDiskV = c.v.Scale(defocusRadius)

}
func (c *Camera) GetRay(i, j float64) Ray {
	offset := NewBoundedRandomVec(-0.5, 0.5)
	pixelSample := c.Pixel00Loc.Add(c.PixelDeltaU.Scale(i + offset.X).Add(c.PixelDeltaV.Scale(j + offset.Y)))

	var rayOrigin Vec3
	if c.defocusAngle <= 0 {
		rayOrigin = c.Center
	} else {
		rayOrigin = c.defocusDiskSample()
	}
	rayDirection := pixelSample.Sub(rayOrigin)
	return NewRay(rayOrigin, rayDirection)
}
func (c *Camera) Render(world Hittable) {

	c.InitCamera()
	initImage(c.ImageWidth, c.ImageHeight)
	lastPercent := -1
	for i := 0; i < c.ImageHeight; i++ {
		for j := 0; j < c.ImageWidth; j++ {
			pixelColor := NewVec3(0.0, 0.0, 0.0)
			for sample := 0; sample < c.SamplesPerPixel; sample++ {
				r := c.GetRay(float64(j), float64(i))
				pixelColor.PlusEq(c.RayColor(r, c.MaxDepth, world))
			}
			WriteColor(pixelColor.Scale(c.PixelSamplesScale))
		}
		percent := (i*100)/c.ImageHeight + 1
		if percent%5 == 0 && percent != lastPercent {
			fmt.Printf("%d percent done.\n", percent)
			lastPercent = percent
		}
	}
}

func (c *Camera) defocusDiskSample() Vec3 {
	p := RandomInUnitDisk()
	return c.Center.Add(c.defocusDiskU.Scale(p.X).Add(c.defocusDiskV.Scale(p.Y)))
}

func LinearToGamma(linearComponent float64) float64 {
	if linearComponent > 0 {
		return math.Sqrt(linearComponent)
	}
	return 0
}

type Material interface {
	Scatter(rIn *Ray, rec *HitRecord, attenuation *Vec3, scattered *Ray) bool
}

type Lambertian struct {
	Albedo Vec3
}

func NewLambertian(albedo Vec3) Lambertian {
	return Lambertian{Albedo: albedo}
}
func (l *Lambertian) Scatter(rIn *Ray, rec *HitRecord, attenuation *Vec3, scattered *Ray) bool {
	scatterDirection := rec.Normal.Add(RandomUnitVector())
	if scatterDirection.NearZero() {
		scatterDirection = rec.Normal
	}
	*scattered = NewRay(rec.P, scatterDirection)
	*attenuation = l.Albedo
	return true
}

type Metal struct {
	Albedo Vec3
	Fuzz   float64
}

func NewMetal(albedo Vec3, fuzz float64) Metal {
	return Metal{Albedo: albedo, Fuzz: min(1, fuzz)}
}
func (m *Metal) Scatter(rIn *Ray, rec *HitRecord, attenuation *Vec3, scattered *Ray) bool {
	reflected := Reflect(&rIn.Direction, &rec.Normal)
	reflected = reflected.GetUnitVec().Add(RandomUnitVector().Scale(m.Fuzz))
	*scattered = NewRay(rec.P, reflected)
	*attenuation = m.Albedo
	return Dot(&scattered.Direction, &rec.Normal) > 0
}

type Dielectric struct {
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
	*scattered = NewRay(rec.P, direction)

	return true
}

func initImage(w, h int) {
	f, err := os.Create("out.ppm")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fmt.Fprintf(f, "P3\n%d %d\n255\n", w, h)
}

func main() {

	// materialGround := NewLambertian(NewVec3(0.8, 0.8, 0.0))
	// materialCenter := NewLambertian(NewVec3(0.1, 0.2, 0.5))
	// materialLeft := NewDielectric(1.50)
	// materialBubble := NewDielectric(1.0 / 1.50)
	// materialRight := NewMetal(NewVec3(0.8, 0.6, 0.2), 1.0)

	// s1 := NewSphere(NewVec3(0.0, -100.5, -1.0), 100, &materialGround)
	// s2 := NewSphere(NewVec3(0.0, 0.0, -1.2), 0.5, &materialCenter)
	// s3 := NewSphere(NewVec3(-1.0, 0.0, -1.0), 0.5, &materialLeft)
	// s4 := NewSphere(NewVec3(-1.0, 0.0, -1.0), 0.4, &materialBubble)
	// s5 := NewSphere(NewVec3(1.0, 0.0, -1.0), 0.5, &materialRight)

	// world := NewHittableList(&s1, &s2, &s3, &s4, &s5)

	// R := math.Cos(math.Pi / 4)
	// materialLeft := NewLambertian(NewVec3(0.0, 0.0, 1.0))
	// materialRight := NewLambertian(NewVec3(1.0, 0.0, 0.0))

	// s1 := NewSphere(NewVec3(-R, 0.0, -1.0), R, &materialLeft)
	// s2 := NewSphere(NewVec3(R, 0.0, -1.0), R, &materialRight)

	// world := NewHittableList(&s1, &s2)

	groundMaterial := NewLambertian(NewVec3(0.5, 0.5, 0.5))
	s1 := NewSphere(NewVec3(0, -1000, -0), 1000, &groundMaterial)

	world := NewHittableList(&s1)

	for i := -11; i < 11; i++ {
		for j := -11; j < 11; j++ {
			chooseMat := rand.Float64()
			center := NewVec3(float64(i)+0.9*rand.Float64(), 0.2, float64(j)+0.9*rand.Float64())

			if center.Sub(NewVec3(4, 0.2, 0.0)).Length() > 0.9 {
				if chooseMat < 0.8 {
					albedo := NewVec3(rand.Float64(), rand.Float64(), rand.Float64()).Mul(NewVec3(rand.Float64(), rand.Float64(), rand.Float64()))
					sM := NewLambertian(albedo)
					sN := NewSphere(center, 0.2, &sM)
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

	cam := NewCamera()
	cam.AspectRatio = 16.0 / 9.0
	cam.ImageWidth = 1200
	cam.MaxDepth = 50
	cam.SamplesPerPixel = 500

	cam.VFov = 20
	cam.lookFrom = NewVec3(13.0, 2.0, 3.0)
	cam.lookAt = NewVec3(0.0, 0.0, 0.0)
	cam.vup = NewVec3(0.0, 1.0, 0.0)

	cam.defocusAngle = 0.6
	cam.focusDistance = 10.0

	cam.Render(world)

}
