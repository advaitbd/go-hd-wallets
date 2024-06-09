package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"reflect"
	"strings"
	"unicode"
)

func sha256Hash(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func partition(s string, n int) []string {
	if n <= 0 {
		return []string{s}
	}

	var result []string
	for i := 0; i < len(s); i += n {
		end := i + n
		if end > len(s) {
			end = len(s)
		}
		result = append(result, s[i:end])
	}
	return result
}

func split(s []byte, length int) ([]byte, []byte) {
	if length <= 0 || length > len(s) {
		return s, nil
	}
	return s[:len(s)-length], s[len(s)-length:]
}

func getCRC(message []byte) uint32 {
	return crc32.ChecksumIEEE(message)
}

func getCRCHex(message []byte) string {
	crc := crc32.ChecksumIEEE(message)
	return fmt.Sprintf("%08x", crc)
}

func toUint32(number int) uint32 {
	return uint32(number)
}

func intToBytes(num uint32) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, num)
	return buf.Bytes()
}

func isURType(s string) bool {
	for _, c := range s {
		if unicode.IsLower(c) || unicode.IsDigit(c) || c == '-' {
			continue
		}
		return false
	}
	return true
}

func hasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func arraysEqual(ar1, ar2 []interface{}) bool {
	if len(ar1) != len(ar2) {
		return false
	}

	for _, el := range ar1 {
		found := false
		for _, el2 := range ar2 {
			if reflect.DeepEqual(el, el2) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func arrayContains(ar1, ar2 []interface{}) bool {
	for _, v := range ar2 {
		found := false
		for _, el := range ar1 {
			if reflect.DeepEqual(v, el) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func setDifference(ar1, ar2 []interface{}) []interface{} {
	var diff []interface{}
	for _, x := range ar1 {
		found := false
		for _, y := range ar2 {
			if reflect.DeepEqual(x, y) {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func bufferXOR(a, b []byte) []byte {
	length := len(a)
	if len(b) > length {
		length = len(b)
	}
	buffer := make([]byte, length)

	for i := 0; i < length; i++ {
		if i < len(a) && i < len(b) {
			buffer[i] = a[i] ^ b[i]
		} else if i < len(a) {
			buffer[i] = a[i]
		} else {
			buffer[i] = b[i]
		}
	}
	return buffer
}
