package uuid

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

func (u UUID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.TypeBinary, bsoncore.AppendBinary(nil, 4, u[:]), nil
}

func (u *UUID) UnmarshalBSONValue(t bsontype.Type, raw []byte) error {
	if t != bson.TypeBinary {
		return fmt.Errorf("invalid format on unmarshal bson value")
	}

	_, data, _, ok := bsoncore.ReadBinary(raw)
	if !ok {
		return fmt.Errorf("not enough bytes to unmarshal bson value")
	}

	copy(u[:], data)

	return nil
}
