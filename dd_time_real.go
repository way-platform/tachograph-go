package tachograph

import (
	"bytes"
	"encoding/binary"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// appendTimeReal appends a 4-byte TimeReal value.
//
// The data type `TimeReal` is specified in the Data Dictionary, Section 2.162.
//
// ASN.1 Definition:
//
//	TimeReal ::= INTEGER (0..2^32-1)
func appendTimeReal(dst []byte, ts *timestamppb.Timestamp) []byte {
	if ts == nil {
		return append(dst, 0, 0, 0, 0)
	}
	return binary.BigEndian.AppendUint32(dst, uint32(ts.GetSeconds()))
}

// readTimeReal reads a TimeReal value (4 bytes) from a bytes.Reader and converts to Timestamp
//
// The data type `TimeReal` is specified in the Data Dictionary, Section 2.162.
//
// ASN.1 Definition:
//
//	TimeReal ::= INTEGER (0..2^32-1)
func readTimeReal(r *bytes.Reader) *timestamppb.Timestamp {
	var timeVal uint32
	_ = binary.Read(r, binary.BigEndian, &timeVal) // ignore error as we're reading from in-memory buffer
	if timeVal == 0 {
		return nil
	}
	return timestamppb.New(time.Unix(int64(timeVal), 0))
}

