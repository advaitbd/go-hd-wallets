package utils

import (
	"bytes"
	"errors"

	"github.com/fxamacker/cbor/v2"
)

type UR struct {
	cborPayload []byte
	typ         string
}

func NewUR(cborPayload []byte, typ string) (*UR, error) {
	if !isURType(typ) {
		return nil, errors.New("invalid UR type")
	}
	return &UR{cborPayload: cborPayload, typ: typ}, nil
}

func URFromBuffer(buf []byte) (*UR, error) {
	encoded, err := cbor.Marshal(buf)
	if err != nil {
		return nil, err
	}
	return NewUR(encoded, "bytes")
}

func URFrom(value []byte, encoding string) (*UR, error) {
	return URFromBuffer(value)
}

func (ur *UR) DecodeCBOR() ([]byte, error) {
	var decoded []byte
	err := cbor.Unmarshal(ur.cborPayload, &decoded)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func (ur *UR) Type() string {
	return ur.typ
}

func (ur *UR) Cbor() []byte {
	return ur.cborPayload
}

func (ur *UR) Equals(ur2 *UR) bool {
	return ur.Type() == ur2.Type() && bytes.Equal(ur.Cbor(), ur2.Cbor())
}
