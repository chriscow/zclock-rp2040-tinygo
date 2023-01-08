package main

import (
	"image/color"
	"machine"

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
	spi    *machine.SPI
	lcd    *gc9a01.Device
	buffer []color.RGBA
}

func newDisplay() (*mcu, error) {
	spi := machine.SPI1
	conf := machine.SPIConfig{
		Frequency: 40 * machine.MHz,
	}

	if err := spi.Configure(conf); err != nil {
		return nil, err
	}

	lcd := gc9a01.New(spi, RESETPIN, DCPIN, CSPIN, BLPIN)
	lcd.Configure(gc9a01.Config{})

	width, height := 240, 240
	// fmt.Println(width, height)
	buffer := make([]color.RGBA, width*height)
	return &mcu{spi: spi, lcd: &lcd, buffer: buffer}, nil
}

func (d *mcu) Size() (int16, int16) {
	return d.lcd.Size()
}

func (d *mcu) SetPixel(x, y int16, c color.RGBA) {
	if x < 0 || x > 239 || y < 0 || y > 239 {
		return
	}

	i := int(y)*240 + int(x) // upcast to prevent overflow

	if i >= len(d.buffer) {
		return
	}

	d.buffer[i] = c
}

func (d *mcu) Display() error {
	w, h := d.lcd.Size()
	return d.lcd.FillRectangleWithBuffer(0, 0, w, h, d.buffer)
}

func (d *mcu) FillDisplay(c color.RGBA) {
	for i := range d.buffer {
		d.buffer[i] = c
	}
}
