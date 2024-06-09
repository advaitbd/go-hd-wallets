package crypto

type DecodeURRequest struct {
	URData string `json:"urData"`
}

type HDKey struct {
	IsMaster          *bool     `cbor:"1,keyasint,omitempty"`
	IsPrivate         *bool     `cbor:"2,keyasint,omitempty"`
	KeyData           []byte    `cbor:"3,keyasint,omitempty"`
	ChainCode         []byte    `cbor:"4,keyasint,omitempty"`
	UseInfo           *CoinInfo `cbor:"5,keyasint,omitempty"`
	Origin            *KeyPath  `cbor:"6,keyasint,omitempty"`
	Children          *KeyPath  `cbor:"7,keyasint,omitempty"`
	ParentFingerprint *uint32   `cbor:"8,keyasint,omitempty"`
	Name              *string   `cbor:"9,keyasint,omitempty"`
	Note              *string   `cbor:"10,keyasint,omitempty"`
}

type KeyPath struct {
	Components        []interface{} `cbor:"1,keyasint,toarray"`
	SourceFingerprint uint32        `cbor:"2,keyasint,omitempty"`
	Depth             uint32        `cbor:"3,keyasint,omitempty"`
}

type PathComponent struct {
	_          struct{} `cbor:",toarray"`
	ChildIndex uint32
	IsHardened bool
}

type CoinInfo struct {
	Type    uint32 `cbor:"1,keyasint,omitempty"`
	Network int    `cbor:"2,keyasint,omitempty"`
}
