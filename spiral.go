package main

import (
	math "github.com/chewxy/math32"
	"golang.org/x/image/colornames"

	"tinygo.org/x/drivers"
	"tinygo.org/x/tinydraw"
)

type spiral struct {
	joints    []Vec
	bounds    Rect
	scale     float32
	numJoints int
	hourHand  Line
	minHand   Line
}

func (s *spiral) calc(real, imag float32) {
	s.bounds = Rect{}
	s.numJoints = int(imag/math.Pi + 1)

	if s.joints == nil || len(s.joints) < s.numJoints {
		s.joints = make([]Vec, s.numJoints*2)
	}

	start := V(0, 0)
	s.joints[0] = start

	for i := float32(1.0); i < float32(s.numJoints); i++ {
		x := float32(math.Cos(imag*math.Log(i)) / math.Pow(i, real))
		y := float32(math.Sin(imag*math.Log(i)) / math.Pow(i, real))
		end := V(start.X+x, start.Y+y)

		if end.X < s.bounds.Min.X {
			s.bounds.Min.X = end.X
		}
		if end.Y < s.bounds.Min.Y {
			s.bounds.Min.Y = end.Y
		}
		if end.X > s.bounds.Max.X {
			s.bounds.Max.X = end.X
		}
		if end.Y > s.bounds.Max.Y {
			s.bounds.Max.Y = end.Y
		}

		s.joints[int(i)] = end
		start = end
	}
}

func drawSpiral(lcd drivers.Displayer, s *spiral, mi int) {
	max := math.Max(s.bounds.W(), s.bounds.H()) // is width or height bigger?
	scale := float32(188.49 / max)              // pi / 4 * 240 keeps it out of the corners
	mov := V(128, 128).Sub(s.bounds.Center().Scaled(float32(scale)))

	// mat := IM.Scaled(ZV, scale)
	// mat = mat.Moved(win.Bounds().Center().Sub(s.bounds.Center().Scaled(width))) // center it on the screen
	// s.imd.SetMatrix(mat)

	for i := 1; i < s.numJoints; i++ {
		from := s.joints[i-1]
		to := s.joints[i]
		c := colornames.Dimgray

		if i == mi+1 {
			c = colornames.Orange
		} else if i == mi {
			c = colornames.Green
		} else if i == mi+2 {
			c = colornames.Red
		}

		tinydraw.Line(lcd, int16(from.X*scale+mov.X), int16(from.Y*scale+mov.Y),
			int16(to.X*scale+mov.X), int16(to.Y*scale+mov.Y), c)
	}
}

func drawHands(lcd drivers.Displayer, minHand, hourHand Line) {
	hl := hourHand.Len()
	scale := 128 * .6 / hl
	center := V(128, 128)
	mat := IM.Scaled(ZV, scale) //win.Bounds().W())
	mat = mat.Moved(center.Sub(hourHand.A.Scaled(scale)))

	line := Line{
		A: mat.Project(hourHand.A),
		B: mat.Project(hourHand.B),
	}
	pt := V(line.A.X+2, line.A.Y+2)

	tinydraw.FilledTriangle(lcd, int16(line.A.X), int16(line.A.Y),
		int16(line.B.X), int16(line.B.Y), int16(pt.X), int16(pt.Y), colornames.Red)

	scale = 128 / minHand.Len() * .9
	mat = IM.Scaled(ZV, scale) //win.Bounds().W())
	mat = mat.Moved(center.Sub(minHand.B.Scaled(scale)))

	line = Line{
		A: mat.Project(minHand.A),
		B: mat.Project(minHand.B),
	}
	pt = V(line.B.X+2, line.B.Y+2)
	line.B.X -= 2
	line.B.Y -= 2
	// tinydraw.Line(lcd, int16(line.A.X), int16(line.A.Y), int16(line.B.X), int16(line.B.Y), colornames.Orange)
	tinydraw.FilledTriangle(lcd, int16(line.A.X), int16(line.A.Y),
		int16(line.B.X), int16(line.B.Y), int16(pt.X), int16(pt.Y), colornames.Orange)

}
