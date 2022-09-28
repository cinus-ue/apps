package imgdedup

import (
	"../pkg/imghash"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
)

var HASHALGORITHM = 1

type ImageGroup struct {
	ID     uint64
	Size   uint64
	Images []*ImageWithHash
}

type ImageWithHash struct {
	FileName  string
	FilePath  string
	GroupID   uint64
	HashValue uint64
}

// Given a directory this function returns the list of image files in the directory (and in the subdirectories
// if the recursive parameter is set to true) as ImageWithHash type. GroupID and HashValues are set to 0 for
// these images.
func ListImagesInDir(dir string, recursive bool) ([]ImageWithHash, error) {
	// read the directory contents
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return make([]ImageWithHash, 0), err
	}

	// loop over the files in the directory
	result := make([]ImageWithHash, 0)
	for _, file := range files {
		// if it's a directory and subdirectories are included, do a recursive call
		if file.IsDir() && recursive {
			recursiveResult, recursiveErr := ListImagesInDir(dir+"/"+file.Name(), recursive)
			if recursiveErr != nil {
				return make([]ImageWithHash, 0), recursiveErr
			}
			result = append(result, recursiveResult...)
		}
		if file.IsDir() {
			continue
		}
		// if the file is an image, add it to the results
		if checkFileType(file.Name(), dir) {
			result = append(result, ImageWithHash{FileName: file.Name(), FilePath: dir, GroupID: 0, HashValue: 0})
		}
	}

	// sort the results with respect to file names
	sort.Slice(result, func(i, j int) bool {
		return result[i].FileName < result[j].FileName
	})
	return result, nil
}

// Given a list of images this function groups them according to their hash value with respect to given threshold
// Returns a list ImageGroups.
func GroupImages(images []ImageWithHash, threshold uint64, verbose bool) []ImageGroup {
	if verbose {
		fmt.Println("Started grouping images")
	}
	groups := make([]ImageGroup, 0)
	idCount := uint64(1)
	for i := 0; i < len(images); i++ {
		if images[i].GroupID != 0 {
			continue
		}
		imagesInTheGroup := make([]*ImageWithHash, 0)
		imagesInTheGroup = append(imagesInTheGroup, &images[i])
		images[i].GroupID = idCount
		size := uint64(1)
		for j := i + 1; j < len(images); j++ {
			if images[j].GroupID != 0 {
				continue
			}
			if imghash.HammingDistance(images[i].HashValue, images[j].HashValue) < threshold {
				images[j].GroupID = idCount
				imagesInTheGroup = append(imagesInTheGroup, &images[j])
				size++
			}
		}
		groups = append(groups, ImageGroup{ID: idCount, Size: size, Images: imagesInTheGroup})
		if verbose {
			fmt.Printf("\tGroup %d has size %d\n", groups[idCount-1].ID, groups[idCount-1].Size)
		}
		idCount++
	}
	if verbose {
		fmt.Println("Finished grouping images")
	}
	return groups
}

// Given a list of images, this function calculates and sets the hash value of the structs passed.
// Prints out the progress if verbose is true
func HashImages(images []ImageWithHash, verbose bool) {
	if verbose {
		fmt.Println("Started hashing images")
	}
	for i := range images {
		img := readImage(images[i])
		images[i].HashValue = hashImage(img)
		if verbose {
			fmt.Printf("\t%d/%d Image: %s\tHash: %064b \n", i+1, len(images), images[i].FileName, images[i].HashValue)
		}
	}
	if verbose {
		fmt.Println("Finished hashing images")
	}
}

func readImage(imageHash ImageWithHash) *image.Image {
	infile, _ := os.Open(imageHash.FilePath + "/" + imageHash.FileName)
	img, _, _ := image.Decode(infile)
	return &img
}

func hashImage(img *image.Image) uint64 {
	switch HASHALGORITHM {
	case 1:
		return imghash.Ahash(img)
	case 2:
		return imghash.Phash(img)
	default:
		return imghash.Dhash(img)
	}
}

func checkFileType(name, dir string) bool {
	// open the files
	file, err := os.Open(dir + "/" + name)
	if err != nil {
		return false
	}

	// read 512 bits of buffer
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return false
	}

	// detect the type
	t := http.DetectContentType(buffer)

	// return true if it's an image type
	if t == "image/jpeg" || t == "image/png" {
		return true
	}

	return false
}
