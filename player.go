package main

type Player uint8

const (
	Empty Player = iota // unused
	Black
	White
)

func (p *Player) Other() Player {
	if *p == Black {
		return Black
	}

	return White
}

func (p *Player) Val() uint8 {
	switch *p {
	case Empty:
		return 0
	case Black:
		return 1
	}

	return 2
}

func (p *Player) String() string {
	switch *p {
	case Empty:
		return "Empty"
	case Black:
		return "Black"
	}

	return "White"
}
