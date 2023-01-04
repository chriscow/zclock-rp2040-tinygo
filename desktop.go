//go:build !rp2040
// +build !rp2040

package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"
	"zclock/zeta"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type desktop struct {
	win *pixelgl.Window
}

func (d *desktop) Size() (int16, int16) {
	b := d.win.Bounds()
	return int16(b.Max.X), int16(b.Max.Y)
}

func (d *desktop) Draw(s *Sprite) error {
	return errors.New("not implemented")
}

func main() {
	times = map[int]data{}

	// calcData()
	if err := loadTimeData(); err != nil {
		fmt.Println("error loading data:", err)
	}

	pixelgl.Run(run)
}

type spiral struct {
	imd      *imdraw.IMDraw
	hands    *imdraw.IMDraw
	zeta     *imdraw.IMDraw
	links    []pixel.Vec
	bounds   pixel.Rect
	scale    float64
	numLinks int
	hourHand pixel.Line
	minHand  pixel.Line
}

func newSpiral() *spiral {
	return &spiral{
		imd:      imdraw.New(nil),
		hands:    imdraw.New(nil),
		zeta:     imdraw.New(nil),
		hourHand: pixel.Line{},
		minHand:  pixel.Line{},
	}
}

func (s *spiral) calc(real, imag float64) {
	s.bounds = pixel.Rect{}
	s.numLinks = int(imag/math.Pi + 1)
	s.links = make([]pixel.Vec, s.numLinks)

	start := pixel.V(0, 0)
	s.links[0] = start

	for i := 1; i < s.numLinks; i++ {
		x := math.Cos(imag*math.Log(float64(i))) / math.Pow(float64(i), real)
		y := -math.Sin(imag*math.Log(float64(i))) / math.Pow(float64(i), real)
		end := pixel.V(start.X+x, start.Y+y)

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

		s.links[i] = end
		start = end
	}
}

func (s *spiral) drawSpiral(mi int, real, imag float64, win *pixelgl.Window) {
	s.imd.Clear()

	s.imd.EndShape = imdraw.NoEndShape

	s.imd.Color = pixel.RGB(.5, .5, .5)
	s.imd.Push(s.links[:mi]...)

	pt := s.links[mi]
	s.minHand.A = pt
	s.imd.Push(pt) // pushing an extra point makes a hard stop from color blending
	s.imd.Color = pixel.RGB(1, .5, 0)

	pt = s.links[mi+1]
	s.minHand.B = pt
	s.hourHand.A = pt
	s.imd.Push(pt, pt)
	s.imd.Color = pixel.RGB(1, 0, 0)

	pt = s.links[mi+2]
	s.hourHand.B = pt
	s.imd.Push(pt, pt)

	s.imd.Color = pixel.RGB(.5, .5, .5)
	s.imd.Push(s.links[mi+2:]...)

	max := math.Max(5, math.Max(s.bounds.W(), s.bounds.H())) // is width or height bigger?
	width := win.Bounds().W() / max
	mat := pixel.IM.Scaled(pixel.ZV, width)
	mat = mat.Moved(win.Bounds().Center().Sub(s.bounds.Center().Scaled(width))) // center it on the screen
	s.imd.SetMatrix(mat)

	s.imd.Line(.05)
	s.imd.Draw(win)
}

func (s *spiral) drawZeta(real, imag float64, zetaPt pixel.Vec, win *pixelgl.Window) {
	s.zeta.Clear()
	s.zeta.EndShape = imdraw.RoundEndShape

	s.zeta.Color = pixel.RGBA{R: 0, G: 1, B: 0, A: .3}

	link := s.links[1] // this is the 1st link numLinks
	center := zeta.GetSpiralPos(1, []pixel.Vec{pixel.ZV, link}, zetaPt)

	inside := false
	target := pixel.Circle{
		Center: center,
		Radius: .3,
	}

	var i int
	for i = len(s.links) - 2; i > 0; i-- {
		line := pixel.Line{
			A: s.links[i-1],
			B: s.links[i],
		}

		if !inside && target.IntersectLine(line) != pixel.ZV {
			inside = true
		}

		if inside && target.IntersectLine(line) == pixel.ZV {
			break
		}

		if i <= 1 {
			break
		}
	}
	s.zeta.Push(s.links[i:]...)

	max := 3.0 //math.Max(5, math.Max(s.bounds.W(), s.bounds.H())) // is width or height bigger?
	width := win.Bounds().W() / max
	mat := pixel.IM.Scaled(pixel.ZV, width)
	mat = mat.Moved(win.Bounds().Center().Sub(center.Scaled(width))) // center it on the screen
	s.zeta.SetMatrix(mat)

	s.zeta.Line(.05)
	//
	// trying to figure out why the spiral position is off
	//
	// s.zeta.Color = colornames.Olive
	// s.zeta.Push(target.Center)
	// s.zeta.SetMatrix(mat)
	// s.zeta.Circle(target.Radius, .01)

	// s.zeta.Color = colornames.Azure
	// s.zeta.Push(zetaPt)
	// s.zeta.SetMatrix(mat)
	// s.zeta.Circle(target.Radius, .01)

	// s.zeta.Color = colornames.Navy
	// s.zeta.Push(s.links[1], center)
	// s.zeta.SetMatrix(mat)
	// s.zeta.Line(.05)

	s.zeta.Draw(win)
}

func (s *spiral) drawHands(win *pixelgl.Window) {
	s.hands.Clear()
	s.hands.EndShape = imdraw.RoundEndShape

	scale := win.Bounds().W() / 2 * .6 / s.hourHand.Len()
	mat := pixel.IM.Scaled(pixel.ZV, scale) //win.Bounds().W())
	mat = mat.Moved(win.Bounds().Center().Sub(s.hourHand.A.Scaled(scale)))

	s.hands.Color = pixel.RGB(1, 0, 0)
	s.hands.Push(s.hourHand.A, s.hourHand.B)
	// s.hands.SetMatrix(mat.Scaled(pixel.ZV, scale*1.5))
	s.hands.SetMatrix(mat)
	s.hands.Line(.05)
	s.hands.Draw(win)

	scale = win.Bounds().W() / 2 / s.minHand.Len() * .9
	mat = pixel.IM.Scaled(pixel.ZV, scale) //win.Bounds().W())
	mat = mat.Moved(win.Bounds().Center().Sub(s.minHand.B.Scaled(scale)))

	s.hands.Color = pixel.RGB(1, .5, 0)
	s.hands.Push(s.minHand.A, s.minHand.B)
	s.hands.SetMatrix(mat)
	s.hands.Line(.05 * .7)
	s.hands.Draw(win)
}

func configWindow() desktop {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 240, 240),
		VSync:  true,
	}
	var err error
	d := desktop{}
	d.win, err = pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatal(err)
	}
	d.win.SetSmooth(true)

	return d
}

func run() {
	d := configWindow()
	s := newSpiral()

	var mi int
	imaginary := 216.8121

	lastMin := time.Unix(0, 0)
	lastSec := time.Unix(0, 0)
	// min := 0
	// trouble: min == 6

	for !d.win.Closed() {
		now := time.Now()
		if time.Since(lastMin) > time.Millisecond*500 {
			midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local) // % int64(time.Hour*24/time.Second)
			min := int(time.Since(midnight)/time.Minute) % int(12*time.Hour/time.Minute)

			t := timedata[min]
			imaginary = float64(t.imaginary)
			mi = t.index
			lastMin = now

			// min++
			// fmt.Println(min, mi)
		}

		if time.Since(lastSec) > time.Millisecond*500 {

			// z := zeta.EulerMaclaurin(complex(.5, imaginary))
			// zetaPt := pixel.V(real(z), imag(z))

			d.win.Clear(colornames.Black)

			s.calc(.5, imaginary)
			s.drawSpiral(mi, .5, imaginary, d.win)
			// s.drawZeta(.5, imaginary, zetaPt, d.win)
			s.drawHands(d.win)

			d.win.Update()
			// imaginary += .05
			lastSec = now
		}
	}
}
