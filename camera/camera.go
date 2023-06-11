package camera

import (
	"math"
	el "raytracer/elements"
	. "raytracer/primitives"
	"raytracer/util"
	"strconv"
	"time"
)

var PlusInf = math.Inf(+1)
var MinusInf = math.Inf(-1)

type World []Hittable

func (this *World) Add(a Hittable) {
	(*this) = append(*this, a)
}

func (this *World) Clear() {
	for i := range *this {
		(*this)[i] = nil
	}
	*this = (*this)[:0]
}

func (this *World) Hit(r Ray, t_min, t_max float64) (Hit, bool) {
	closestSoFar := t_max
	hitAnything := false
	out := NoHit

	for _, hittable := range *this {
		hit, ok := hittable.Hit(r, t_min, closestSoFar)
		if ok {
			hitAnything = true
			closestSoFar = hit.Scalar
			out = hit
		}
	}
	return out, hitAnything
}

var objs = World{
	el.Sphere{
		Center: Point_(0, 0, -1),
		Radius: 0.5,
	},
	el.Sphere{
		Center: Point_(-1, 0, -2),
		Radius: 0.5,
	},
	el.Sphere{
		Center: Point_(1, 0, -2),
		Radius: 0.5,
	},
	el.Sphere{
		Center: Point_(-3, 0, -4),
		Radius: 1,
	},
	el.Sphere{
		Center: Point_(3, 0, -4),
		Radius: 1,
	},
	el.Sphere{
		Center: Point_(0, -100.5, -1),
		Radius: 100,
	},
}

func rayColor(r Ray, depth int) Color {
	if depth <= 0 {
		return Color_(0, 0, 0)
	}
	hit, ok := objs.Hit(r, 0.001, PlusInf)
	if ok {
		randVec := VerctorInUnitSphere()
		target := hit.Point.Translate(hit.Normal.Add(randVec))
		r := Ray_(hit.Point, FromPoints(hit.Point, target))
		return rayColor(r, depth-1).Mult(0.5)
	}
	unit := r.Direction.UnitVector()
	t := (unit.Y + 1.0) * 0.5
	a := Color_(1.0, 1.0, 1.0).Mult(1.0 - t)
	b := Color_(0.5, 0.7, 1.0).Mult(t)
	return a.Add(b)
}

func normalColor(v Vector) Color {
	return Color_(v.X+1, v.Y+1, v.Z+1).Mult(0.5)
}

type Camera struct {
	ViewportHeight float64
	ViewportWidth  float64
	FocalLength    float64

	Origin Point

	LowerLeftCorner Vector
	Horizontal      Vector
	Vertical        Vector

	SamplesPerPixel int
	MaxRayDepth     int

	Image *Image
}

func NewCamera(img *Image, samples int, maxRayDepth int) *Camera {
	if samples <= 0 {
		samples = 1
	}
	vpHeight := 2.0
	vpWidth := vpHeight * img.AspectRatio
	focalLength := 1.0

	horVec := Vector_(vpWidth, 0, 0)
	verVec := Vector_(0, vpHeight, 0)
	var llcVec Vector
	{ // creates the lower left corner vector by linear combination
		// llvVec = (-0.5)*h + (-0.5)*v + (-1)*d
		depVec := Vector_(0, 0, focalLength)
		llcVec = LinearCombination(-0.5, horVec, -0.5, verVec, -1, depVec)
	}
	return &Camera{
		ViewportHeight:  vpHeight,
		ViewportWidth:   vpWidth,
		FocalLength:     focalLength,
		Origin:          Point_(0, 0, 0),
		LowerLeftCorner: llcVec,
		Horizontal:      horVec,
		Vertical:        verVec,
		Image:           img,
		SamplesPerPixel: samples,
		MaxRayDepth:     maxRayDepth,
	}
}

func (this *Camera) Render() string {
	start := time.Now()
	util.Println() // so the bar doesn't eat a line
	for y := range this.Image.Data {
		for x := range this.Image.Data[y] {
			color := Color_(0, 0, 0)
			for i := 0; i < this.SamplesPerPixel; i++ { // anti-aliasing
				u := (float64(x) + util.Random(-1, 1)) / float64(this.Image.Width-1)
				v := (float64(y) + util.Random(-1, 1)) / float64(this.Image.Height-1)
				dir := LinearCombination(
					1, this.LowerLeftCorner,
					u, this.Horizontal,
					v, this.Vertical)
				r := Ray{
					Origin:    this.Origin,
					Direction: dir,
				}
				color = color.Add(rayColor(r, this.MaxRayDepth))
			}
			this.Image.Data[y][x] = color.Div(float64(this.SamplesPerPixel)).Gamma(2)
		}
		util.Bar(y+1, this.Image.Height)
	}
	out := this.Image.ToPPM()
	util.Printf("took: %v\n", time.Since(start))
	return out
}

type Image struct {
	Width  int
	Height int

	AspectRatio float64

	// format is: [y][x]Pixel
	Data [][]Color
}

func NewImage(width int, aspectRatio float64) *Image {
	height := int(float64(width) / aspectRatio)
	a := make([][]Color, height)
	for i := range a {
		a[i] = make([]Color, width)
	}
	return &Image{
		Width:       width,
		Height:      height,
		AspectRatio: aspectRatio,
		Data:        a,
	}
}

func (this *Image) MapPixels(m func(x, y int, img *Image) Color) {
	for y := range this.Data {
		for x := range this.Data[y] {
			this.Data[y][x] = m(x, y, this)
		}
	}
}

// We render it upside down :)
func (this *Image) ToPPM() string {
	s := "P3\n" +
		strconv.Itoa(this.Width) + " " +
		strconv.Itoa(this.Height) + "\n255\n"
	// max per line is "255 255 255\n" which is 12 characters
	buffer := make([]byte, len(s)+(this.Width*this.Height*12))
	copy(buffer, []byte(s))
	i := len(s)
	for y := len(this.Data) - 1; y >= 0; y-- {
		for _, px := range this.Data[y] {
			i = px.Marshal(&buffer, i)
			buffer[i] = '\n'
			i++
		}
	}
	return string(buffer)
}
