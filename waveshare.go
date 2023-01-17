package main

import (
	"errors"
	"image/color"
	"machine"

	"github.com/chewxy/math32"
	"golang.org/x/image/colornames"
	"tinygo.org/x/drivers/gc9a01"
	imu "tinygo.org/x/drivers/qmi8658c"
	"tinygo.org/x/tinydraw"
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

// RP2040 MCU Board, With 1.28inch Round LCD, accelerometer and gyroscope Sensor
// https://www.waveshare.com/rp2040-lcd-1.28.htm
type waveshare struct {
	spi    *machine.SPI
	lcd    *gc9a01.Device
	imu    *imu.Device
	adc    machine.ADC
	buffer []color.RGBA

	minVolts float32
}

func newMCU() (*waveshare, error) {
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
	buffer := make([]color.RGBA, width*height)

	machine.InitSerial()

	// for reading battery voltage
	machine.InitADC()
	adc := machine.ADC{Pin: machine.GPIO29}
	adc.Configure(machine.ADCConfig{})

	// Configure the IMU over I2C (inertial measurement unit)
	i2c := machine.I2C1
	// This is the default pinout for the "WaveShare RP2040 Round LCD 1.28in"
	err := i2c.Configure(machine.I2CConfig{
		SDA:       machine.GP6,
		SCL:       machine.GP7,
		Frequency: 100000,
	})
	if err != nil {
		return nil, errors.New("unable to configure I2C:" + err.Error())
	}

	imud := imu.New(i2c)
	if !imud.Connected() {
		return nil, errors.New("unable to connect to sensor")
	}

	// This IMU has multiple configurations like output data rate, multiple
	// measurements scales, low pass filters, low power modes, all the vailable
	// values can be found in the datasheet and were defined at registers file.
	// This is the default configuration which will be used if the `nil` value
	// is passed do the `Configure` method.
	imud.Configure(imu.Config{})

	return &waveshare{
		spi:      spi,
		lcd:      &lcd,
		imu:      &imud,
		adc:      adc,
		buffer:   buffer,
		minVolts: math32.MaxFloat32,
	}, nil
}

func (d *waveshare) Size() (int16, int16) {
	return d.lcd.Size()
}

func (d *waveshare) SetPixel(x, y int16, c color.RGBA) {
	if x < 0 || x > 239 || y < 0 || y > 239 {
		return
	}

	i := int(y)*240 + int(x) // upcast to prevent overflow

	if i >= len(d.buffer) {
		return
	}

	d.buffer[i] = c
}

func (d *waveshare) Display() error {
	w, h := d.lcd.Size()
	return d.lcd.FillRectangleWithBuffer(0, 0, w, h, d.buffer)
}

func (d *waveshare) FillDisplay(c color.RGBA) {
	for i := range d.buffer {
		d.buffer[i] = c
	}
}

// Volts returns the raw voltage reading of the battery.
func (d *waveshare) Volts() float32 {
	raw := d.adc.Get()

	// 3.3v max from the pin / 4096 steps
	// waveshare biases to 1/2 voltage so double it
	return 3.3 / 4096.0 * 2.0 * float32(raw>>4)
}

// DrawBattery displays the battery if the voltage has dropped below 3.6 volts
// and changes to orange then red to represent charge status.
func (d *waveshare) DrawBattery() {
	d.minVolts = math32.Min(d.minVolts, d.Volts())

	volts := d.minVolts

	if volts >= 4 {
		// TODO: draw lightning bolt indicating charging status
		// fully charged or charging
		d.minVolts = math32.MaxFloat32
		tinydraw.Rectangle(d, 117, 225, 20, 10, colornames.Grey)     // battery body
		tinydraw.FilledRectangle(d, 137, 228, 3, 4, colornames.Grey) // positive terminal
		tinydraw.FilledRectangle(d, 118, 226, 19, 8, colornames.Green)
		return
	}

	// if d.minVolts < 4 {
	// draw outline only when battery is below

	// TODO: parameterize battery position and size
	// battery width: 20, height: 10
	tinydraw.Rectangle(d, 117, 225, 20, 10, colornames.Grey)     // battery body
	tinydraw.FilledRectangle(d, 137, 228, 3, 4, colornames.Grey) // positive terminal
	// }

	if d.minVolts < 3.4 {
		tinydraw.FilledRectangle(d, 118, 226, 8, 8, colornames.Red)
	} else if d.minVolts < 3.5 {
		tinydraw.FilledRectangle(d, 118, 226, 10, 8, colornames.Orange)
	} else if d.minVolts < 3.6 {
		tinydraw.FilledRectangle(d, 118, 226, 10, 8, colornames.Grey)
	} else if d.minVolts < 3.7 {
		tinydraw.FilledRectangle(d, 118, 226, 12, 8, colornames.Grey)
	} else if d.minVolts < 3.8 {
		tinydraw.FilledRectangle(d, 118, 226, 14, 8, colornames.Grey)
	} else if d.minVolts < 3.9 {
		tinydraw.FilledRectangle(d, 118, 226, 16, 8, colornames.Grey)
	} else {
		tinydraw.FilledRectangle(d, 118, 226, 18, 8, colornames.Grey)
	}
}
