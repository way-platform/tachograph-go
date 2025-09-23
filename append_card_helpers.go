package tachograph

import (
	"encoding/binary"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func toBCD(val int) byte {
	return byte(((val / 10) << 4) | (val % 10))
}

func appendDatef(dst []byte, t *timestamppb.Timestamp) []byte {
	if t == nil {
		return append(dst, 0, 0, 0, 0)
	}
	year := t.AsTime().Year()
	month := int(t.AsTime().Month())
	day := t.AsTime().Day()

	dst = append(dst, toBCD(year/100), toBCD(year%100))
	dst = append(dst, toBCD(month), toBCD(day))
	return dst
}

func appendString(dst []byte, s string, fixedLen int) []byte {
	b := []byte(s)
	if len(b) > fixedLen {
		b = b[:fixedLen]
	}
	dst = append(dst, b...)
	for i := len(b); i < fixedLen; i++ {
		dst = append(dst, ' ') // Pad with spaces
	}
	return dst
}

func appendTimeReal(dst []byte, t *timestamppb.Timestamp) []byte {
	var timeVal uint32
	if t != nil {
		timeVal = uint32(t.GetSeconds())
	}
	return binary.BigEndian.AppendUint32(dst, timeVal)
}

func appendOdometer(dst []byte, km int32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(km))
	return append(dst, b[1:]...) // Append 3 bytes
}
