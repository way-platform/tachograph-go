package dd

import (
	"encoding/binary"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// UnmarshalTimeReal unmarshals a TimeReal timestamp from a byte slice.
//
// The data type `TimeReal` is specified in the Data Dictionary, Section 2.162.
//
// ASN.1 Definition:
//
//	TimeReal ::= INTEGER (0..2^32-1)
//
// Binary Layout (4 bytes):
//   - Seconds since Unix epoch (4 bytes): Big-endian uint32
func (opts UnmarshalOptions) UnmarshalTimeReal(data []byte) (*timestamppb.Timestamp, error) {
	const lenTimeReal = 4
	if len(data) != lenTimeReal {
		return nil, fmt.Errorf("invalid data length for TimeReal: got %d, want %d", len(data), lenTimeReal)
	}
	timeVal := binary.BigEndian.Uint32(data[:lenTimeReal])
	if timeVal == 0 {
		return nil, nil // Zero time is represented as nil
	}
	return timestamppb.New(time.Unix(int64(timeVal), 0)), nil
}

// AppendTimeReal appends a 4-byte TimeReal value.
//
// The data type `TimeReal` is specified in the Data Dictionary, Section 2.162.
//
// ASN.1 Definition:
//
//	TimeReal ::= INTEGER (0..2^32-1)
//
// Binary Layout (4 bytes):
//   - Seconds since Unix epoch (4 bytes): Big-endian uint32
func AppendTimeReal(dst []byte, ts *timestamppb.Timestamp) ([]byte, error) {
	if ts.GetNanos() > 0 {
		return nil, fmt.Errorf("nanosecond resolution is not supported for TimeReal")
	}
	return binary.BigEndian.AppendUint32(dst, uint32(ts.GetSeconds())), nil
}
