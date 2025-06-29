package main

import (
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"os"
)

type Vec3 struct {
	X, Y, Z float32
}

// Vec3 Ops
func NewVec3(x, y, z float32) Vec3 {
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
func (v Vec3) Scale(t float32) Vec3 {
	return Vec3{X: t * v.X, Y: t * v.Y, Z: t * v.Z}
}
func (v *Vec3) ScaleAssign(t float32) {
	v.X *= t
	v.Y *= t
	v.Z *= t
}
func (v1 Vec3) Dot(v2 Vec3) float32 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}
func (v1 Vec3) Cross(v2 Vec3) Vec3 { // check this
	return Vec3{
		X: (v1.Y*v2.Z - v1.Z*v2.Y),
		Y: (v1.Z*v2.X - v1.X*v2.Z),
		Z: (v1.X*v2.Y - v1.Y*v2.X)}
}
func (v *Vec3) MakeUnit() {
	v.ScaleAssign((1.0 / v.Length()))
}
func (v Vec3) GetUnit() Vec3 {
	return v.Scale((1.0 / v.Length()))
}
func (v Vec3) Negate() Vec3 {
	return Vec3{X: -v.X, Y: -v.Y, Z: -v.Z}
}
func (v Vec3) SquaredLength() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}
func (v Vec3) Length() float32 {
	return float32(math.Sqrt(float64(v.SquaredLength())))
}
func RandomUnitVector() Vec3 {
	for {
		p := NewBoundedRandomVec(-1, 1)
		lensq := p.SquaredLength()
		if 1e-160 < lensq && lensq <= 1 {
			return p.Scale(1.0 / float32(math.Sqrt(float64(lensq))))
		}
	}
}
func (Normal Vec3) RandomOnHemisphere() Vec3 {
	OnUnitSphere := RandomUnitVector()
	if OnUnitSphere.Dot(Normal) < 0.0 {
		OnUnitSphere.ScaleAssign(-1.0)
	}
	return OnUnitSphere
}
func NewRandomVec() Vec3 {
	return NewVec3(rand.Float32(), rand.Float32(), rand.Float32())
}
func NewBoundedRandomVec(min, max float32) Vec3 {
	r := max - min
	return NewVec3(rand.Float32()*r+min, rand.Float32()*r+min, rand.Float32()*r+min)
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
	r, g, b := int(256*(intensity.Clamp(LinearToGamma(float64(color.X))))), int(256*(intensity.Clamp(LinearToGamma(float64(color.Y))))), int(256*(intensity.Clamp(LinearToGamma(float64(color.Z)))))

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
func (r Ray) PointAtParameter(t float32) Vec3 {
	return r.Origin.Add(r.Direction.Scale(t))
}
func NewRay(o, d Vec3) Ray {
	return Ray{Origin: o, Direction: d}
}
func (r Ray) Color(world Hittable, depth int) Vec3 {
	if depth <= 0 {
		return NewVec3(0.0, 0.0, 0.0)
	}

	var rec HitRecord
	if world.Hit(r, NewInterval(0.001, math.MaxFloat64), &rec) {
		direction := rec.Normal.Add(RandomUnitVector())
		return (NewRay(rec.P, direction).Color(world, depth-1)).Scale(0.5)
	}
	// if not hit it renders the background color
	unitDirection := r.Direction.GetUnit()
	t := 0.5 * (unitDirection.Y + 1.0)
	// ((1-t) * <1,1,1>) + (t * <0.5,0.7,1>)
	return (NewVec3(1.0, 1.0, 1.0).Scale(1.0 - t)).Add((NewVec3(0.5, 0.7, 1.0).Scale(t)))
}

type HitRecord struct {
	T               float32
	P               Vec3
	Normal          Vec3
	FrontFace       bool
	MaterialPointer *Material
}

func (h *HitRecord) SetFaceNormal(r Ray, outwardNormal Vec3) {
	FrontFace := r.Direction.Dot(outwardNormal) < 0
	if FrontFace {
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
			closestSoFar = float64(temp.T)
			*rec = temp
		}
	}
	return hitAnything
}

type Sphere struct {
	Center Vec3
	Radius float32
}

func NewSphere(c Vec3, r float32) Sphere {
	return Sphere{Center: c, Radius: r}
}

func (s *Sphere) Hit(r Ray, i Interval, rec *HitRecord) bool {
	oc := r.Origin.Sub(s.Center)
	a, b, c := r.Direction.SquaredLength(), oc.Dot(r.Direction), oc.Dot(oc)-(s.Radius*s.Radius)
	discriminant := b*b - a*c

	if discriminant > 0 {
		sqrtDiscriminant := float32(math.Sqrt(float64(discriminant)))
		posRoot, negRoot := (-b-sqrtDiscriminant)/a, (-b+sqrtDiscriminant)/a
		hitT := float32(-1.0)
		if i.Surrounds(float64(posRoot)) {
			hitT = posRoot
		} else if i.Surrounds(float64(negRoot)) {
			hitT = negRoot
		}

		if hitT != -1.0 {
			rec.T = hitT
			rec.P = r.PointAtParameter(rec.T)
			outwardNormal := rec.P.Sub(s.Center).Scale(1.0 / s.Radius)
			rec.SetFaceNormal(r, outwardNormal)
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
	Center            Vec3
	PixelDeltaU       Vec3
	PixelDeltaV       Vec3
	Pixel00Loc        Vec3
}

func NewCamera() Camera {
	return Camera{
		ImageWidth:      100,
		SamplesPerPixel: 10,
		MaxDepth:        50,
		AspectRatio:     1.0,
	}
}
func (c *Camera) InitCamera() {
	c.ImageHeight = max(int(float64(c.ImageWidth)/c.AspectRatio), 1)
	c.Center = NewVec3(0.0, 0.0, 0.0)
	c.PixelSamplesScale = 1.0 / float64(c.SamplesPerPixel)
	focalLength := 1.0
	viewPortHeight := 2.0
	viewPortWidth := viewPortHeight * (float64(c.ImageWidth) / float64(c.ImageHeight))
	viewPortU, viewPortV := NewVec3(float32(viewPortWidth), 0.0, 0.0), NewVec3(0.0, -float32(viewPortHeight), 0.0)
	c.PixelDeltaU, c.PixelDeltaV = viewPortU.Scale(1.0/float32(c.ImageWidth)), viewPortV.Scale(1.0/float32(c.ImageHeight))
	viewPortUpperLeft := c.Center.Sub(NewVec3(0.0, 0.0, float32(focalLength))).Sub(viewPortU.Scale(0.5)).Sub(viewPortV.Scale(0.5)) // center - <0,0,focal length> - (viewportU / 2) - (viewportV / 2)
	c.Pixel00Loc = viewPortUpperLeft.Add((c.PixelDeltaU.Add(c.PixelDeltaV)).Scale(0.5))                                            // viewPortUpperLeft + 0.5*(PixelDeltaU + PixelDeltaV)
}
func (c *Camera) GetRay(i, j float32) Ray {
	offset := NewBoundedRandomVec(-0.5, 0.5)
	pixelSample := c.Pixel00Loc.Add(c.PixelDeltaU.Scale(i + offset.X).Add(c.PixelDeltaV.Scale(j + offset.Y)))
	rayOrigin := c.Center
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
				r := c.GetRay(float32(j), float32(i))
				pixelColor.PlusEq(r.Color(world, c.MaxDepth))
			}
			WriteColor(pixelColor.Scale(float32(c.PixelSamplesScale)))
		}
		percent := (i*100)/c.ImageHeight + 1
		if percent%5 == 0 && percent != lastPercent {
			fmt.Printf("%d percent done.\n", percent)
			lastPercent = percent
		}
	}
}

func LinearToGamma(linearComponent float64) float64 {
	if linearComponent > 0 {
		return math.Sqrt(linearComponent)
	}
	return 0
}

type Material interface {
	Scatter(r *Ray, rec *HitRecord, attenuation *Vec3, scattered *Ray) bool
}

type Lambertian struct {
	Albedo Vec3
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

	s1, s2 := NewSphere(NewVec3(0.0, 0.0, -1.0), 0.5), NewSphere(NewVec3(0.0, -100.5, -1.0), 100)
	world := NewHittableList(&s1, &s2)
	cam := NewCamera()
	cam.AspectRatio = 16.0 / 9.0
	cam.ImageWidth = 800
	cam.MaxDepth = 50
	cam.SamplesPerPixel = 100

	cam.Render(world)

}
