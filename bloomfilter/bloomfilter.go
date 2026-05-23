package bloomfilter

import (
	"fmt"
	"math"

	"github.com/cespare/xxhash/v2"
)

type BloomFilter struct {
	expectedSize int
	rate         float64
	k            int
	bitArray     []byte
}

// expected size and false positive rate depends on the use case
func New(expectedSize int, falsePositiveRate float64) *BloomFilter {

	// m = -(n * ln(p)) / (ln(2) ^ 2) - size of bit array
	m := (-(float64(expectedSize) * math.Log(falsePositiveRate))) / math.Pow(math.Log(2), 2)
	// k = (n / m) * ln(2) - no. of hash functions
	k := (m / float64(expectedSize)) * math.Log(2)

	numBytes := int(math.Ceil(m / 8))

	return &BloomFilter{
		expectedSize: expectedSize,
		rate:         falsePositiveRate,
		k:            int(math.Round(k)),
		bitArray:     make([]byte, numBytes),
	}
}

func (bf *BloomFilter) Add(value string) {
	// hash value thru k hash functions
	for i := 0; i < bf.k; i++ {
		hashValue := xxhash.Sum64String(fmt.Sprintf("%d:%s", i, value))

		// to find out which position in the byte array to place it
		// use division - hash / 8
		pos := hashValue % uint64(len(bf.bitArray)*8)
		byteIndex := pos / 8

		// to find out which bit on the byte
		// use 1 << (hash % 8) - the remainder
		bitIndex := byte(1 << (pos % 8))

		// to flip the specific bit on the byte
		bf.bitArray[byteIndex] = bf.bitArray[byteIndex] | bitIndex
	}
}

func (bf *BloomFilter) Check(value string) bool {
	for i := 0; i < bf.k; i++ {
		hashValue := xxhash.Sum64String(fmt.Sprintf("%d:%s", i, value))

		pos := hashValue % uint64(len(bf.bitArray)*8)

		// to find out which position in the byte array to place it
		// use division - hash / 8
		byteIndex := pos / 8

		// to find out which bit on the byte
		// use 1 << (hash % 8) - the remainder
		bitIndex := byte(1 << (pos % 8))

		// check if the bit position is 1
		if bf.bitArray[byteIndex]&bitIndex == 0 {
			return false
		}
	}

	return true
}
