package utils

import (
	"fmt"
	"strings"
)

type UREncoder struct {
	ur              *UR
	fountainEncoder *FountainEncoder
}

func NewUREncoder(ur *UR, maxFragmentLength, firstSeqNum, minFragmentLength int) *UREncoder {
	return &UREncoder{
		ur:              ur,
		fountainEncoder: NewFountainEncoder(ur.cborPayload, maxFragmentLength, firstSeqNum, minFragmentLength),
	}
}

func (u *UREncoder) FragmentsLength() int {
	return u.fountainEncoder.FragmentsLength()
}

func (u *UREncoder) Fragments() [][]byte {
	return u.fountainEncoder.Fragments()
}

func (u *UREncoder) MessageLength() int {
	return u.fountainEncoder.MessageLength()
}

func (u *UREncoder) Cbor() []byte {
	return u.ur.cborPayload
}

func (u *UREncoder) EncodeWhole() []string {
	parts := make([]string, u.FragmentsLength())
	for i := range parts {
		parts[i] = u.NextPart()
	}
	return parts
}

func (u *UREncoder) NextPart() string {
	part := u.fountainEncoder.NextPart()
	if u.fountainEncoder.IsSinglePart() {
		return EncodeSinglePart(u.ur)
	}
	return EncodePart(u.ur.typ, part)
}

func EncodeUri(scheme string, pathComponents []string) string {
	path := strings.Join(pathComponents, "/")
	return fmt.Sprintf("%s:%s", scheme, path)
}

func EncodeUR(pathComponents []string) string {
	return EncodeUri("ur", pathComponents)
}

func EncodePart(typeStr string, part *FountainEncoderPart) string {
	seq := fmt.Sprintf("%d-%d", part.SeqNum, part.SeqLength)
	data, _ := part.FountainEncoderPartCBOR()
	body, _ := BytewordEncode(fmt.Sprintf("%x", data), "uri")
	return EncodeUR([]string{typeStr, seq, body})
}

func EncodeSinglePart(ur *UR) string {
	body, _ := BytewordEncode(fmt.Sprintf("%x", ur.Cbor), "uri")
	return EncodeUR([]string{ur.typ, body})
}
