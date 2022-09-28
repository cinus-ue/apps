package imghash

import (
	"../go-fourier"
	"../resize"
	"image"
	"log"
)

// Phash calculates a hash value for the given image. It implements the pHash algorithm as explained here:
// http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
// Returns the hash of the image in 64 bits integer.
func Phash(img *image.Image) uint64 {
	// resize the image
	resizedImage := resize.Resize(32, 32, *img, INTERP)

	// get the gray image
	grayImage := imageToGray(&resizedImage)

	// setup 2d array for dct
	dctImage := make([][]float64, 32)
	for i := 0; i < 32; i++ {
		dctImage[i] = make([]float64, 32)
		for j := 0; j < 32; j++ {
			dctImage[i][j] = float64(grayImage.GrayAt(i, j).Y)
		}
	}

	// apply dct
	dctImage, err := go_fourier.DCT2D(dctImage)
	if err != nil {
		log.Fatal(err)
	}

	// reduce dct
	dctImage = dctImage[0:8]
	for i := 0; i < 8; i++ {
		dctImage[i] = dctImage[i][0:8]
	}

	// find the average of image
	average := getAverageArray(dctImage, true)

	// get the hash value of the image
	// hash bit is 1 if pixel value is above the average
	hash := uint64(0)

	// get the hash value of the image
	// hash bit is 1 if pixel value is above the average
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			value := uint64(dctImage[i][j])
			if value > average {
				hash += 1
			}
			hash = hash << 1
		}
	}

	return hash
}
