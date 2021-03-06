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

	"github.com/google/uuid"
)

const version = "1.0.7"

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func genImg(ch chan *image.NRGBA, width, height, mc int, allowalpha bool) {
	randomImg := image.NewNRGBA(image.Rect(0, 0, width, height))
	var alpha uint8 = 255
	if allowalpha {
		alpha = uint8(rnd.Intn(255))
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			randomImg.Set(x, y, color.NRGBA{
				R: uint8(rnd.Intn(mc)),
				G: uint8(rnd.Intn(mc)),
				B: uint8(rnd.Intn(mc)),
				A: alpha,
			})
		}
	}
	ch <- randomImg
}

func saveImg(img *image.NRGBA, width, height, amount int, wg *sync.WaitGroup) {
	defer wg.Done()
	id := uuid.New()

	file, err := os.Create("random_" + strconv.Itoa(width) + "x" + strconv.Itoa(height) + "_" + id.String() + ".png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	enc := &png.Encoder{
		CompressionLevel: png.NoCompression,
	}
	if err := enc.Encode(file, img); err != nil {
		panic(err)
	}
}

func main() {
	height := flag.Int("h", 500, "Height of image(s).")
	width := flag.Int("w", 500, "Width of image(s).")
	amount := flag.Int("a", 1, "Amount of images to generate.")
	maxcolors := flag.Int("r", 255, "The highest RGBA value that can be generated. Maximum is 255.")
	allowalpha := flag.Bool("l", false, "Randomize the alpha value in addition to RGB values. The upper limit is 255, as usual.")
	ver := flag.Bool("v", false, "Show version number and exit.")
	flag.Parse()
	if *ver {
		fmt.Printf("Random PNG generator v%s by Daniel Gurney\nArguments:\n", version)
		flag.PrintDefaults()
		return
	}
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
	switch {
	case *amount > 1:
		fmt.Printf("Generating %d %dx%d PNG files...\n", *amount, *width, *height)
	default:
		fmt.Printf("Generating a single %dx%d PNG file...\n", *width, *height)
	}
	ch := make(chan *image.NRGBA)
	var wg sync.WaitGroup
	for i := 0; i < *amount; i++ {
		wg.Add(1)
		go genImg(ch, *width, *height, *maxcolors, *allowalpha)
		go saveImg(<-ch, *width, *height, *amount, &wg)
	}
	wg.Wait()
}
