package zeta

import "github.com/faiface/pixel"

// getSpiralPos returns the spiral indicated by num where 0 is the last spiral
// and working backward.
//
// zetaPt is Zeta in vector form
func GetSpiralPos(num int, spiral []pixel.Vec, zetaPt pixel.Vec) pixel.Vec {
	// slope := -z.x / z.y
	// bipt := bisectPoint(z)

	z2 := zetaPt.Scaled(.5)

	// draw a line from each of the first links at the same slope as zeta
	from := spiral[num]
	norm := z2.Unit()
	dot := from.Dot(norm) // from.X*norm.X + from.Y*norm.Y //Vector2.Dot(from, norm)

	return zetaPt.Add(from).Sub(norm.Scaled(2 * dot)) // reflect from about a normal (z2)
}

// bisectPoint finds the point where the bisecting line intersects the middle link
func BisectPoint(mi int, zetaPt pixel.Vec, spiral []pixel.Vec) pixel.Vec {
	M1 := spiral[mi]
	M2 := spiral[mi+1]

	slope1 := -zetaPt.X / zetaPt.Y
	slope2 := (M2.Y - M1.Y) / (M2.X - M1.X)

	x := ((slope2*M2.X - slope1*zetaPt.X/2) - (M2.Y - zetaPt.Y/2)) / (slope2 - slope1)

	y := slope1*(x-zetaPt.X/2) + zetaPt.Y/2

	return pixel.Vec{X: x, Y: y}
}
