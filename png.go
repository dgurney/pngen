package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/satori/go.uuid"
)

const version = "1.0.2"

var t = time.Now().UnixNano()
var rndseed = rand.NewSource(t)
var rnd = rand.New(rndseed)

func genImg(ch chan *image.NRGBA, width, height int) {
	randomImg := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			randomImg.Set(x, y, color.NRGBA{
				R: uint8(rnd.Intn(255)),
				G: uint8(rnd.Intn(255)),
				B: uint8(rnd.Intn(255)),
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
	height := flag.Int("h", 500, "Height of image(s)")
	width := flag.Int("w", 500, "Width of image(s)")
	amount := flag.Int("a", 1, "Amount of images to generate.")
	ver := flag.Bool("v", false, "Show version number and exit.")
	flag.Parse()
	switch {
	case *height == 0 && *width == 0:
		*width = 500
		*height = 500
	case *width == 0:
		*width = 500
	case *height == 0:
		*height = 500
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
	// Concurrency decreases the time required by approximately 63%!
	for i := 0; i < *amount; i++ {
		wg.Add(1)
		go genImg(ch, *width, *height)
		go saveImg(<-ch, *width, *height, *amount, &wg)
	}
	wg.Wait()
}
