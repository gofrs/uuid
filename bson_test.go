package uuid

import (
	"fmt"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

var uuidBSONSignature = []byte{16, 0, 0, 0, 4}

func TestUUIDMarshalUnmarshalBSON(t *testing.T) {

	testsIds := []struct {
		version int
		id      UUID
	}{
		{
			id:      Must(NewV1()),
			version: 1,
		},
		{
			id:      Must(NewV4()),
			version: 4,
		},
		{
			id:      Must(NewV6()),
			version: 6,
		},
	}
	for _, tID := range testsIds {
		t.Run(fmt.Sprintf("MarshalBSONValue UUID Version %d", tID.version), func(t *testing.T) {
			tType, tBytes, err := tID.id.MarshalBSONValue()
			if err != nil {
				t.Errorf("Error in MarshalBSONValue: %v", err)
			}
			if tType != bson.TypeBinary {
				t.Errorf("Expected bsontype.TypeBinary, got %v", tType)
			}
			expectedBytes := tID.id.Bytes()
			typePadBytes := tBytes[0:5]
			if !reflect.DeepEqual(typePadBytes, uuidBSONSignature) {
				t.Errorf("Expected %v, got %v", uuidBSONSignature, typePadBytes)
			}
			realDataBytes := tBytes[5:]
			if !reflect.DeepEqual(realDataBytes, expectedBytes) {
				t.Errorf("Expected %v, got %v", expectedBytes, realDataBytes)
			}
		})

		t.Run(fmt.Sprintf("UnmarshalBSONValue UUID Version %d", tID.version), func(t *testing.T) {
			raw := tID.id.Bytes()
			raw = append(uuidBSONSignature, raw...)
			u := UUID{}
			err := u.UnmarshalBSONValue(bson.TypeBinary, raw)
			if err != nil {
				t.Errorf("Error in UnmarshalBSONValue: %v", err)
			}
			if !reflect.DeepEqual(u, tID.id) {
				t.Errorf("Expected %v, got %v", tID.id, u)
			}
			if u.Version() != tID.id.Version() {
				t.Errorf("Expected valid version %d, got %v", tID.id.Version(), u.Version())
			}
		})

		t.Run(fmt.Sprintf("UnmarshalBSONValue Wrong Type UUID Version %d", tID.version), func(t *testing.T) {
			u := UUID{}
			raw := tID.id.Bytes()
			raw = append(uuidBSONSignature, raw...)
			err := u.UnmarshalBSONValue(bson.TypeNull, raw)
			if err == nil {
				t.Errorf("Error not returned for wrong bson Type in UnmarshalBSONValue")
			}
		})

		t.Run(fmt.Sprintf("UnmarshalBSONValue Wrong byte slice length %d", tID.version), func(t *testing.T) {
			u := UUID{}
			raw := tID.id.Bytes()[1:]
			raw = append(uuidBSONSignature, raw...)
			err := u.UnmarshalBSONValue(bson.TypeBinary, raw)
			if err == nil {
				t.Errorf("Error not returned for wrong binary data length in UnmarshalBSONValue")
			}
		})

	}
}
