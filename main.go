package main

import (
	"fmt"
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

func main() {
	//flag.Parse()

	rand.Seed(int64(time.Now().Nanosecond()))
	if len(os.Args) > 1 {
		help := os.Args[1]
		if help != "" && help == "help" {
			fmt.Println("Usage: ./gogo size scale\nOne of ./gogo 9 35, or ./gogo 19 18, or ./gogo (defaults to 9 35)\nOptionally: ./gogo help (this help)")
			os.Exit(0)
		}
	}
	args := os.Args[1:]
	var dim int
	var err error
	if len(args) > 0 {
		dimStr := args[0]
		if dimStr == "" || (dimStr != "9" && dimStr != "19") {
			dimStr = "9"
		}
		dim, err = strconv.Atoi(dimStr)
		if err != nil {
			log.Fatal(err)
		}
		if len(args) == 2 {
			scaleStr := args[1]
			if scaleStr == "" || (scaleStr != "35" && scaleStr != "18") {
				if dim == 19 {
					scaleStr = "18"
				} else {
					scaleStr = "35"
				}
			}
			scaleInt, err := strconv.Atoi(scaleStr)
			if err != nil {
				log.Fatal(err)
			}
			scale = &scaleInt
		}
	} else {
		dim = 9
		scaleInt := 35
		scale = &scaleInt
	}

	fmt.Printf("Board: %dx%d\n", dim, dim)
	fmt.Printf("Scale: %d%%\n", *scale)

	board := NewBoard(dim, *scale)

	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Title: "Gogo",
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
