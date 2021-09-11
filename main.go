package main

import (
	"image"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

var scale *int //= flag.Int("scale", 35, "`percent` to scale images 35 or 18")

//var dim = flag.Int("dim", 9, "`size of board 9 or 19")

func main() {
	//flag.Parse()

	rand.Seed(int64(time.Now().Nanosecond()))
	dimStr := os.Args[1]
	if dimStr == "" {
		dimStr = "9"
	}
	dim, err := strconv.Atoi(dimStr)
	if err != nil {
		log.Fatal(err)
	}
	scaleStr := os.Args[2]
	if scaleStr == "" {
		scaleStr = "35"
	}
	scaleInt, err := strconv.Atoi(scaleStr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Board: %dx%d", dim, dim)
	log.Printf("Scale: %d%%", scaleInt)
	scale := &scaleInt
	board := NewBoard(dim, *scale)

	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Title: "Goban Shiny Example",
		})
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()

		var b screen.Buffer
		defer func() {
			if b != nil {
				b.Release()
			}
		}()

		for {
			switch e := w.NextEvent().(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}

			case key.Event:
				if e.Code == key.CodeEscape {
					return
				}

			case mouse.Event:
				if e.Direction != mouse.DirRelease {
					break
				}

				// Re-map control-click to middle-click, etc, for computers with one-button mice.
				if e.Modifiers&key.ModControl != 0 {
					e.Button = mouse.ButtonMiddle
				} else if e.Modifiers&key.ModAlt != 0 {
					e.Button = mouse.ButtonRight
				} else if e.Modifiers&key.ModMeta != 0 {
					e.Button = mouse.ButtonMiddle
				}

				if board.click(b.RGBA(), int(e.X), int(e.Y), int(e.Button)) {
					w.Send(paint.Event{})
				}

			case paint.Event:
				w.Upload(image.Point{}, b, b.Bounds())
				w.Publish()

			case size.Event:
				// TODO: Set board size.
				if b != nil {
					b.Release()
				}
				b, err = s.NewBuffer(e.Size())
				if err != nil {
					log.Fatal(err)
				}
				render(b.RGBA(), board)

			case error:
				log.Print(e)
			}
		}
	})
}

func render(m *image.RGBA, board *Board) {
	board.Draw(m)
}
