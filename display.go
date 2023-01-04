package main

type Display interface {
	Size() (int16, int16)
	Draw(s *Sprite) error
}

type Sprite interface {
	Clear()
}
