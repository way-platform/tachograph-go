package dd

import (
	"bytes"
	"testing"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// TestUnmarshalFullCardNumberControlCard verifies that a FullCardNumber with
// equipment type CONTROL_CARD (3) is parsed using OwnerIdentification, as
// specified in the Data Dictionary, Section 2.73 (CardNumber CHOICE).
//
// Control card numbers appear in VuControlActivityData records of any vehicle
// unit that has ever been checked at the roadside, so rejecting them made
// Overview Gen1 parsing fail for most real-world VU downloads.
func TestUnmarshalFullCardNumberControlCard(t *testing.T) {
	// EquipmentType (1) + NationNumeric (1) + OwnerIdentification (13+1+1+1)
	data := make([]byte, 0, 18)
	data = append(data, 0x03) // CONTROL_CARD
	data = append(data, 0x0F) // Spain
	data = append(data, []byte("C000012345678")...)
	data = append(data, '0', '0', '1')

	opts := UnmarshalOptions{}
	got, err := opts.UnmarshalFullCardNumber(data)
	if err != nil {
		t.Fatalf("UnmarshalFullCardNumber(CONTROL_CARD): %v", err)
	}
	if got.GetCardType() != ddv1.EquipmentType_CONTROL_CARD {
		t.Errorf("card type: got %v, want CONTROL_CARD", got.GetCardType())
	}
	owner := got.GetOwnerIdentification()
	if owner == nil {
		t.Fatal("owner identification: got nil")
	}
	if v := owner.GetOwnerIdentification().GetValue(); v != "C000012345678" {
		t.Errorf("owner identification: got %q, want %q", v, "C000012345678")
	}

	// Round-trip.
	out, err := MarshalOptions{}.MarshalFullCardNumber(got)
	if err != nil {
		t.Fatalf("MarshalFullCardNumber(CONTROL_CARD): %v", err)
	}
	if !bytes.Equal(out, data) {
		t.Errorf("round-trip mismatch:\n got %x\nwant %x", out, data)
	}
}
