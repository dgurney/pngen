package main

import (
	"flag"
	"fmt"
	"github.com/satori/go.uuid"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

const version = "1.0.3"

var t = time.Now().UnixNano()
var rndseed = rand.NewSource(t)
var rnd = rand.New(rndseed)

func genImg(ch chan *image.NRGBA, width, height, mc int) {
	randomImg := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			randomImg.Set(x, y, color.NRGBA{
				R: uint8(rnd.Intn(mc)),
				G: uint8(rnd.Intn(mc)),
				B: uint8(rnd.Intn(mc)),
				A: uint8(255),
			})
		}
	}
	ch <- randomImg
}

func saveImg(img *image.NRGBA, width, height, amount int, wg *sync.WaitGroup) {
	defer wg.Done()
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	file, err := os.Create("random_" + strconv.Itoa(width) + "x" + strconv.Itoa(height) + "_" + id.String() + ".png")
	if err != nil {
		panic(err)
	}
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
	file.Close()
}

func main() {
	height := flag.Int("h", 500, "Height of image(s).")
	width := flag.Int("w", 500, "Width of image(s).")
	amount := flag.Int("a", 1, "Amount of images to generate.")
	maxcolors := flag.Int("r", 255, "The highest RGBA value that can be generated. Maximum is 255.")
	ver := flag.Bool("v", false, "Show version number and exit.")
	flag.Parse()
	if *maxcolors < 1 || *maxcolors > 255 {
		*maxcolors = 255
	}
	switch {
	case *height < 1 && *width < 1:
		*width = 500
		*height = 500
	case *width < 1:
		*width = 500
	case *height < 1:
		*height = 500
	}
	if *amount < 1 {
		*amount = 1
	}
	if *ver {
		fmt.Printf("Random PNG generator v%s by Daniel Gurney\n", version)
		return
	}
	switch {
	case *amount > 1:
		fmt.Printf("Generating %d %dx%d PNG files...\n", *amount, *width, *height)
	default:
		fmt.Printf("Generating a single %dx%d PNG file...\n", *width, *height)
	}
	ch := make(chan *image.NRGBA)
	var wg sync.WaitGroup
	// Concurrency reduces the time required by approximately 63%!
	for i := 0; i < *amount; i++ {
		wg.Add(1)
		go genImg(ch, *width, *height, *maxcolors)
		go saveImg(<-ch, *width, *height, *amount, &wg)
	}
	wg.Wait()
}
