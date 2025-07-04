package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

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
