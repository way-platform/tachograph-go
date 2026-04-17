package card

import (
	"testing"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// TestParseRawDriverCardFile_BorderCrossings verifies that EF_BORDER_CROSSINGS records
// are wired into ParseRawDriverCardFile and propagated to TachographG2.border_crossings.
//
// Uses the same synthetic byte layout as TestUnmarshalBorderCrossings_Synthetic:
//
//	1-byte pointer (0x00) + 1 record (17 bytes):
//	  [0x0F]               country_left  = Spain
//	  [0x11]               country_entered = France
//	  [00 00 00 0A]        GNSS timestamp = 10 s since midnight 1970-01-01
//	  [00]                 GNSS accuracy
//	  [00 00 00]           GNSS latitude  = 0
//	  [00 00 00]           GNSS longitude = 0
//	  [00]                 authentication status = NOT_AUTHENTICATED
//	  [00 00 64]           odometer = 100 km
func TestParseRawDriverCardFile_BorderCrossings(t *testing.T) {
	borderCrossingsData := []byte{
		0x00,                               // pointer: newest record index = 0
		0x0F, 0x11,                         // country_left = Spain, country_entered = France
		0x00, 0x00, 0x00, 0x0A,            // GNSS timestamp (4 bytes)
		0x00,                               // GNSS accuracy
		0x00, 0x00, 0x00,                  // GNSS latitude
		0x00, 0x00, 0x00,                  // GNSS longitude
		0x00,                               // GNSS authentication status = NOT_AUTHENTICATED
		0x00, 0x00, 0x64,                  // odometer = 100 km
	}

	rec := &cardv1.RawCardFile_Record{}
	rec.SetFile(cardv1.ElementaryFileType_EF_BORDER_CROSSINGS)
	rec.SetGeneration(ddv1.Generation_GENERATION_2)
	rec.SetContentType(cardv1.ContentType_DATA)
	rec.SetValue(borderCrossingsData)

	rawFile := &cardv1.RawCardFile{}
	rawFile.SetRecords([]*cardv1.RawCardFile_Record{rec})

	cardFile, err := ParseOptions{}.ParseRawDriverCardFile(rawFile)
	if err != nil {
		t.Fatalf("ParseRawDriverCardFile returned unexpected error: %v", err)
	}

	g2 := cardFile.GetTachographG2()
	if g2 == nil {
		t.Fatal("TachographG2 is nil — EF_BORDER_CROSSINGS was not wired into the switch")
	}

	bc := g2.GetBorderCrossings()
	if bc == nil {
		t.Fatal("BorderCrossings is nil — EF_BORDER_CROSSINGS case is missing from the switch")
	}

	if got := bc.GetNewestRecordIndex(); got != 0 {
		t.Errorf("NewestRecordIndex: got %d, want 0", got)
	}

	records := bc.GetRecords()
	if len(records) != 1 {
		t.Fatalf("Records count: got %d, want 1", len(records))
	}

	r := records[0]
	if got, want := r.GetCountryLeft(), ddv1.NationNumeric_SPAIN; got != want {
		t.Errorf("CountryLeft: got %v, want %v", got, want)
	}
	if got, want := r.GetCountryEntered(), ddv1.NationNumeric_FRANCE; got != want {
		t.Errorf("CountryEntered: got %v, want %v", got, want)
	}
	if got, want := r.GetVehicleOdometerKm(), int32(100); got != want {
		t.Errorf("VehicleOdometerKm: got %d, want %d", got, want)
	}
}
