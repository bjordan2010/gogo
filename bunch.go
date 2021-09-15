package main

// Bunch is a sub-group of connected stones
type Bunch struct {
	Player    Player
	Stones    []Point
	Liberties []Point
}

func (b *Bunch) RemoveLiberty(point Point) {
	goodLiberties := make([]Point, 0)
	for _, p := range b.Liberties {
		if p != point {
			goodLiberties = append(goodLiberties, p)
		}
	}
	b.Liberties = goodLiberties
}

// func (b *Bunch) containsStone(p *Point) bool {
// 	for _, x := range b.Stones {
// 		if x == p {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (b *Bunch) containsLiberty(p *Point) bool {
// 	for _, x := range b.Liberties {
// 		if x == p {
// 			return true
// 		}
// 	}
// 	return false
// }

func (b *Bunch) AddLiberty(point Point) {
	b.Liberties = append(b.Liberties, point)
}

func (b *Bunch) MergeWith(bunch Bunch) Bunch {
	if b.Player != bunch.Player {
		return Bunch{}
	}

	combinedStones := append(b.Stones, bunch.Stones...)
	return Bunch{
		Stones:    combinedStones,
		Player:    b.Player,
		Liberties: append(b.Liberties, bunch.Liberties...),
	}
}

func (b *Bunch) LibertyCount() int {
	return len(b.Liberties)
}
