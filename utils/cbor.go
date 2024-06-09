package utils

import (
	"fmt"
	"reflect"

	"github.com/fxamacker/cbor/v2"
)

func FromCBOR(cborBytes []byte, v interface{}, tags cbor.TagSet) error {
	// Ensure v is a pointer to a struct
	if reflect.ValueOf(v).Kind() != reflect.Ptr || reflect.Indirect(reflect.ValueOf(v)).Kind() != reflect.Struct {
		return fmt.Errorf("v must be a pointer to a struct")
	}

	var err error
	if tags != nil {
		dm, err := cbor.DecOptions{}.DecModeWithTags(tags)
		if err != nil {
			return err
		}
		dm.Unmarshal(cborBytes, v)
	} else {
		err = cbor.Unmarshal(cborBytes, v)
	}
	if err != nil {
		return err
	}

	return nil
}

func ToCBOR(v interface{}, tags cbor.TagSet) ([]byte, error) {
	// Ensure v is a struct or a pointer to a struct
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("v must be a struct or a pointer to a struct")
	}

	var cborBytes []byte

	if tags != nil {
		em, _ := cbor.EncOptions{}.EncModeWithTags(tags)
		cborBytes, _ = em.Marshal(v)
	} else {
		cborBytes, _ = cbor.Marshal(v)
	}
	return cborBytes, nil
}
