package primitives

import (
	"fmt"
	"math"
	"raytracer/util"
	"strconv"
)

type Triple struct {
	X float64
	Y float64
	Z float64
}

func Triple_(x, y, z float64) Triple {
	return Triple{
		X: x,
		Y: y,
		Z: z,
	}
}

func (this Triple) String() string {
	return fmt.Sprintf("(%v, %v, %v)", this.X, this.Y, this.Z)
}

func Ray_(o Point, dir Vector) Ray {
	return Ray{
		Origin:    o,
		Direction: dir,
	}
}

type Ray struct {
	Origin    Point
	Direction Vector
}

func (this Ray) At(t float64) Point {
	return this.Origin.Translate(this.Direction.Scalar(t))
}

func Vector_(a, b, c float64) Vector {
	return Vector{
		X: a,
		Y: b,
		Z: c,
	}
}

func FromPoints(a, b Point) Vector {
	return Vector{
		X: b.X - a.X,
		Y: b.Y - a.Y,
		Z: b.Z - a.Z,
	}
}

func LinearCombination(
	a1 float64, v1 Vector,
	a2 float64, v2 Vector,
	a3 float64, v3 Vector) Vector {
	// a1*v1 + a2*v2 + a3*v3
	return v1.Scalar(a1).Add(v2.Scalar(a2)).Add(v3.Scalar(a3))
}

func RandomVector(min, max float64) Vector {
	return Vector_(util.Random(min, max), util.Random(min, max), util.Random(min, max))
}

func VerctorInUnitSphere() Vector {
	for {
		vec := RandomVector(-1, 1)
		// we don't need to take sqrt, since
		// if 0 <= x <= 1 then  0 <= sqrt(x) <= 1
		// and sqrt is expensive
		if vec.LengthSquared() < 1 {
			return vec
		}
	}
}

// Lambertian sphere render
func UnitSphereVector() Vector {
	return VerctorInUnitSphere().UnitVector()
}

type Vector Triple

// produto de escalar por vetor
func (this Vector) Scalar(a float64) Vector {
	return Vector{
		X: this.X * a,
		Y: this.Y * a,
		Z: this.Z * a,
	}
}

func (this Vector) Add(other Vector) Vector {
	return Vector{
		X: this.X + other.X,
		Y: this.Y + other.Y,
		Z: this.Z + other.Z,
	}
}

func (this Vector) Sub(other Vector) Vector {
	return Vector{
		X: this.X - other.X,
		Y: this.Y - other.Y,
		Z: this.Z - other.Z,
	}
}

// produto escalar
func (this Vector) Dot(other Vector) float64 {
	return this.X*other.X + this.Y*other.Y + this.Z*other.Z
}

// produto vetorial
func (this Vector) Cross(other Vector) Vector {
	return Vector{
		X: this.Y*other.Z - this.Z*other.Y,
		Y: this.Z*other.X - this.X*other.Z,
		Z: this.X*other.Y - this.Y*other.X,
	}
}

func (this Vector) UnitVector() Vector {
	return this.Scalar(1.0 / this.Length())
}

func (this Vector) Length() float64 {
	return math.Sqrt(this.LengthSquared())
}

// this.Dot(this)
func (this Vector) LengthSquared() float64 {
	return this.X*this.X + this.Y*this.Y + this.Z*this.Z
}

type Point Triple

func Point_(x, y, z float64) Point {
	return Point{
		X: x,
		Y: y,
		Z: z,
	}
}

func (this Point) Translate(v Vector) Point {
	return Point{
		X: this.X + v.X,
		Y: this.Y + v.Y,
		Z: this.Z + v.Z,
	}
}

type Color struct {
	R float64 // [0, 1]
	G float64
	B float64
}

func Color_(r, g, b float64) Color {
	return Color{
		R: r,
		G: g,
		B: b,
	}
}

func (this Color) Marshal(buff *[]byte, index int) int {
	r := byte(this.R * 255)
	g := byte(this.G * 255)
	b := byte(this.B * 255)
	index = marshalByte(r, buff, index)
	(*buff)[index] = ' '
	index++
	index = marshalByte(g, buff, index)
	(*buff)[index] = ' '
	index++
	index = marshalByte(b, buff, index)
	return index
}

func (this Color) Mult(a float64) Color {
	return Color{
		R: this.R * a,
		G: this.G * a,
		B: this.B * a,
	}
}

func (this Color) Div(a float64) Color {
	return Color{
		R: this.R / a,
		G: this.G / a,
		B: this.B / a,
	}
}

func (this Color) Add(other Color) Color {
	return Color{
		R: this.R + other.R,
		G: this.G + other.G,
		B: this.B + other.B,
	}
}

func (this Color) Gamma(factor float64) Color {
	g := 1.0 / factor
	return Color{
		R: math.Pow(this.R, g),
		G: math.Pow(this.G, g),
		B: math.Pow(this.B, g),
	}
}

func marshalByte(b byte, buff *[]byte, index int) int {
	s := strconv.Itoa(int(b))
	for i := range s {
		(*buff)[index] = s[i]
		index++
	}
	return index
}

func Hit_(s float64, p Point, n Vector, r Ray) Hit {
	return Hit{
		Scalar:    s,
		Point:     p,
		Normal:    n,
		FrontFace: n.Dot(r.Direction) > 0, // is this ok?
	}
}

type Hit struct {
	Scalar float64
	Point  Point
	Normal Vector

	FrontFace bool
}

var NoHit = Hit{}

type Hittable interface {
	Hit(r Ray, t_min, t_max float64) (Hit, bool)
}
