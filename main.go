package main

import (
	"fmt"
	. "raytracer/camera"
)

func main() {
	img := NewImage(400, 16.0/9.0)
	cam := NewCamera(img, 100, 50)
	fmt.Print(cam.Render())
}
