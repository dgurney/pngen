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
)

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
func saveImg(img *image.NRGBA, width, height int, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Create("random_" + strconv.Itoa(width) + "x" + strconv.Itoa(height) + "_" + strconv.Itoa(rnd.Intn(9223372036854775807)) + ".png")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}
func main() {
	height := flag.Int("h", 500, "Height of image")
	width := flag.Int("w", 500, "Width of image")
	amount := flag.Int("a", 1, "Amount of images to generate")
	flag.Parse()
	switch {
	case *amount > 1:
		fmt.Printf("Generating %d PNG files...\n", *amount)
	default:
		fmt.Println("Generating a single PNG file...")
	}
	ch := make(chan *image.NRGBA)
	var wg sync.WaitGroup
	// Concurrency decreases the time required by approximately 63%!
	for i := 0; i < *amount; i++ {
		wg.Add(1)
		go genImg(ch, *width, *height)
		go saveImg(<-ch, *width, *height, &wg)
	}
	wg.Wait()
}
