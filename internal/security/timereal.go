package security

import (
	"encoding/binary"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// unmarshalTimeReal converts a 4-byte TimeReal value to a timestamp.
func unmarshalTimeReal(data []byte) (*timestamppb.Timestamp, error) {
	if len(data) != 4 {
		return nil, fmt.Errorf("invalid TimeReal length: got %d, want 4", len(data))
	}
	seconds := binary.BigEndian.Uint32(data)
	// TimeReal is seconds since 1970-01-01 00:00:00 UTC
	epoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	return timestamppb.New(epoch.Add(time.Duration(seconds) * time.Second)), nil
}
