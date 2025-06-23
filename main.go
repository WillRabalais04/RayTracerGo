package main

import (
	"fmt"
	"os"
	"strings"
)

type vec3 struct {
	
}

func main(){
	
	nx, ny := 200, 100
	var out strings.Builder
	out.WriteString(fmt.Sprintf("P3\n%d %d\n255\n", nx, ny))

	for j := ny-1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			r,g,b := (float32(i) / float32(nx)), (float32(j) / float32(nx)), 0.2
			ir,ig,ib := int(255.99*r), int(255.99*g), int(255.99*b)
			out.WriteString(fmt.Sprintf("%d %d %d\n", ir, ig, ib))
		}
	}
	err := os.WriteFile("out.ppm", []byte(out.String()), 0644)
	if err != nil {
		fmt.Println("Error writing file: ", err)
	}

}