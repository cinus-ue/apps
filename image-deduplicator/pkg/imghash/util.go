package imghash

import (
	"../resize"
	"fmt"
	"image"
	"image/png"
	"os"
)

var INTERP = resize.Bicubic

func getAverageGray(img *image.Gray) uint64 {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	sum := uint64(0)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			sum += uint64(img.GrayAt(i, j).Y)
		}
	}
	return sum / uint64(width*height)
}

func getAverageArray(arr [][]float64, ignoreFirst bool) uint64 {
	width := len(arr)
	height := len(arr[0])
	sum := uint64(0)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			if ignoreFirst && i == 0 && j == 0 {
				continue
			}
			sum += uint64(arr[i][j])
		}
	}
	if ignoreFirst {
		return sum / (uint64(width*height) - 1)
	}
	return sum / uint64(width*height)
}

func saveImage(img image.Image, filename string) {
	out, err := os.Create("./" + filename + ".png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = png.Encode(out, img)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
