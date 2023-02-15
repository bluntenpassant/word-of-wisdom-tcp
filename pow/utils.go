package pow

func countLeadingZeroes(data [32]byte) uint32 {
	var count uint32

	for _, b := range data {
		if b != 0 {
			break
		}

		count++
	}

	return count
}
