//go:build rp2040
// +build rp2040

package main

import (
	math "github.com/chewxy/math32"
)

func ImagToIndex(imag float32) float32 {
	//best so far -- this is from Zzrob
	gamma := float32(0.57721566490153286060651209008240243104215933593992)
	e := float32(2.7182818284590452353602874713526624977572)
	gamma_to_the_e := math.Pow(gamma, e) // = .2245172519832320
	two_root_3_pi := 2 * math.Sqrt(3*math.Pi)
	return_this := math.Sqrt(6*gamma_to_the_e/imag+6*imag+math.Pi)/two_root_3_pi - 1.0/2.0

	return (return_this)
}

func IndexToImag(n float32) float32 {
	return (n*2 + 1) * math.Pi / (math.Log(n+1) - math.Log(n))
}
