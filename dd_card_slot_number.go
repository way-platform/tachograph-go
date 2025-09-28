package tachograph

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalCardSlotNumber parses card slot number from raw data.
//
// The data type `CardSlotNumber` is specified in the Data Dictionary, Section 2.33.
//
// ASN.1 Definition:
//
//	CardSlotNumber ::= INTEGER (0..1)
//
// Binary Layout (1 bit):
//   - Card Slot Number (1 bit): Raw integer value (0-1)
func unmarshalCardSlotNumber(data []byte) (ddv1.CardSlotNumber, error) {
	if len(data) < 1 {
		return ddv1.CardSlotNumber_CARD_SLOT_NUMBER_UNSPECIFIED, fmt.Errorf("insufficient data for CardSlotNumber: got %d, want 1", len(data))
	}

	rawValue := int32(data[0])

	// Use the protocol enum value mapping
	cardSlotNumber := ddv1.CardSlotNumber_CARD_SLOT_NUMBER_UNSPECIFIED
	SetCardSlotNumber(ddv1.CardSlotNumber_CARD_SLOT_NUMBER_UNSPECIFIED.Descriptor(), rawValue, func(en protoreflect.EnumNumber) {
		cardSlotNumber = ddv1.CardSlotNumber(en)
	}, func(unrecognized int32) {
		cardSlotNumber = ddv1.CardSlotNumber_CARD_SLOT_NUMBER_UNRECOGNIZED
	})

	return cardSlotNumber, nil
}

// appendCardSlotNumber appends card slot number as a single byte.
//
// The data type `CardSlotNumber` is specified in the Data Dictionary, Section 2.33.
//
// ASN.1 Definition:
//
//	CardSlotNumber ::= INTEGER (0..1)
//
// Binary Layout (1 bit):
//   - Card Slot Number (1 bit): Raw integer value (0-1)
func appendCardSlotNumber(dst []byte, cardSlotNumber ddv1.CardSlotNumber) []byte {
	// Get the protocol value for the enum
	protocolValue := GetCardSlotNumber(cardSlotNumber, 0)
	return append(dst, byte(protocolValue))
}
