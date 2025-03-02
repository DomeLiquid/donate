package bigdecimal

import (
	"fmt"
	"reflect"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoDecimal is a ValueEncoder and ValueDecoder for encoding decimal.Decimal to Decimal128 and vice versa.
type MongoDecimal struct{}

var _ bsoncodec.ValueEncoder = &MongoDecimal{}
var _ bsoncodec.ValueDecoder = &MongoDecimal{}

func (dc *MongoDecimal) EncodeValue(ect bsoncodec.EncodeContext, w bsonrw.ValueWriter, value reflect.Value) error {
	// Use reflection to convert the value to decimal.Decimal.
	dec, ok := value.Interface().(decimal.Decimal)
	if !ok {
		return fmt.Errorf("value %v to encode is not of type decimal.Decimal", value)
	}

	// Convert decimal.Decimal to primitive.Decimal128.
	primDec, err := primitive.ParseDecimal128(dec.Truncate(28).String())
	if err != nil {
		return fmt.Errorf("error converting decimal.Decimal %v to primitive.Decimal128: %v", dec, err)
	}
	return w.WriteDecimal128(primDec)
}

func (dc *MongoDecimal) DecodeValue(ect bsoncodec.DecodeContext, r bsonrw.ValueReader, value reflect.Value) error {
	primDec, err := r.ReadDecimal128()
	if err != nil {
		return fmt.Errorf("error reading primitive.Decimal128 from ValueReader: %v", err)
	}

	// Convert primitive.Decimal128 to Golang's decimal.Decimal.
	dec, err := decimal.NewFromString(primDec.String())
	if err != nil {
		return fmt.Errorf("error converting primitive.Decimal128 %v to decimal.Decimal: %v", primDec, err)
	}

	// Set the value to the decoded decimal.Decimal.
	value.Set(reflect.ValueOf(dec))
	return nil
}
