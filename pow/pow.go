package pow

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math"
	"strconv"
)

const nonceSize = 8

func CalculateNonce(difficulty uint32, data []byte) ([]byte, error) {
	var currentNonce uint64

	dataOffset := len(data)

	buf := make([]byte, len(data)+nonceSize)
	copy(buf, data)
	for currentNonce < math.MaxUint64 {
		nonceBuf := make([]byte, nonceSize)
		binary.BigEndian.PutUint64(nonceBuf, currentNonce)

		copy(buf[dataOffset:], nonceBuf)

		hash := sha256.Sum256(buf)

		countZeroes := countLeadingZeroes(hash)
		if countZeroes == difficulty {
			return nonceBuf, nil
		} else {
			currentNonce++
		}
	}

	return nil, errors.New("cannot calculate nonce for given difficulty and data")
}

func IsValidNonce(data []byte, nonce []byte, difficulty uint32) error {
	if len(nonce) != nonceSize {
		return errors.New("incorrect bytes size in nonce: expected " + strconv.Itoa(nonceSize))
	}

	buf := make([]byte, len(data)+len(nonce))
	copy(buf, data)
	copy(buf[len(data):], nonce)

	hash := sha256.Sum256(buf)
	if countZeroes := countLeadingZeroes(hash); countZeroes != difficulty {
		return errors.New("incorrect count of digits for given hash: expected: " + strconv.Itoa(int(difficulty)))
	}

	return nil
}
