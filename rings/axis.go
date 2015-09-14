// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rings

import (
	"fmt"

	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"

	"github.com/biogo/biogo/feat"
)

// Axis represents the radial axis of ring, usually a Scores.
type Axis struct {
	// Angle specifies the angular location of the axis.
	Angle Angle

	// Label describes the axis label configuration.
	Label AxisLabel

	// LineStyle is the style of the axis line.
	LineStyle draw.LineStyle

	// Tick describes the scale's tick configuration.
	Tick TickConfig

	// Grid is the style of the grid lines.
	Grid draw.LineStyle
}

// AxisLabel describes an axis label format and text.
type AxisLabel struct {
	// Text is the axis label string.
	Text string

	// TextStyle is the style of the axis label text.
	draw.TextStyle

	// Placement determines the text rotation and alignment.
	// If Placement is nil, DefaultPlacement is used.
	Placement TextPlacement
}

// TickConfig describes an axis tick configuration.
type TickConfig struct {
	// Label is the TextStyle on the tick labels.
	Label draw.TextStyle

	// LineStyle is the LineStyle of the tick lines.
	LineStyle draw.LineStyle

	// Placement determines the text rotation and alignment.
	// If Placement is nil, DefaultPlacement is used.
	Placement TextPlacement

	// Length is the length of a major tick mark.
	// Minor tick marks are half of the length of major
	// tick marks.
	Length vg.Length

	// Marker returns the tick marks. Any tick marks
	// returned by the Marker function that are not in
	// range of the axis are not drawn.
	Marker plot.Ticker
}

// drawAt renders the axis at cen in the specified drawing area, according to the
// Axis configuration.
func (r *Axis) drawAt(ca draw.Canvas, cen draw.Point, fs []Scorer, base ArcOfer, inner, outer vg.Length, min, max float64) {
	locMap := make(map[feat.Feature]struct{})

	var (
		pa vg.Path
		e  Point

		marks []plot.Tick

		scale = (outer - inner) / vg.Length(max-min)
	)
	for _, f := range fs {
		locMap[f.Location()] = struct{}{}
	}
	if r.Grid.Color != nil && r.Grid.Width != 0 {
		for loc := range locMap {
			arc, err := base.ArcOf(loc, nil)
			if err != nil {
				panic(fmt.Sprint("rings: no arc for feature location:", err))
			}

			ca.SetLineStyle(r.Grid)
			marks = r.Tick.Marker.Ticks(min, max)
			for _, mark := range marks {
				if mark.Value < min || mark.Value > max {
					continue
				}
				pa = pa[:0]

				radius := vg.Length(mark.Value-min)*scale + inner

				e = Rectangular(arc.Theta, float64(radius))
				pa.Move(cen.X+vg.Length(e.X), cen.Y+vg.Length(e.Y))
				pa.Arc(cen.X, cen.Y, radius, float64(arc.Theta), float64(arc.Phi))

				ca.Stroke(pa)
			}
		}
	}

	if r.LineStyle.Color != nil && r.LineStyle.Width != 0 {
		pa = pa[:0]

		e = Rectangular(r.Angle, float64(inner))
		pa.Move(cen.X+vg.Length(e.X), cen.Y+vg.Length(e.Y))
		e = Rectangular(r.Angle, float64(outer))
		pa.Line(cen.X+vg.Length(e.X), cen.Y+vg.Length(e.Y))

		ca.SetLineStyle(r.LineStyle)
		ca.Stroke(pa)
	}

	if r.Tick.LineStyle.Color != nil && r.Tick.LineStyle.Width != 0 && r.Tick.Length != 0 {
		ca.SetLineStyle(r.Tick.LineStyle)
		if marks == nil {
			marks = r.Tick.Marker.Ticks(min, max)
		}
		for _, mark := range marks {
			if mark.Value < min || mark.Value > max {
				continue
			}
			pa = pa[:0]

			radius := vg.Length(mark.Value-min)*scale + inner

			var length vg.Length
			if mark.IsMinor() {
				length = r.Tick.Length / 2
			} else {
				length = r.Tick.Length
			}
			off := Rectangular(r.Angle+Complete/4, float64(length))
			e = Rectangular(r.Angle, float64(radius))
			pa.Move(cen.X+vg.Length(e.X), cen.Y+vg.Length(e.Y))
			pa.Line(cen.X+vg.Length(e.X+off.X), cen.Y+vg.Length(e.Y+off.Y))

			ca.Stroke(pa)

			if mark.IsMinor() || r.Tick.Label.Color == nil {
				continue
			}

			e = Rectangular(r.Angle, float64(radius))
			x, y := vg.Length(e.X+(off.X*2))+cen.X, vg.Length(e.Y+(off.Y*2))+cen.Y

			var (
				rot            Angle
				xalign, yalign float64
			)
			if r.Tick.Placement == nil {
				rot, xalign, yalign = DefaultPlacement(r.Angle)
			} else {
				rot, xalign, yalign = r.Tick.Placement(r.Angle)
			}
			if rot != 0 {
				ca.Push()
				ca.Translate(x, y)
				ca.Rotate(float64(rot))
				ca.Translate(-x, -y)
				ca.FillText(r.Tick.Label, x, y, xalign, yalign, mark.Label)
				ca.Pop()
			} else {
				ca.FillText(r.Tick.Label, x, y, xalign, yalign, mark.Label)
			}
		}
	}

	if r.Label.Text != "" && r.Label.Color != nil {
		e = Rectangular(r.Angle, float64(inner+outer)/2)
		x, y := vg.Length(e.X)+cen.X, vg.Length(e.Y)+cen.Y

		var (
			rot            Angle
			xalign, yalign float64
		)
		if r.Label.Placement == nil {
			rot, xalign, yalign = DefaultPlacement(r.Angle)
		} else {
			rot, xalign, yalign = r.Label.Placement(r.Angle)
		}
		if rot != 0 {
			ca.Push()
			ca.Translate(x, y)
			ca.Rotate(float64(rot))
			ca.Translate(-x, -y)
			ca.FillText(r.Label.TextStyle, x, y, xalign, yalign, r.Label.Text)
			ca.Pop()
		} else {
			ca.FillText(r.Label.TextStyle, x, y, xalign, yalign, r.Label.Text)
		}
	}
}
