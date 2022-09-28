package imghash

import (
	"../resize"
	"image"
)

// Dhash calculates a hash value for the given image. It implements the dHash (gradient) algorithm as explained here:
// http://www.hackerfactor.com/blog/index.php?/archives/529-Kind-of-Like-That.html
// Returns the hash of the image in 64 bits integer.
func Dhash(img *image.Image) uint64 {
	// resize the image
	resizedImage := resize.Resize(9, 8, *img, INTERP)

	// get the gray image
	grayImage := imageToGray(&resizedImage)

	// get the difference matrix
	diffMatrix := make([][]int, 8)
	for i := 0; i < 8; i++ {
		diffMatrix[i] = make([]int, 8)
		for j := 0; j < 8; j++ {
			diffMatrix[i][j] = int(grayImage.GrayAt(i, j+1).Y) - int(grayImage.GrayAt(i, j).Y)
		}
	}

	// get the hash value of the image
	// hash bit is 1 if pixel value is above the average
	hash := uint64(0)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			value := diffMatrix[i][j]
			if value > 0 {
				hash += 1
			}
			hash = hash << 1
		}
	}
	return hash
}
