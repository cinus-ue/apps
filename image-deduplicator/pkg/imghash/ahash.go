package imghash

import (
	"../resize"
	"image"
)

// Ahash calculates a hash value for the given image. It implements the Average hash algorithm as explained here:
// http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
// Returns the hash of the image in 64 bits integer.
func Ahash(img *image.Image) uint64 {
	// resize the image
	resizedImage := resize.Resize(8, 8, *img, INTERP)

	// get the gray image
	grayImage := imageToGray(&resizedImage)

	// find the average of image
	average := getAverageGray(grayImage)

	// get the hash value of the image
	// hash bit is 1 if pixel value is above the average
	hash := uint64(0)
	width := grayImage.Bounds().Dx()
	height := grayImage.Bounds().Dy()
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			value := uint64(grayImage.GrayAt(i, j).Y)
			if value > average {
				hash += 1
			}
			hash = hash << 1
		}
	}
	return hash
}
