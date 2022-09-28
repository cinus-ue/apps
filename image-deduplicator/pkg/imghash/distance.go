package imghash

func HammingDistance(d1, d2 uint64) uint64 {
	distance := uint64(0)
	for bitPos := uint(0); bitPos < 64; bitPos++ {
		d1Bit := d1 & (1 << bitPos)
		d2Bit := d2 & (1 << bitPos)
		if d1Bit != d2Bit {
			distance += 1
		}
	}
	return distance
}
