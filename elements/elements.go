package elements

import (
	"math"
	. "raytracer/primitives"
)

type Sphere struct {
	Center Point
	Radius float64
}

func (this Sphere) Hit(r Ray, t_min, t_max float64) (Hit, bool) {
	CA := FromPoints(this.Center, r.Origin)
	a := r.Direction.Dot(r.Direction)
	half_b := r.Direction.Dot(CA)
	c := CA.Dot(CA) - this.Radius*this.Radius

	discriminant := half_b*half_b - a*c
	if discriminant < 0 {
		return NoHit, false
	}
	sqD := math.Sqrt(discriminant)
	root := (-half_b - sqD) / a
	if root < t_min || root > t_max {
		root = (-half_b + sqD) / a
		if root < t_min || root > t_max {
			return NoHit, false
		}
	}
	p := r.At(root)
	normal := FromPoints(this.Center, p).UnitVector()
	hit := Hit_(root, p, normal, r)
	return hit, true
}
