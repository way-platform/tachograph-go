package card

import (
	"testing"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestUnmarshalBorderCrossings_Synthetic(t *testing.T) {
	// Construct a synthetic raw buffer for EF_Border_Crossings
	// 1 byte pointer + 1 record (17 bytes)
	
	// Record:
	// Country Left: 0x12 (18 - Spain)
	// Country Entered: 0x14 (20 - France)
	// GNSSPlaceAuthRecord (12 bytes):
	//   Time: 0x00 0x00 0x00 0x0A (10s)
	//   Accuracy: 0x00
	//   Lat: 0x00 0x00 0x00
	//   Lon: 0x00 0x00 0x00
	// Odometer: 0x00 0x00 0x64 (100km)
	
	record := []byte{
		0x0F, // Country Left (Spain - 0x0F)
		0x11, // Country Entered (France - 0x11)
		// GNSS Place Auth Record (12 bytes)
		0x00, 0x00, 0x00, 0x0A, // Time (4)
		0x00, // Accuracy (1)
		0x00, 0x00, 0x00, // Lat (3)
		0x00, 0x00, 0x00, // Lon (3)
		0x00, // AuthenticationStatus (1) - Not Authenticated
		// Odometer
		0x00, 0x00, 0x64,
	}
	
	data := append([]byte{0x00}, record...) // Pointer 0, 1 record

	opts := UnmarshalOptions{}
	bc, err := opts.unmarshalBorderCrossings(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bc == nil {
		t.Fatal("expected non-nil result")
	}
	if got := bc.GetNewestRecordIndex(); got != 0 {
		t.Errorf("NewestRecordIndex: got %d, want 0", got)
	}
	if len(bc.GetRecords()) != 1 {
		t.Fatalf("Records length: got %d, want 1", len(bc.GetRecords()))
	}
	
	r := bc.GetRecords()[0]
	if got := r.GetCountryLeft(); got != ddv1.NationNumeric_SPAIN {
		t.Errorf("CountryLeft: got %v, want %v", got, ddv1.NationNumeric_SPAIN)
	}
	if got := r.GetCountryEntered(); got != ddv1.NationNumeric_FRANCE {
		t.Errorf("CountryEntered: got %v, want %v", got, ddv1.NationNumeric_FRANCE)
	}
	if got := r.GetVehicleOdometerKm(); got != 100 {
		t.Errorf("VehicleOdometerKm: got %d, want 100", got)
	}
	if got := r.GetGnssPlaceAuthRecord().GetTimestamp().GetSeconds(); got != 10 {
		t.Errorf("Timestamp Seconds: got %d, want 10", got)
	}
	
	// Default proto value 0 maps to UNSPECIFIED, but UnmarshalEnum logic mapped 0x00 to NOT_AUTHENTICATED (2).
	if got := r.GetGnssPlaceAuthRecord().GetAuthenticationStatus(); got != ddv1.PositionAuthenticationStatus_NOT_AUTHENTICATED {
		t.Errorf("AuthenticationStatus: got %v, want %v", got, ddv1.PositionAuthenticationStatus_NOT_AUTHENTICATED)
	}
}

func TestMarshalBorderCrossings_Synthetic(t *testing.T) {
	record := &cardv1.BorderCrossings_Record{}
	record.SetCountryLeft(ddv1.NationNumeric_SPAIN)
	record.SetCountryEntered(ddv1.NationNumeric_FRANCE)
	record.SetVehicleOdometerKm(100)
	
	gnss := &ddv1.GNSSPlaceAuthRecord{}
	gnss.SetTimestamp(&timestamppb.Timestamp{Seconds: 10})
	gnss.SetAuthenticationStatus(ddv1.PositionAuthenticationStatus_NOT_AUTHENTICATED)
	
	coords := &ddv1.GeoCoordinates{}
	coords.SetLatitude(0)
	coords.SetLongitude(0)
	gnss.SetGeoCoordinates(coords)
	
	record.SetGnssPlaceAuthRecord(gnss)

	msg := &cardv1.BorderCrossings{}
	msg.SetNewestRecordIndex(0)
	msg.SetRecords([]*cardv1.BorderCrossings_Record{record})

	opts := MarshalOptions{}
	data, err := opts.MarshalCardBorderCrossings(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if len(data) != 18 {
		t.Fatalf("data length: got %d, want 18", len(data))
	}
	if data[0] != 0x00 {
		t.Errorf("pointer: got 0x%02x, want 0x00", data[0])
	}
	if data[1] != 0x0F {
		t.Errorf("Country Left: got 0x%02x, want 0x0F", data[1])
	}
	if data[2] != 0x11 {
		t.Errorf("Country Entered: got 0x%02x, want 0x11", data[2])
	}
	// Verify odometer at end
	if data[17] != 0x64 {
		t.Errorf("Odometer byte: got 0x%02x, want 0x64", data[17])
	}
}
