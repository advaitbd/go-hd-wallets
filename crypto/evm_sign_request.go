package crypto

type RequestIDType []byte

type EVMSignRequest struct {
	RequestID      RequestIDType `cbor:"1,keyasint,omitempty"`
	SignData       []byte        `cbor:"2,keyasint,omitempty"`
	DataType       DataType      `cbor:"3,keyasint,omitempty"`
	ChainID        uint64        `cbor:"4,keyasint,omitempty"`
	DerivationPath KeyPath       `cbor:"5,keyasint,omitempty"`
	Address        []byte        `cbor:"6,keyasint,omitempty"`
	Origin         string        `cbor:"7,keyasint,omitempty"`
}
