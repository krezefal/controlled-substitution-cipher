package main

import (
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"io"
	"log"
	"os"
)

func to4Bytes(r, g, b, a uint32) []byte {
	return []byte{byte(r / 257), byte(g / 257), byte(b / 257), byte(a / 257)}
}

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
			byteRow = append(byteRow, to4Bytes(img.At(x, y).RGBA())...)
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

	var kk int
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			cyan := color.RGBA{R: byteRow[kk], G: byteRow[kk+1], B: byteRow[kk+2], A: byteRow[kk+3]}
			newImg.Set(x, y, cyan)
			kk+=4
		}
	}

	return newImg
}

func saveFile(img image.Image, byteRow []byte, path string) {

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	newImg := writeImage(width, height, byteRow)

	file, err := os.Create(path)
	if err != nil {
		log.Printf(err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf(err.Error())
		}
	}(file)

	err = bmp.Encode(file, newImg)
	if err != nil {
		log.Printf(err.Error())
	}
}

func remove(s []uint8, i int) []uint8 {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
