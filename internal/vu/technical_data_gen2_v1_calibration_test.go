package vu

import (
	"testing"

	"github.com/way-platform/tachograph-go/internal/dd"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// TestParseOneCalibrationRecordGen2_Alignment guards against the one-byte
// misalignment fixed for tacho#1106. The Gen2 V1 VuCalibrationRecord prefix is
// identical to Gen1 (DD 2.174): a 18-byte FullCardNumber workshop card followed
// by a TimeReal expiry — NOT a 19-byte FullCardNumberAndGeneration + Date. The
// extra phantom byte shifted the VIN, vehicle registration, and every following
// field by +1.
//
// This record places the VIN at offset 95 and the registration at offset 112
// (the spec-correct offsets, matching the Gen2 V2 record). With the phantom
// generation byte the parser read them at 96/113 and produced corrupted values.
func TestParseOneCalibrationRecordGen2_Alignment(t *testing.T) {
	const recordLen = 168
	data := make([]byte, recordLen)

	// workshopCardNumber: FullCardNumber (18 bytes) at offset 73; 0xFF = no card.
	data[73] = 0xFF

	// VIN (17 bytes) at offset 95.
	copy(data[95:112], []byte("YS2P6X20005637947"))

	// vehicleRegistration (15 bytes) at offset 112: nation + codepage + number.
	data[112] = 0x25 // Norway protocol value
	data[113] = 0x00 // code page
	copy(data[114:121], []byte("JD91529"))

	var opts dd.UnmarshalOptions
	rec, err := parseOneCalibrationRecordGen2(opts, data)
	if err != nil {
		t.Fatalf("parseOneCalibrationRecordGen2() unexpected error: %v", err)
	}

	if got := rec.vin.GetValue(); got != "YS2P6X20005637947" {
		t.Errorf("vin = %q, want %q", got, "YS2P6X20005637947")
	}
	if got := rec.vehicleRegistration.GetNation(); got != ddv1.NationNumeric_NORWAY {
		t.Errorf("vehicleRegistration nation = %v, want NORWAY", got)
	}
	if got := rec.vehicleRegistration.GetNumber().GetValue(); got != "JD91529" {
		t.Errorf("vehicleRegistration number = %q, want %q", got, "JD91529")
	}
}
