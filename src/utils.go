package main

import (
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"os"
)

func InitImage(w, h int) {
	f, err := os.Create("../out.ppm")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fmt.Fprintf(f, "P3\n%d %d\n255\n", w, h)
}
func WriteColor(color Vec3) {

	intensity := NewInterval(0.0, 0.999)

	r := int(256 * intensity.Clamp(LinearToGamma(color.X)))
	g := int(256 * intensity.Clamp(LinearToGamma(color.Y)))
	b := int(256 * intensity.Clamp(LinearToGamma(color.Z)))

	file, err := os.OpenFile("../out.ppm", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Fprintf(file, "%d %d %d\n", r, g, b)
}
func LinearToGamma(linearComponent float64) float64 {
	if linearComponent > 0 {
		return math.Sqrt(linearComponent)
	}
	return 0
}
func DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
func ArgMax(sizes ...float64) int {
	maxIndex := 0
	maxValue := sizes[0]
	for i, v := range sizes[1:] {
		if v > maxValue {
			maxIndex = i + 1
			maxValue = v
		}
	}
	return maxIndex
}

func SrgbToLinear(c byte) float64 {
	v := float64(c) / 255
	if v <= 0.04045 {
		return v / 12.92
	}
	return (math.Pow((float64((v + 0.055) / (1.055))), 2.4))
}

func RandFloatInRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
