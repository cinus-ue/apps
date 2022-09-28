package imghash

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

type ImageNamePair struct {
	img      image.Image
	filename string
}

var tests []ImageNamePair

func setup() {
	files, err := ioutil.ReadDir("./resources")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Name()[0] == '.' {
			continue
		}
		fmt.Printf("Reading file %v\n", file.Name())
		infile, err := os.Open("resources/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		src, _, err := image.Decode(infile)
		if err != nil {
			log.Fatal(err)
		}

		tests = append(tests, ImageNamePair{img: src, filename: file.Name()})

		err = infile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func TestAhash(t *testing.T) {
	for _, imageNamePair := range tests {
		hashValue := Ahash(&imageNamePair.img)
		t.Logf("Filename: %v\t Ahash: %016x", imageNamePair.filename, hashValue)
	}
}

func TestDhash(t *testing.T) {
	for _, imageNamePair := range tests {
		hashValue := Dhash(&imageNamePair.img)
		t.Logf("Filename: %v\t Dhash: %016x", imageNamePair.filename, hashValue)
	}
}

func TestPhash(t *testing.T) {
	for _, imageNamePair := range tests {
		hashValue := Phash(&imageNamePair.img)
		t.Logf("Filename: %v\t Phash: %016x", imageNamePair.filename, hashValue)
	}
}

func TestImageToGraySliced(t *testing.T) {
	for _, imageNamePair := range tests {
		t.Logf("Working on file: %v\n", imageNamePair.filename)
		grayImage := imageToGray(&imageNamePair.img)
		grayImageSliced := imageToGraySliced(&imageNamePair.img)
		bounds := grayImage.Rect
		found := false
		var fx, fy int
		for i := 0; i < bounds.Max.X; i++ {
			for j := 0; j < bounds.Max.Y; j++ {
				if grayImage.GrayAt(i, j) != grayImageSliced.GrayAt(i, j) {
					found = true
					fx = i
					fy = j
					break
				}
			}
			if found {
				break
			}
		}
		if found {
			t.Errorf("Found error at (%v, %v), %v, %v\n", fx, fy, grayImage.GrayAt(fx, fy), grayImageSliced.GrayAt(fx, fy))
		}
	}
}

func BenchmarkResize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resize.Resize(8, 8, tests[0].img, INTERP)
	}
}

func BenchmarkAhash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Ahash(&tests[0].img)
	}
}

func BenchmarkDhash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Dhash(&tests[0].img)
	}
}

func BenchmarkPhash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Phash(&tests[0].img)
	}
}

func BenchmarkImageToGray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		imageToGray(&tests[0].img)
	}
}

func BenchmarkImageToGraySliced(b *testing.B) {
	for i := 0; i < b.N; i++ {
		imageToGraySliced(&tests[0].img)
	}
}
