package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"os"

	"github.com/indeedplusplus/go-jpegxl"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s INPUT OUTPUT\n", os.Args[0])
		os.Exit(1)
	}
	inputName := os.Args[1]
	outputName := os.Args[2]
	img := func() image.Image {
		f, err := os.Open(inputName)
		if err != nil {
			log.Fatal(err)
		}
		defer func(f io.Closer) {
			_ = f.Close()
		}(f)
		img, err := jpegxl.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
		return img
	}()
	f, err := os.Create(outputName)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f io.Closer) {
		_ = f.Close()
	}(f)
	err = png.Encode(f, img)
	if err != nil {
		log.Fatal(err)
	}
}
