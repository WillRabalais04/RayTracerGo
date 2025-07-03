package main

import (
	"fmt"
	"log"
	"math"
	"os"
)

func LinearToGamma(linearComponent float64) float64 {
	if linearComponent > 0 {
		return math.Sqrt(linearComponent)
	}
	return 0
}
func BoxCompare(axisIndex int, a, b Hittable) bool {
	if a.BBOX() == nil || b.BBOX() == nil {
		fmt.Printf("boxcomparefailed")
		return false
	}
	// compares min of chosen axis of a and b
	return a.BBOX().axisInterval(axisIndex).Min < b.BBOX().axisInterval(axisIndex).Min
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
func InitImage(w, h int) {
	f, err := os.Create("out.ppm")
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

	file, err := os.OpenFile("out.ppm", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Fprintf(file, "%d %d %d\n", r, g, b)
}
