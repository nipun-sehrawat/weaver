// go:build !ignoreWeaverGen

package model

// Code generated by "weaver generate". DO NOT EDIT.
import (
	"fmt"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
)


// Local stub implementations.

// Client stub implementations.

// Server stub implementations.

// AutoMarshal implementations.

var _ codegen.AutoMarshal = &Transaction{}

func (x *Transaction) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Transaction.WeaverMarshal: nil receiver"))
	}
	enc.String(x.FromAccountNum)
	enc.String(x.FromRoutingNum)
	enc.String(x.ToAccountNum)
	enc.String(x.ToRoutingNum)
	enc.Int64(x.Amount)
	enc.EncodeBinaryMarshaler(&x.Timestamp)
}

func (x *Transaction) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Transaction.WeaverUnmarshal: nil receiver"))
	}
	x.FromAccountNum = dec.String()
	x.FromRoutingNum = dec.String()
	x.ToAccountNum = dec.String()
	x.ToRoutingNum = dec.String()
	x.Amount = dec.Int64()
	dec.DecodeBinaryUnmarshaler(&x.Timestamp)
}

var _ codegen.AutoMarshal = &TransactionWithID{}

func (x *TransactionWithID) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("TransactionWithID.WeaverMarshal: nil receiver"))
	}
	(x.Transaction).WeaverMarshal(enc)
	enc.Int64(x.TransactionID)
}

func (x *TransactionWithID) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("TransactionWithID.WeaverUnmarshal: nil receiver"))
	}
	(&x.Transaction).WeaverUnmarshal(dec)
	x.TransactionID = dec.Int64()
}