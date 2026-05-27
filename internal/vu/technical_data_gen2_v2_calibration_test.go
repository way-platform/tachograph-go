package vu

import (
	"testing"

	"github.com/way-platform/tachograph-go/internal/dd"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// TestParseOneCalibrationRecordGen2V2_CalibrationCountry guards against the same
// enum-mapping regression as the vehicle registration nation: the calibration
// country byte (offset 247) must be mapped through the protocol_enum_value
// annotation, not cast directly to the proto enum. With the raw cast, a
// Norwegian calibration (protocol value 37) decoded as MONACO (proto enum value
// 37). See tacho#1090.
//
// A 252-byte all-zero record parses cleanly (lenient sub-parsers); only the
// calibration country byte is varied.
func TestParseOneCalibrationRecordGen2V2_CalibrationCountry(t *testing.T) {
	const lenRecord = 252
	const idxCalCountry = 247

	tests := []struct {
		name        string
		countryByte byte
		want        ddv1.NationNumeric
	}{
		{name: "Norway", countryByte: 0x25, want: ddv1.NationNumeric_NORWAY},
		{name: "Monaco", countryByte: 0x22, want: ddv1.NationNumeric_MONACO},
		{name: "Default (0x00)", countryByte: 0x00, want: ddv1.NationNumeric_NATION_NUMERIC_DEFAULT},
		{name: "Unrecognized (100)", countryByte: 0x64, want: ddv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]byte, lenRecord)
			data[idxCalCountry] = tt.countryByte

			var opts dd.UnmarshalOptions
			rec, err := parseOneCalibrationRecordGen2V2(opts, data)
			if err != nil {
				t.Fatalf("parseOneCalibrationRecordGen2V2() unexpected error: %v", err)
			}
			if rec.GetCalibrationCountry() != tt.want {
				t.Errorf("GetCalibrationCountry() = %v, want %v", rec.GetCalibrationCountry(), tt.want)
			}
		})
	}
}
