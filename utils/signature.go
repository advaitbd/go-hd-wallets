package utils

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type Signature struct {
	R             string
	S             string
	V             uint8
	RecoveryParam int
	VS            string
	YParityAndS   string
	Compact       string
}

// HexCharacters are the hex digits
var HexCharacters = "0123456789abcdef"

// DataOptions represent the options for hexlify
type DataOptions struct {
	AllowMissingPrefix bool
	HexPad             string
}

func isHexString(s string) bool {
	_, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	return err == nil
}

func Hexlify(value interface{}, options *DataOptions) (string, error) {
	if options == nil {
		options = &DataOptions{}
	}

	switch v := value.(type) {
	case int:
		if v < 0 {
			return "", errors.New("invalid hexlify value")
		}

		hexStr := fmt.Sprintf("%x", v)
		if len(hexStr)%2 != 0 {
			hexStr = "0" + hexStr
		}
		return "0x" + hexStr, nil

	case int64:
		hexStr := fmt.Sprintf("%x", v)
		if len(hexStr)%2 != 0 {
			hexStr = "0" + hexStr
		}
		return "0x" + hexStr, nil

	case string:
		if options.AllowMissingPrefix && !strings.HasPrefix(v, "0x") {
			v = "0x" + v
		}

		if isHexString(v) {
			if len(v)%2 != 0 {
				switch options.HexPad {
				case "left":
					v = "0x0" + v[2:]
				case "right":
					v = v + "0"
				default:
					return "", errors.New("hex data is odd-length")
				}
			}
			return strings.ToLower(v), nil
		}

		return "", errors.New("invalid hexlify value")

	case []byte:
		result := "0x"
		for _, b := range v {
			result += string(HexCharacters[(b&0xf0)>>4]) + string(HexCharacters[b&0x0f])
		}
		return result, nil

	default:
		return "", errors.New("invalid hexlify value")
	}
}

func Concat(items [][]byte) ([]byte, error) {
	// Calculate the total length of the resulting slice
	totalLength := 0
	for _, item := range items {
		totalLength += len(item)
	}

	// Create a result slice of the total length
	result := make([]byte, totalLength)

	// Copy each item into the result slice
	offset := 0
	for _, item := range items {
		copy(result[offset:], item)
		offset += len(item)
	}

	return result, nil
}

func SplitSignature(signature []byte) (Signature, error) {
	var result Signature

	if len(signature) != 64 && len(signature) != 65 {
		return result, errors.New("invalid signature length")
	}

	result.R = "0x" + hex.EncodeToString(signature[:32])
	result.S = "0x" + hex.EncodeToString(signature[32:64])

	if len(signature) == 65 {
		result.V = uint8(signature[64])
	} else {
		result.V = 27 + (signature[32] >> 7)
		signature[32] &= 0x7F
	}

	if result.V < 27 {
		if result.V == 0 || result.V == 1 {
			result.V += 27
		} else {
			return result, errors.New("invalid v byte")
		}
	}

	result.RecoveryParam = 1 - int(result.V%2)

	if result.RecoveryParam != 0 {
		signature[32] |= 0x80
	}
	result.VS = "0x" + hex.EncodeToString(signature[32:64])

	result.YParityAndS = result.VS
	result.Compact = result.R + result.YParityAndS[2:]

	return result, nil
}

func JoinSignature(signature Signature) (string, error) {
	rBytes, err := hex.DecodeString(strings.TrimPrefix(signature.R, "0x"))
	if err != nil {
		return "", err
	}

	sBytes, err := hex.DecodeString(strings.TrimPrefix(signature.S, "0x"))
	if err != nil {
		return "", err
	}

	vByte := byte(0x1b)
	if signature.RecoveryParam == 1 {
		vByte = 0x1c
	}

	joined := append(rBytes, sBytes...)
	joined = append(joined, vByte)

	return "0x" + hex.EncodeToString(joined), nil
}
