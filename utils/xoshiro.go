package utils

import (
	"math/big"
)

// MAX_UINT64 is the maximum value for a uint64
var MAX_UINT64 = new(big.Int).SetUint64(0xFFFFFFFFFFFFFFFF)

// rotl performs a left rotation on a big.Int value
func rotl(x *big.Int, k uint) *big.Int {
	left := new(big.Int).Lsh(x, k)
	right := new(big.Int).Rsh(x, 64-k)
	return new(big.Int).Xor(left, right)
}

// Xoshiro represents the Xoshiro random number generator
type Xoshiro struct {
	s [4]*big.Int
}

// NewXoshiro initializes a new Xoshiro instance with a seed
func NewXoshiro(seed []byte) *Xoshiro {
	digest := sha256Hash(seed)
	x := &Xoshiro{
		s: [4]*big.Int{big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
	}
	x.setS(digest)
	return x
}

// setS sets the internal state of the Xoshiro instance
func (x *Xoshiro) setS(digest []byte) {
	for i := 0; i < 4; i++ {
		o := i * 8
		v := big.NewInt(0)
		for n := 0; n < 8; n++ {
			v.Lsh(v, 8)
			v.Or(v, big.NewInt(int64(digest[o+n])))
		}
		x.s[i] = new(big.Int).SetUint64(v.Uint64())
	}
}

// roll generates the next value in the Xoshiro sequence
func (x *Xoshiro) roll() *big.Int {
	s1 := new(big.Int).Set(x.s[1])
	result := new(big.Int).Mul(rotl(new(big.Int).Mul(s1, big.NewInt(5)), 7), big.NewInt(9))
	result.And(result, MAX_UINT64)

	t := new(big.Int).Lsh(s1, 17)

	x.s[2].Xor(x.s[2], x.s[0])
	x.s[3].Xor(x.s[3], x.s[1])
	x.s[1].Xor(x.s[1], x.s[2])
	x.s[0].Xor(x.s[0], x.s[3])
	x.s[2].Xor(x.s[2], t)
	x.s[3] = rotl(x.s[3], 45)

	return result
}

// Next generates the next BigNumber in the Xoshiro sequence
func (x *Xoshiro) Next() *big.Int {
	return x.roll()
}

// NextDouble generates the next double value in the Xoshiro sequence
func (x *Xoshiro) NextDouble() *big.Float {
	max := new(big.Float).SetInt(MAX_UINT64)
	value := new(big.Float).SetInt(x.roll())
	return new(big.Float).Quo(value, new(big.Float).Add(max, big.NewFloat(1)))
}

// NextInt generates the next integer in the given range
func (x *Xoshiro) NextInt(low, high int) int {
	diff := high - low + 1
	value := new(big.Float).Mul(new(big.Float).SetInt(big.NewInt(int64(diff))), x.NextDouble())
	result, _ := value.Int64()
	return int(result) + low
}

// NextByte generates the next byte value in the Xoshiro sequence
func (x *Xoshiro) NextByte() byte {
	return byte(x.NextInt(0, 255))
}

// NextData generates the next sequence of bytes
func (x *Xoshiro) NextData(count int) []byte {
	data := make([]byte, count)
	for i := range data {
		data[i] = x.NextByte()
	}
	return data
}
