package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
)

type Texture interface {
	Value(u, v float64, p *Vec3) Vec3
}

type SolidColor struct {
	Albedo Vec3
}

func NewSolidColor(albedo *Vec3) *SolidColor {
	return &SolidColor{Albedo: *albedo}
}
func NewSolidColorFromRGB(r, g, b float64) *SolidColor {
	return NewSolidColor(NewVec3Pointer(r, g, b))
}
func (c *SolidColor) Value(u, v float64, p *Vec3) Vec3 {
	return c.Albedo
}

type CheckeredTexture struct {
	InvScale float64
	Even     Texture
	Odd      Texture
}

func NewCheckeredTexture(scale float64, even, odd Texture) CheckeredTexture {
	return CheckeredTexture{InvScale: 1.0 / scale, Even: even, Odd: odd}
}

func NewCheckeredTextureFromColors(scale float64, c1, c2 Vec3) CheckeredTexture {
	t1, t2 := NewSolidColor(&c1), NewSolidColor(&c2)
	return NewCheckeredTexture(scale, t1, t2)
}
func (t *CheckeredTexture) Value(u, v float64, p *Vec3) Vec3 {
	isEven := int((math.Floor(t.InvScale*p.X)+math.Floor(t.InvScale*p.Y)+math.Floor(t.InvScale*p.Z)))%2 == 0
	// even & odd map to the two different colors in the checkered pattern
	if isEven {
		return t.Even.Value(u, v, p)
	} else {
		return t.Odd.Value(u, v, p)
	}
}

type NoiseTexture struct {
	Noise PerlinNoise
	Scale float64
}

func NewNoiseTexture(scale float64) NoiseTexture {
	return NoiseTexture{Noise: NewPerlinNoise(), Scale: scale}
}

func (t *NoiseTexture) Value(u, v float64, p *Vec3) Vec3 {
	// return NewVec3(1.0, 1.0, 1.0).Scale(t.Noise.Turbulence(p, 7))
	return NewVec3(0.5, 0.5, 0.5).Scale(1 + math.Sin(t.Scale*p.Z+10*t.Noise.Turbulence(p, 7)))
}

type ImageTexture struct {
	img LoadedImage
}

func (t *ImageTexture) Value(u, v float64, p *Vec3) Vec3 {
	if t.img.Height <= 0 {
		return NewVec3(0, 1, 1)
	}
	interval := NewInterval(0, 1)
	u, v = interval.Clamp(u), 1-interval.Clamp(v)
	i, j := u*float64(t.img.Width), v*float64(t.img.Height)
	p1, p2, p3 := t.img.PixelData(i, j)
	// fmt.Printf("RGB = %.2f %.2f %.2f\n", p1, p2, p3)
	// fmt.Println(SrgbToLinear(255))
	// colorScale := 1.0 / 255.0
	// return NewVec3(colorScale*p1, colorScale*p2, colorScale*p3)
	return NewVec3(p1, p2, p3)

}

type LoadedImage struct {
	Pixels        []float64
	Width, Height int
}

func NewImageTexture(filename string) *ImageTexture {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	img, imgFmt, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decoding image (%s)\n", imgFmt)

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	li := LoadedImage{Pixels: make([]float64, w*h*3),
		Width:  w,
		Height: h,
	}
	tex := ImageTexture{
		img: li,
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := color.NRGBAModel.Convert(img.At(x+bounds.Min.X, y+bounds.Min.Y)).(color.NRGBA)
			i := (y*w + x) * 3
			tex.img.Pixels[i+0] = SrgbToLinear(c.R)
			tex.img.Pixels[i+1] = SrgbToLinear(c.G)
			tex.img.Pixels[i+2] = SrgbToLinear(c.B)
		}
	}
	return &tex

}

func (t *LoadedImage) PixelData(x, y float64) (r, g, b float64) {
	i1 := NewInterval(0, float64(t.Width-1))
	i2 := NewInterval(0, float64(t.Height-1))
	idx := (int(i2.Clamp(y))*t.Width + int(i1.Clamp(x))) * 3

	return t.Pixels[idx], t.Pixels[idx+1], t.Pixels[idx+2]
}
