package imghash

import (
	"image"
	"image/color"
	"runtime"
	"sync"
)

type ImageSlice struct {
	xStart int
	xEnd   int
	yStart int
	yEnd   int
}

func imageToGray(img *image.Image) *image.Gray {
	bounds := (*img).Bounds()
	gray := image.NewGray(bounds)

	for i := 0; i < bounds.Max.X; i++ {
		for j := 0; j < bounds.Max.Y; j++ {
			rgbaPixel := (*img).At(i, j)
			red, green, blue, _ := rgbaPixel.RGBA()

			// See: https://en.wikipedia.org/wiki/Grayscale#Luma_coding_in_video_systems
			grayValue := 0.299*float64(red) + 0.587*float64(green) + 0.114*float64(blue)
			grayPixel := color.Gray{Y: uint8(grayValue)}

			gray.Set(i, j, grayPixel)
		}
	}
	return gray
}

func imageToGraySliced(img *image.Image) *image.Gray {
	bounds := (*img).Bounds()
	gray := image.NewGray(bounds)
	maxThreadNum := runtime.GOMAXPROCS(0)
	slices := sliceImage(bounds, maxThreadNum)
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(maxThreadNum)

	for n := 0; n < maxThreadNum; n++ {
		imageSlice := slices[n]
		go func(xStart, xEnd, yStart, yEnd int) {
			defer waitGroup.Done()
			for i := xStart; i < xEnd; i++ {
				for j := yStart; j < yEnd; j++ {
					rgbaPixel := (*img).At(i, j)
					red, green, blue, _ := rgbaPixel.RGBA()

					//// See: https://en.wikipedia.org/wiki/Grayscale#Luma_coding_in_video_systems
					grayValue := 0.299*float64(red) + 0.587*float64(green) + 0.114*float64(blue)
					grayPixel := color.Gray{Y: uint8(grayValue)}

					gray.Set(i, j, grayPixel)
				}
			}
		}(imageSlice.xStart, imageSlice.xEnd, imageSlice.yStart, imageSlice.yEnd)
	}

	waitGroup.Wait()
	return gray
}

func sliceImage(bounds image.Rectangle, n int) []ImageSlice {
	result := make([]ImageSlice, n)
	var columns bool
	if bounds.Max.X > bounds.Max.Y {
		columns = true
	} else {
		columns = false
	}
	for i := 0; i < n; i++ {
		var xStart, xEnd, yStart, yEnd int
		if columns {
			xStart = i * (bounds.Max.X / n)
			if i == n-1 {
				xEnd = bounds.Max.X
			} else {
				xEnd = (i + 1) * (bounds.Max.X / n)
			}
			yStart = 0
			yEnd = bounds.Max.Y
		} else {
			xStart = 0
			xEnd = bounds.Max.X
			yStart = i * (bounds.Max.Y / n)
			if i == n-1 {
				yEnd = bounds.Max.Y
			} else {
				yEnd = (i + 1) * (bounds.Max.Y / n)
			}
		}
		result[i] = ImageSlice{xStart: xStart, xEnd: xEnd, yStart: yStart, yEnd: yEnd}
		//fmt.Printf("Bounds: %v %v \t xStart:%v xEnd:%v yStart:%v yEnd:%v\n", bounds.Max.X, bounds.Max.Y, xStart, xEnd, yStart, yEnd)
	}

	return result
}
