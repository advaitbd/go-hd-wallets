package utils

import (
	"bytes"
	"encoding/binary"
	"math/big"
)

// IntToBytes converts an int to a byte slice
func IntToBytes(num int) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, int32(num))
	return buf.Bytes()
}

// ChooseDegree chooses a degree based on the sequence length and RNG
func ChooseDegree(seqLength int, rng *Xoshiro) int {
	degreeProbabilities := make([]float64, seqLength)
	for i := 0; i < seqLength; i++ {
		degreeProbabilities[i] = 1.0 / float64(i+1)
	}

	degreeChooser := NewAliasSampler(degreeProbabilities, rng)

	return degreeChooser.Next() + 1
}

// Shuffle shuffles the items using the given RNG
func Shuffle(items []interface{}, rng *Xoshiro) []interface{} {
	remaining := make([]interface{}, len(items))
	copy(remaining, items)
	var result []interface{}

	for len(remaining) > 0 {
		index := rng.NextInt(0, len(remaining)-1)
		item := remaining[index]
		remaining = append(remaining[:index], remaining[index+1:]...)
		result = append(result, item)
	}

	return result
}

// ChooseFragments chooses the fragments based on the sequence number, length, and checksum
func ChooseFragments(seqNum int, seqLength int, checksum int) []int {
	if seqNum <= seqLength {
		return []int{seqNum - 1}
	} else {
		seed := append(IntToBytes(seqNum), IntToBytes(checksum)...)
		rng := NewXoshiro(seed)
		degree := ChooseDegree(seqLength, rng)
		indexes := make([]int, seqLength)
		for i := 0; i < seqLength; i++ {
			indexes[i] = i
		}
		shuffledIndexes := ShuffleInt(indexes, rng)
		return shuffledIndexes[:degree]
	}
}

// AliasSampler is a simple alias sampling method
type AliasSampler struct {
	probabilities []float64
	aliases       []int
	rng           *Xoshiro
}

// NewAliasSampler creates a new alias sampler
func NewAliasSampler(probabilities []float64, rng *Xoshiro) *AliasSampler {
	n := len(probabilities)
	aliases := make([]int, n)
	aliasSampler := &AliasSampler{
		probabilities: probabilities,
		aliases:       aliases,
		rng:           rng,
	}

	// Here is the alias sampling initialization logic
	small := []int{}
	large := []int{}
	scaledProbabilities := make([]float64, n)
	for i := 0; i < n; i++ {
		scaledProbabilities[i] = probabilities[i] * float64(n)
		if scaledProbabilities[i] < 1.0 {
			small = append(small, i)
		} else {
			large = append(large, i)
		}
	}

	for len(small) > 0 && len(large) > 0 {
		smallIndex := small[len(small)-1]
		small = small[:len(small)-1]
		largeIndex := large[len(large)-1]
		large = large[:len(large)-1]

		aliasSampler.probabilities[smallIndex] = scaledProbabilities[smallIndex]
		aliasSampler.aliases[smallIndex] = largeIndex

		scaledProbabilities[largeIndex] += scaledProbabilities[smallIndex] - 1.0
		if scaledProbabilities[largeIndex] < 1.0 {
			small = append(small, largeIndex)
		} else {
			large = append(large, largeIndex)
		}
	}

	for len(large) > 0 {
		largeIndex := large[len(large)-1]
		large = large[:len(large)-1]
		aliasSampler.probabilities[largeIndex] = 1.0
	}

	for len(small) > 0 {
		smallIndex := small[len(small)-1]
		small = small[:len(small)-1]
		aliasSampler.probabilities[smallIndex] = 1.0
	}

	return aliasSampler
}

// Next returns the next sampled index
func (a *AliasSampler) Next() int {
	n := len(a.probabilities)
	i := a.rng.NextInt(0, n-1)
	p := a.rng.NextDouble()
	if p.Cmp(big.NewFloat(a.probabilities[i])) < 0 {
		return i
	}
	return a.aliases[i]
}

// ShuffleInt shuffles an int slice using the given RNG
func ShuffleInt(items []int, rng *Xoshiro) []int {
	remaining := make([]int, len(items))
	copy(remaining, items)
	var result []int

	for len(remaining) > 0 {
		index := rng.NextInt(0, len(remaining)-1)
		item := remaining[index]
		remaining = append(remaining[:index], remaining[index+1:]...)
		result = append(result, item)
	}

	return result
}
