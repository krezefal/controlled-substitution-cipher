package main

import (
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"io"
	"log"
	"os"
)

func getBytes(file io.Reader) (image.Image, []byte) {

	img, err := bmp.Decode(file)
	if err != nil {
		log.Printf(err.Error())
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var byteRow []byte

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			byteRow = append(byteRow, byte(r), byte(g), byte(b))
		}
	}

	return img, byteRow
}

func getFile(path string) (image.Image, []byte) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf(err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf(err.Error())
		}
	}(file)

	img, byteRow := getBytes(file)

	if err != nil {
		log.Printf(err.Error())
		os.Exit(1)
	}

	return img, byteRow
}

func writeImage(width, height int, byteRow []byte) *image.RGBA {
	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}
	newImg := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	for i := 0; i < len(byteRow); i += 3 {
		pixel := color.RGBA{R: byteRow[i], G: byteRow[i + 1], B: byteRow[i + 2]}
		newImg.Set(i % height, i / height, pixel)
	}

	return newImg
}

func saveFile(img image.Image, byteRow []byte, path string) error {

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	newImg := writeImage(width, height, byteRow)

	file, err := os.Create(path)
	if err != nil {
		log.Printf(err.Error())
	}
	err = file.Close()
	if err != nil {
		log.Printf(err.Error())
	}

	err = bmp.Encode(file, newImg)
	if err != nil {
		log.Printf(err.Error())
	}

	return nil
}

func remove(s []uint8, i int) []uint8 {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
