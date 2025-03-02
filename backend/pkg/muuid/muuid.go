package muuid

import (
	"fmt"
	"reflect"

	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

var (
	uuidSubtype = byte(0x04)
)

// MongoUUID is a ValueEncoder and ValueDecoder for encoding uuid.UUID to Binary and vice versa.
type MongoUUID struct{}

var _ bsoncodec.ValueEncoder = &MongoUUID{}
var _ bsoncodec.ValueDecoder = &MongoUUID{}

func (mu *MongoUUID) EncodeValue(ect bsoncodec.EncodeContext, w bsonrw.ValueWriter, value reflect.Value) error {
	// Use reflection to convert the value to uuid.UUID.
	u, ok := value.Interface().(uuid.UUID)
	if !ok {
		return fmt.Errorf("value %v to encode is not of type uuid.UUID", value)
	}

	// Convert uuid.UUID to binary with subtype.
	return w.WriteBinaryWithSubtype(u[:], uuidSubtype)
}

func (mu *MongoUUID) DecodeValue(ect bsoncodec.DecodeContext, r bsonrw.ValueReader, value reflect.Value) error {
	var data []byte
	var subtype byte
	var err error

	switch rType := r.Type(); rType {
	case bsontype.Binary:
		data, subtype, err = r.ReadBinary()
		if subtype != uuidSubtype {
			return fmt.Errorf("unsupported binary subtype %v for UUID", subtype)
		}
	case bsontype.Null:
		err = r.ReadNull()
	case bsontype.Undefined:
		err = r.ReadUndefined()
	default:
		return fmt.Errorf("cannot decode %v into a UUID", rType)
	}

	if err != nil {
		return fmt.Errorf("error reading binary data from ValueReader: %v", err)
	}

	// Convert binary data to uuid.UUID.
	u, err := uuid.FromBytes(data)
	if err != nil {
		return fmt.Errorf("error converting binary data %v to uuid.UUID: %v", data, err)
	}

	// Set the value to the decoded uuid.UUID.
	value.Set(reflect.ValueOf(u))
	return nil
}
