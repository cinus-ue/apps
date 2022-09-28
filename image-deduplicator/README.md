# image-deduplicator
A program implemented in Go that groups similar images in a given directory or finds similar images in a directory to given image.

## Running the program
In the directory of the project build the executable:
```
 go build main.go
```
Then it can be run in the same directory:
```
./main [flags]
```
Here are the all the flags (also can be seen by running `./main -h`)
```
Usage of ./main:
  -algo int
        Indicates the algorithm that is to be used for hashing.
        1-aHash, 2-pHash, 3-dHash (default 1)
  -dir string
        Indicates the directory the program will run on. (default "./")
  -r    Indicates whether the program should look for images in sub directories.
  -t uint
        Indicates the threshold value for the grouping, smaller threshold value yield to higher confidence in results (default 5)
  -v    Indicates whether the program should print out logs.

```
## Dependencies
* [go-fourier](https://github.com/ardabasaran/go-fourier) package.

    ``` go get github.com/ardabasaran/go-fourier```
* [imghash](https://github.com/ardabasaran/imghash) package.

    ``` go get github.com/ardabasaran/imghash```