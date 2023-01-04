//go:build rp2040
// +build rp2040

package main

import (
	"errors"
	"fmt"

	math "github.com/chewxy/math32"

	"machine"
	"time"
	pixel "zclock/pixel32"

	"golang.org/x/image/colornames"
	"tinygo.org/x/drivers/gc9a01"
)

const (
	RESETPIN = machine.GPIO12
	CSPIN    = machine.GPIO9
	DCPIN    = machine.GPIO8
	BLPIN    = machine.GPIO25

	// Default Serial Clock Bus 1 for SPI communications
	SPI1_SCK_PIN = machine.GPIO10
	// Default Serial Out Bus 1 for SPI communications
	SPI1_SDO_PIN = machine.GPIO11 // Tx
	// Default Serial In Bus 1 for SPI communications
	SPI1_SDI_PIN = machine.GPIO11 //machine.GPIO12 // Rx
)

type mcu struct {
	spi *machine.SPI
	lcd *gc9a01.Device
}

type spiral struct {
	imd      Sprite
	hands    Sprite
	links    []pixel.Vec
	bounds   pixel.Rect
	scale    float32
	numLinks int
	hourHand pixel.Line
	minHand  pixel.Line
}

func newSpiral() *spiral {
	return &spiral{
		// imd:      make([]uint16, 0),
		// hands:    make([]uint16, 0),
		hourHand: pixel.Line{},
		minHand:  pixel.Line{},
	}
}

func main() {
	run()
}

func configWindow() *mcu {
	d := mcu{}
	d.spi = machine.SPI1
	conf := machine.SPIConfig{
		Frequency: 40 * machine.MHz,
	}

	if err := d.spi.Configure(conf); err != nil {
		fmt.Println("error configuring spi:", err)
	}

	lcd := gc9a01.New(d.spi, RESETPIN, DCPIN, CSPIN, BLPIN)
	d.lcd = &lcd
	d.lcd.Configure(gc9a01.Config{})
	return &d
}

func run() {
	d := configWindow()
	s := newSpiral()

	var mi int
	imaginary := float32(216.8121)

	last := time.Unix(0, 0)

	for {
		now := time.Now()
		if time.Since(last) > time.Second {
			midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local) // % int64(time.Hour*24/time.Second)
			min := int(time.Since(midnight)/time.Minute) % int(12*time.Hour/time.Minute)

			t := timedata[min]
			imaginary = t.imaginary
			mi = t.index
			last = now

			d.lcd.FillScreen(colornames.Black)

			s.calc(.5, imaginary)
			s.drawSpiral(mi, .5, imaginary, d.lcd)
			s.drawHands(d.lcd)

			last = now
		}
	}
}

func (s *spiral) calc(real, imag float32) {
	s.bounds = pixel.Rect{}
	s.numLinks = int(imag/math.Pi + 1)
	s.links = make([]pixel.Vec, s.numLinks)

	start := pixel.V(0, 0)
	s.links[0] = start

	for i := 1; i < s.numLinks; i++ {
		x := math.Cos(imag*math.Log(float32(i))) / math.Pow(float32(i), real)
		y := -math.Sin(imag*math.Log(float32(i))) / math.Pow(float32(i), real)
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

func (s *spiral) drawSpiral(mi int, real, imag float32, lcd *gc9a01.Device) {
	lcd.FillRectangle(50, 50, 100, 100, colornames.Orange)
	// s.imd.Clear()

	// s.imd.EndShape = imdraw.NoEndShape

	// s.imd.Color = pixel.RGB(.5, .5, .5)
	// s.imd.Push(s.links[:mi]...)

	// pt := s.links[mi]
	// s.minHand.A = pt
	// s.imd.Push(pt) // pushing an extra point makes a hard stop from color blending
	// s.imd.Color = pixel.RGB(1, .5, 0)

	// pt = s.links[mi+1]
	// s.minHand.B = pt
	// s.hourHand.A = pt
	// s.imd.Push(pt, pt)
	// s.imd.Color = pixel.RGB(1, 0, 0)

	// pt = s.links[mi+2]
	// s.hourHand.B = pt
	// s.imd.Push(pt, pt)

	// s.imd.Color = pixel.RGB(.5, .5, .5)
	// s.imd.Push(s.links[mi+2:]...)

	// max := math.Max(5, math.Max(s.bounds.W(), s.bounds.H())) // is width or height bigger?
	// width := win.Bounds().W() / max
	// mat := pixel.IM.Scaled(pixel.ZV, width)
	// mat = mat.Moved(win.Bounds().Center().Sub(s.bounds.Center().Scaled(width))) // center it on the screen
	// s.imd.SetMatrix(mat)

	// s.imd.Line(.05)
	// s.imd.Draw(win)
}

func (s *spiral) drawHands(lcd *gc9a01.Device) {
	// s.hands.Clear()
	// s.hands.EndShape = imdraw.RoundEndShape

	// scale := win.Bounds().W() / 2 * .6 / s.hourHand.Len()
	// mat := pixel.IM.Scaled(pixel.ZV, scale) //win.Bounds().W())
	// mat = mat.Moved(win.Bounds().Center().Sub(s.hourHand.A.Scaled(scale)))

	// s.hands.Color = pixel.RGB(1, 0, 0)
	// s.hands.Push(s.hourHand.A, s.hourHand.B)
	// // s.hands.SetMatrix(mat.Scaled(pixel.ZV, scale*1.5))
	// s.hands.SetMatrix(mat)
	// s.hands.Line(.05)
	// s.hands.Draw(win)

	// scale = win.Bounds().W() / 2 / s.minHand.Len() * .9
	// mat = pixel.IM.Scaled(pixel.ZV, scale) //win.Bounds().W())
	// mat = mat.Moved(win.Bounds().Center().Sub(s.minHand.B.Scaled(scale)))

	// s.hands.Color = pixel.RGB(1, .5, 0)
	// s.hands.Push(s.minHand.A, s.minHand.B)
	// s.hands.SetMatrix(mat)
	// s.hands.Line(.05 * .7)
	// s.hands.Draw(win)
}

func (d *mcu) Size() (int16, int16) {
	return d.lcd.Size()
}

func (d *mcu) Draw(s *Sprite) error {
	return errors.New("Not implemented")
}
