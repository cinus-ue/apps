package imghash

import "image"

type ImageHashPair struct {
	Filename string
	Image    image.Image
	Hash     uint64
}
