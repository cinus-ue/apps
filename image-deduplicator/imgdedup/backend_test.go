package imgdedup

import (
	"testing"
)

func TestListImagesInDir(t *testing.T) {
	t.Logf("TESTING NON RECURSIVE CALL\n")
	result, err := ListImagesInDir("../tests/random", false)
	if err != nil {
		t.Error("TestListImagesInDir: ", err)
	} else {
		t.Log("Number of files in the dir: ", len(result))
		for i, file := range result {
			t.Logf("File %d: %s\n", i+1, file.FileName)
		}
	}
	t.Logf("TESTING RECURSIVE CALL\n")
	result, err = ListImagesInDir("../tests/random", true)
	if err != nil {
		t.Error("TestListImagesInDir: ", err)
	} else {
		t.Log("Number of files in the dir: ", len(result))
		for i, file := range result {
			t.Logf("File %d: %s\n", i+1, file.FileName)
		}
	}
}

func BenchmarkHashImages(b *testing.B) {
	result, _ := ListImagesInDir("../tests/single", false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashImages(result, false)
	}
}

func BenchmarkHashImage(b *testing.B) {
	result, _ := ListImagesInDir("../tests/single", false)
	img := readImage(result[0])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hashImage(img)
	}
}

func BenchmarkReadImage(b *testing.B) {
	result, _ := ListImagesInDir("../tests/single", false)
	for i := 0; i < b.N; i++ {
		readImage(result[0])
	}
}
