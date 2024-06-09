package utils

import (
	"encoding/json"
	"fmt"
	"math"
)

// FountainEncoderPart represents a part of the fountain encoder
type FountainEncoderPart struct {
	SeqNum        uint32
	SeqLength     uint32
	MessageLength uint32
	Checksum      uint32
	Fragment      []byte
}

// NewFountainEncoderPart creates a new FountainEncoderPart
func NewFountainEncoderPart(seqNum, seqLength, messageLength, checksum uint32, fragment []byte) *FountainEncoderPart {
	return &FountainEncoderPart{
		SeqNum:        seqNum,
		SeqLength:     seqLength,
		MessageLength: messageLength,
		Checksum:      checksum,
		Fragment:      fragment,
	}
}

// CBOR encodes the FountainEncoderPart into CBOR format
func (f *FountainEncoderPart) FountainEncoderPartCBOR() ([]byte, error) {
	data := []interface{}{f.SeqNum, f.SeqLength, f.MessageLength, f.Checksum, f.Fragment}
	return json.Marshal(data) // Use JSON for simplicity
}

// Description returns a description of the FountainEncoderPart
func (f *FountainEncoderPart) Description() string {
	return fmt.Sprintf("seqNum:%d, seqLen:%d, messageLen:%d, checksum:%d, data:%x", f.SeqNum, f.SeqLength, f.MessageLength, f.Checksum, f.Fragment)
}

// FountainEncoder represents the fountain encoder
type FountainEncoder struct {
	messageLength  int
	fragments      [][]byte
	fragmentLength int
	seqNum         uint32
	checksum       uint32
}

// NewFountainEncoder creates a new FountainEncoder
func NewFountainEncoder(message []byte, maxFragmentLength, firstSeqNum, minFragmentLength int) *FountainEncoder {
	fragmentLength := findNominalFragmentLength(len(message), minFragmentLength, maxFragmentLength)
	fragments := partitionMessage(message, fragmentLength)
	checksum := getCRC(message)
	return &FountainEncoder{
		messageLength:  len(message),
		fragments:      fragments,
		fragmentLength: fragmentLength,
		seqNum:         uint32(firstSeqNum),
		checksum:       checksum,
	}
}

// FragmentsLength returns the length of the fragments
func (f *FountainEncoder) FragmentsLength() int {
	return len(f.fragments)
}

// Fragments returns the fragments
func (f *FountainEncoder) Fragments() [][]byte {
	return f.fragments
}

// MessageLength returns the message length
func (f *FountainEncoder) MessageLength() int {
	return f.messageLength
}

// IsComplete checks if the encoding is complete
func (f *FountainEncoder) IsComplete() bool {
	return f.seqNum >= uint32(len(f.fragments))
}

// IsSinglePart checks if the message is a single part
func (f *FountainEncoder) IsSinglePart() bool {
	return len(f.fragments) == 1
}

// SeqLength returns the sequence length
func (f *FountainEncoder) SeqLength() int {
	return len(f.fragments)
}

// Mix mixes the fragments at the given indexes
func (f *FountainEncoder) Mix(indexes []int) []byte {
	result := make([]byte, f.fragmentLength)
	for _, index := range indexes {
		result = bufferXOR(f.fragments[index], result)
	}
	return result
}

// NextPart returns the next part of the encoder
func (f *FountainEncoder) NextPart() *FountainEncoderPart {
	f.seqNum++
	indexes := ChooseFragments(int(f.seqNum), len(f.fragments), int(f.checksum))
	mixed := f.Mix(indexes)
	return NewFountainEncoderPart(f.seqNum, uint32(len(f.fragments)), uint32(f.messageLength), f.checksum, mixed)
}

// findNominalFragmentLength finds the nominal fragment length
func findNominalFragmentLength(messageLength, minFragmentLength, maxFragmentLength int) int {
	if messageLength <= 0 || minFragmentLength <= 0 || maxFragmentLength < minFragmentLength {
		panic("invalid input")
	}
	maxFragmentCount := int(math.Ceil(float64(messageLength) / float64(minFragmentLength)))
	fragmentLength := 0
	for fragmentCount := 1; fragmentCount <= maxFragmentCount; fragmentCount++ {
		fragmentLength = int(math.Ceil(float64(messageLength) / float64(fragmentCount)))
		if fragmentLength <= maxFragmentLength {
			break
		}
	}
	return fragmentLength
}

// partitionMessage partitions the message into fragments
func partitionMessage(message []byte, fragmentLength int) [][]byte {
	var fragments [][]byte
	for len(message) > 0 {
		end := fragmentLength
		if len(message) < fragmentLength {
			end = len(message)
		}
		fragment := make([]byte, fragmentLength)
		copy(fragment, message[:end])
		fragments = append(fragments, fragment)
		message = message[end:]
	}
	return fragments
}
