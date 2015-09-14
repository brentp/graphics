// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rings

import (
	"math"
	"math/rand"

	"github.com/gonum/plot/vg"

	"github.com/biogo/graphics/bezier"
)

// LengthDist generates a random value in the range [Length*Min, Length*Max), depending on a
// provided random factor.
type LengthDist struct {
	Length   vg.Length
	Min, Max *float64 // A nil value is interpreted as 1.
}

// Perturb returns a perturbed vg.Length value. Calling Perturb on a nil LengthDist will panic.
func (p *LengthDist) Perturb(f float64) vg.Length {
	if p.Min == nil && p.Max == nil {
		return p.Length
	}
	var min, max = 1., 1.
	if p.Min != nil {
		min = *p.Min
	}
	if p.Max != nil {
		max = *p.Max
	}
	return p.Length * vg.Length(min+(max-min)*f)
}

// FactorDist generates a random value in the range [Length*Min, Length*Max), depending on a
// provided random factor.
type FactorDist struct {
	Factor   float64
	Min, Max *float64 // A nil value is interpreted as 1.
}

// Perturb returns a perturbed float value. Calling Perturb on a nil FactorDist will panic.
func (p *FactorDist) Perturb(f float64) float64 {
	if p.Min == nil && p.Max == nil {
		return p.Factor
	}
	var min, max = 1., 1.
	if p.Min != nil {
		min = *p.Min
	}
	if p.Max != nil {
		max = *p.Max
	}
	return p.Factor * (min + (max-min)*f)
}

// Bezier defines Bézier control points for a link between features represented by Links and Ribbons.
type Bezier struct {
	// Segments defines the number of segments to draw when rendering the curve.
	Segments int

	// Radius, Crest and Purity define aspects of Bézier geometry.
	//
	// See http://circos.ca/documentation/tutorials/links/geometry/images for a detailed explanation
	// of radius, crest and purity.
	//
	// Radius specifies the Bézier radius of a curve generated by the Bezier.
	Radius LengthDist
	// Crest and Purity specify the crest and purity behaviour of a curve generated by the Bezier.
	// If nil, these values are not used.
	Crest  *FactorDist
	Purity *FactorDist
}

// ControlPoints returns a set of Bézier curve control points defining the path between the points defined
// by the parameters and the Bezier's Radius, Crest and Purity fields.
func (b *Bezier) ControlPoints(a [2]Angle, rad [2]vg.Length) []bezier.Point {
	var p [2]Point
	for i := range a {
		p[i] = Rectangular(a[i], float64(rad[i]))
	}

	var radius = b.Radius
	if b.Purity != nil {
		bisectRadius := vg.Length(math.Hypot((p[0].X+p[1].X)/2, (p[0].Y+p[1].Y)/2))
		radius.Length += vg.Length(b.Purity.Perturb(rand.Float64())-1) * (radius.Length - bisectRadius)
	}

	var bisect Angle
	if math.Abs(float64(a[1]-a[0])) > math.Pi {
		bisect = (a[0]+a[1]+Angle(2*math.Pi))/2 - Angle(2*math.Pi)
	} else {
		bisect = (a[1] + a[0]) / 2
	}
	mp := Rectangular(bisect, float64(radius.Perturb(rand.Float64())))
	mid := bezier.Point{X: mp.X, Y: mp.Y}

	if b.Crest != nil {
		points := []bezier.Point{
			0: {X: p[0].X, Y: p[0].Y},
			2: mid,
			4: {X: p[1].X, Y: p[1].Y},
		}
		c := b.Crest.Perturb(rand.Float64())

		var cp Point
		for i, r := range rad {
			cp = Rectangular(a[i], float64(r)-float64(r-radius.Length)*c)
			points[2*i+1] = bezier.Point{X: cp.X, Y: cp.Y}
		}
		return points
	}

	return []bezier.Point{
		{X: p[0].X, Y: p[0].Y},
		mid,
		{X: p[1].X, Y: p[1].Y},
	}
}
