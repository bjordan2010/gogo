package main

type Move struct {
	Point       *Point
	IsPlaying   bool
	IsPassing   bool
	IsResigning bool
}

func (m *Move) Play(point Point) {
	m.Point = &point
	m.IsPlaying = true
	m.IsPassing = false
	m.IsResigning = false
}

func (m *Move) Pass() {
	m.Point = nil
	m.IsPlaying = false
	m.IsPassing = true
	m.IsResigning = false
}

func (m *Move) Resign() {
	m.Point = nil
	m.IsPlaying = false
	m.IsPassing = false
	m.IsResigning = false
}
