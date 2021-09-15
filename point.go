package main

type Point struct {
	Row int
	Col int
}

func (p Point) Neighbors() []Point {
	points := make([]Point, 0)
	points = append(points, Point{p.Row - 1, p.Col})
	points = append(points, Point{p.Row + 1, p.Col})
	points = append(points, Point{p.Row, p.Col - 1})
	points = append(points, Point{p.Row, p.Col + 1})

	return points
}
