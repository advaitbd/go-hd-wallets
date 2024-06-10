package crypto

type ETHSignature struct {
	RequestID RequestIDType `cbor:"1,keyasint,omitempty"`
	Signature []byte        `cbor:"2,keyasint,omitempty"`
	Origin    string        `cbor:"3,keyasint,omitempty"`
}
