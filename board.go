package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"golang.org/x/image/draw"
)

const maxBoard = 19 // Maximum board size we can handle.

var ZP = image.Point{}

type stoneColor uint8

const (
	blank stoneColor = iota // Unused
	black
	white
)

func (s *stoneColor) val() uint8 {
	switch *s {
	case blank:
		return 0
	case black:
		return 1
	}

	return 2
}

// Piece represents a stone on the board. A nil Piece is "blank".
// The delta records pixel offset from the central dot.
type Piece struct {
	stone *Stone
	ij    IJ
	delta image.Point
	color stoneColor
}

type Board struct {
	Dims
	pieces         []*Piece // The board. Dimensions are 1-indexed. 1, 1 is the lower left corner.
	image          *image.RGBA
	stone          []Stone // All the black stones, followed by all the white stones.
	numBlackStones int
	numWhiteStones int
	grid           map[Point]*Bunch // map of Bunches by Point
}

type Stone struct {
	originalImage *image.RGBA
	originalMask  *image.Alpha
	image         *image.RGBA
	mask          *image.Alpha
}

func NewBoard(dim, percent int) *Board {
	switch dim {
	case 9, 19: // removed 13 and 19
	default:
		return nil
	}
	boardTexture := get("goboard.jpg", 0)
	b := new(Board)
	b.Dims.Init(dim, 100)
	b.pieces = make([]*Piece, maxBoard*maxBoard)
	b.image = image.NewRGBA(boardTexture.Bounds())
	draw.Draw(b.image, b.image.Bounds(), boardTexture, ZP, draw.Src)
	dir, err := os.Open("asset")
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()
	names, err := dir.Readdirnames(0)
	if err != nil {
		log.Fatal(err)
	}
	circleMask := makeCircle()
	// Blackstones go first
	for _, name := range names {
		if strings.HasPrefix(name, "blackstone") {
			s, m := makeStone(name, circleMask)
			b.stone = append(b.stone, Stone{s, m, nil, nil})
			b.numBlackStones++
		}
	}
	for _, name := range names {
		if strings.HasPrefix(name, "whitestone") {
			s, m := makeStone(name, circleMask)
			b.stone = append(b.stone, Stone{s, m, nil, nil})
			b.numWhiteStones++
		}
	}
	b.Resize(percent)
	b.grid = make(map[Point]*Bunch)
	return b
}

func (b *Board) Resize(percent int) {
	b.Dims.Resize(percent)
	for i := range b.stone {
		stone := &b.stone[i]
		stone.image = resizeRGBA(stone.originalImage, b.stoneDiam)
		stone.mask = resizeAlpha(stone.originalMask, b.stoneDiam)
	}
}

func resizeRGBA(src *image.RGBA, size int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Src, nil)
	return dst
}

func resizeAlpha(src *image.Alpha, size int) *image.Alpha {
	dst := image.NewAlpha(image.Rect(0, 0, size, size))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Src, nil)
	return dst
}

func (b *Board) piece(ij IJ) *Piece {
	return b.pieces[(ij.j-1)*b.Dims.dim+ij.i-1]
}

func jitter() int {
	max := 25 * *scale / 100
	if max&1 == 0 {
		max++
	}
	return rand.Intn(max) - max/2
}

func (b *Board) putPiece(ij IJ, piece *Piece) {
	b.pieces[(ij.j-1)*b.Dims.dim+ij.i-1] = piece
	if piece != nil {
		piece.ij = ij
		piece.delta = image.Point{jitter(), jitter()}

		// Game State
		var player Player
		if Black == Player(piece.color.val()) {
			player = Black
		} else {
			player = White
		}
		point := Point{ij.i, ij.j}
		stonesToRemove := b.placeStone(player, point)

		// Game rule automatic remove stone for rendering
		for _, p := range stonesToRemove {
			b.putPiece(IJ{i: p.Row, j: p.Col}, nil)
		}
	}
}

// Game State
func (b *Board) placeStone(player Player, point Point) []Point {
	stonesToRemove := make([]Point, 0)
	if 1 < point.Row && point.Col > b.dim {
		log.Fatal("Point row outside board cannot place stone")
		return stonesToRemove
	}
	if 1 < point.Col && point.Col > b.dim {
		log.Fatal("Point col outside board cannot place stone")
		return stonesToRemove
	}
	var friendly_bunch []Bunch
	var enemy_bunch []Bunch
	var liberties []Point
	for _, n := range point.Neighbors() {
		if 1 < n.Row && n.Row > b.dim {
			log.Printf("Neighbor row %d outside board cannot place stone\n", n.Row)
			continue
		}
		if 1 < n.Col && n.Col > b.dim {
			log.Printf("Neighbor col %d outside board cannot place stone\n", n.Col)
			continue
		}

		// if neighbor point is in bunch
		if bunch, ok := b.grid[n]; ok {
			if bunch.Player == player && !contains(friendly_bunch, *bunch) {
				friendly_bunch = append(friendly_bunch, *bunch)
			} else {
				if !contains(enemy_bunch, *bunch) {
					enemy_bunch = append(enemy_bunch, *bunch)
				}
			}
		} else {
			liberties = append(liberties, n)
		}
	}
	stones := make([]Point, 0)
	stones = append(stones, point)
	bunch := Bunch{Player: player, Stones: stones, Liberties: liberties}
	for _, bu := range friendly_bunch {
		bunch.MergeWith(bu)
	}
	for _, p := range bunch.Stones {
		b.grid[p] = &bunch
	}
	for _, bu := range b.grid {
		for _, p := range bu.Liberties {
			if p == point && bu.Player != player {
				bu.RemoveLiberty(point)
			}
		}
	}

	for p, bu := range b.grid {
		if bu.LibertyCount() == 0 {
			for _, n := range p.Neighbors() {
				if val, ok := b.grid[n]; ok {
					val.AddLiberty(p)
				}
			}
			stonesToRemove = append(stonesToRemove, p)
		}
	}

	for _, p := range stonesToRemove {
		delete(b.grid, p)
	}

	fmt.Println("Placing Stone Info")
	fmt.Printf("Placing %s stone at (%d,%d)\n\n", player.String(), point.Row, point.Col)
	for i, bu := range b.grid {
		fmt.Printf("Bunch (%d,%d):\n", i.Row, i.Col)
		fmt.Printf("Player: %s\n", bu.Player.String())
		for j, p := range bu.Stones {
			fmt.Printf("Stone Point %d: (%d,%d)\n", j+1, p.Row, p.Col)
		}
		for j, p := range bu.Liberties {
			fmt.Printf("Liberty Point %d: (%d,%d)\n", j+1, p.Row, p.Col)
		}
	}
	fmt.Print("\n\n\n")

	return stonesToRemove
}

func contains(s []Bunch, bunch Bunch) bool {
	for _, b := range s {
		if reflect.DeepEqual(b, bunch) {
			return true
		}
	}

	return false
}

func (b *Board) selectBlackPiece() *Piece {
	return &Piece{
		stone: &b.stone[rand.Intn(b.numBlackStones)],
		color: black,
	}
}

func (b *Board) selectWhitePiece() *Piece {
	return &Piece{
		stone: &b.stone[b.numBlackStones+rand.Intn(b.numWhiteStones)],
		color: white,
	}
}

func makeStone(name string, circleMask *image.Alpha) (*image.RGBA, *image.Alpha) {
	stone := get(name, stoneSize0)
	dst := image.NewRGBA(stone.Bounds())
	// Make the whole area black, for the shadow.
	draw.Draw(dst, dst.Bounds(), image.Black, ZP, draw.Src)
	// Lay in the stone within the circle so it shows up inside the shadow.
	draw.DrawMask(dst, dst.Bounds(), stone, ZP, circleMask, ZP, draw.Over)
	return dst, makeShadowMask(stone)
}

func get(name string, size int) image.Image {
	f, err := os.Open(filepath.Join("asset", name))
	if err != nil {
		log.Fatal(err)
	}
	i, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	if size != 0 {
		r := i.Bounds()
		if r.Dx() != size || r.Dy() != size {
			log.Fatalf("bad stone size %s for %s; must be %d[2]×%d[2]", r, name, size, size)
		}
	}
	return i
}

func makeCircle() *image.Alpha {
	mask := image.NewAlpha(image.Rect(0, 0, stoneSize0, stoneSize0))
	// Make alpha work on stone.
	// Shade gives shape, to be applied with black.
	for y := 0; y < stoneSize0; y++ {
		y2 := stoneSize0/2 - y
		y2 *= y2
		for x := 0; x < stoneSize0; x++ {
			x2 := stoneSize0/2 - x
			x2 *= x2
			if x2+y2 <= stoneRad2 {
				mask.SetAlpha(x, y, color.Alpha{255})
			}
		}
	}
	return mask
}

func makeShadowMask(stone image.Image) *image.Alpha {
	mask := image.NewAlpha(stone.Bounds())
	// Make alpha work on stone.
	// Shade gives shape, to be applied with black.
	const size = 256
	//const diam = 225
	for y := 0; y < size; y++ {
		y2 := size/2 - y
		y2 *= y2
		for x := 0; x < size; x++ {
			x2 := size/2 - x
			x2 *= x2
			if x2+y2 > stoneRad2 {
				red, _, _, _ := stone.At(x, y).RGBA()
				mask.SetAlpha(x, y, color.Alpha{255 - uint8(red>>8)})
			} else {
				mask.SetAlpha(x, y, color.Alpha{255})
			}
		}
	}
	return mask
}

func (b *Board) Draw(m *image.RGBA) {
	r := b.image.Bounds()
	draw.Draw(m, r, b.image, ZP, draw.Src)
	// Vertical lines.
	x := b.xInset + b.squareWidth/2
	y := b.yInset + b.squareHeight/2
	wid := b.lineWidth
	for i := 0; i < b.dim; i++ {
		r := image.Rect(x, y, x+wid, y+(b.dim-1)*b.squareHeight)
		draw.Draw(m, r, image.Black, ZP, draw.Src)
		x += b.squareWidth
	}
	// Horizontal lines.
	x = b.xInset + b.squareWidth/2
	for i := 0; i < b.dim; i++ {
		r := image.Rect(x, y, x+(b.dim-1)*b.squareWidth+wid, y+wid)
		draw.Draw(m, r, image.Black, ZP, draw.Src)
		y += b.squareHeight
	}
	// Points.
	spot := 4
	if b.dim < 13 {
		spot = 3
	}
	points := []IJ{
		{spot, spot},
		{spot, (b.dim + 1) / 2},
		{spot, b.dim + 1 - spot},
		{(b.dim + 1) / 2, spot},
		{(b.dim + 1) / 2, (b.dim + 1) / 2},
		{(b.dim + 1) / 2, b.dim + 1 - spot},
		{b.dim + 1 - spot, spot},
		{b.dim + 1 - spot, (b.dim + 1) / 2},
		{b.dim + 1 - spot, b.dim + 1 - spot},
	}
	for _, ij := range points {
		b.drawPoint(m, ij)
	}
	// Pieces.
	for i := 1; i <= b.dim; i++ {
		for j := 1; j <= b.dim; j++ {
			ij := IJ{i, j}
			if p := b.piece(ij); p != nil {
				b.drawPiece(m, ij, p)
			}
		}
	}
}

func (b *Board) drawPoint(m *image.RGBA, ij IJ) {
	pt := ij.XYCenter(&b.Dims)
	wid := b.lineWidth
	sz := wid * 3 / 2
	r := image.Rect(pt.x-sz, pt.y-sz, pt.x+wid+sz, pt.y+wid+sz)
	draw.Draw(m, r, image.Black, ZP, draw.Src)
}

func (b *Board) drawPiece(m *image.RGBA, ij IJ, piece *Piece) {
	xy := ij.XYStone(&b.Dims)
	xy = xy.Add(piece.delta)
	draw.DrawMask(m, xy, piece.stone.image, ZP, piece.stone.mask, ZP, draw.Over)
}

func (b *Board) click(m *image.RGBA, x, y, button int) bool {
	ij, ok := XY{x, y}.IJ(&b.Dims)
	if !ok {
		return false
	}
	switch button {
	default:
		return false
	case 1:
		b.putPiece(ij, b.selectBlackPiece())
	case 3:
		b.putPiece(ij, b.selectWhitePiece())
	case 2:
		b.putPiece(ij, nil)
	}
	render(m, b) // TODO: Connect this to paint events.
	return true
}
