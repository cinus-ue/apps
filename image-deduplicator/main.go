package main

import (
	"./imgdedup"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Define and parse flags
	recFlagPtr := flag.Bool("r", false, "Indicates whether the program should look for images in sub directories.")
	algoFlagPtr := flag.Int("algo", 3, "Indicates the algorithm that is to be used for hashing.\n1-aHash, 2-pHash, 3-dHash")
	verboseFlagPtr := flag.Bool("v", false, "Indicates whether the program should print out logs.")
	directoryFlagPtr := flag.String("dir", "./", "Indicates the directory the program will run on.")
	thresholdFlagPtr := flag.Uint64("t", 15, "Indicates the threshold value for the grouping, smaller threshold value yield to higher confidence in results")
	moveFlagPtr := flag.Bool("m", false, "Indicates whether the program should move the similar images in the same group folder.")
	flag.Parse()

	// Assign flag values to corresponding variables
	imgdedup.HASHALGORITHM = *algoFlagPtr
	verbose := *verboseFlagPtr
	recursive := *recFlagPtr
	dir := *directoryFlagPtr
	threshold := *thresholdFlagPtr
	move := *moveFlagPtr

	// Create ImageHash structs for files in the directory
	images, err := imgdedup.ListImagesInDir(dir, recursive)
	if err != nil {
		log.Fatal(err)
	}

	// (if verbose print out number of files and file names)
	if verbose {
		fmt.Printf("Found %d images in total.\n", len(images))
		for i, image := range images {
			fmt.Printf("\tFile %d: %s\n", i, image.FileName)
		}
	}

	//Hash images
	// (if verbose, print out progress)
	imgdedup.HashImages(images, verbose)

	// Group images with threshold
	// (if verbose, print out groups)
	groups := imgdedup.GroupImages(images, threshold, verbose)

	if verbose {
		fmt.Printf("There are %d groups in total.\n", len(groups))
		singleCount := 0
		for _, group := range groups {
			if group.Size == 1 {
				singleCount++
				continue
			}
			fmt.Printf("Group ID: %d \t Group Size: %d\n", group.ID, group.Size)
			for _, imageWithHash := range group.Images {
				fmt.Printf("\tFile name: %s\n", imageWithHash.FileName)
			}
		}
		fmt.Printf("Number of singles: %d\n", singleCount)
	}

	if move {
		//Move images in groups with size >1 to directories
		for i, group := range groups {
			for _, imageWithHash := range group.Images {
				oldFilePath := (*imageWithHash).FilePath + "/" + (*imageWithHash).FileName
				newDirPath := filepath.Join(dir, fmt.Sprintf("group%d", i+1))
				err := os.MkdirAll(newDirPath, os.ModePerm)
				if err != nil {
					log.Fatal(err)
				}
				err = os.Rename(oldFilePath, newDirPath+"/"+(*imageWithHash).FileName)
			}
		}
	}
}
