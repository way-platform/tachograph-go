package tachograph

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalEntryTypeDailyWorkPeriod parses entry type daily work period from raw data.
//
// The data type `EntryTypeDailyWorkPeriod` is specified in the Data Dictionary, Section 2.66.
//
// ASN.1 Definition:
//
//	EntryTypeDailyWorkPeriod ::= INTEGER {
//	    begin(0), end(1), relatedToGNSS(2), relatedToITS(3)
//	} (0..255)
//
// Binary Layout (1 byte):
//   - Entry Type Daily Work Period (1 byte): Raw integer value (0-3)
func unmarshalEntryTypeDailyWorkPeriod(data []byte) (ddv1.EntryTypeDailyWorkPeriod, error) {
	if len(data) < 1 {
		return ddv1.EntryTypeDailyWorkPeriod_ENTRY_TYPE_DAILY_WORK_PERIOD_UNSPECIFIED, fmt.Errorf("insufficient data for EntryTypeDailyWorkPeriod: got %d, want 1", len(data))
	}

	rawValue := int32(data[0])

	// Use the protocol enum value mapping
	entryTypeDailyWorkPeriod := ddv1.EntryTypeDailyWorkPeriod_ENTRY_TYPE_DAILY_WORK_PERIOD_UNSPECIFIED
	SetEntryTypeDailyWorkPeriod(ddv1.EntryTypeDailyWorkPeriod_ENTRY_TYPE_DAILY_WORK_PERIOD_UNSPECIFIED.Descriptor(), rawValue,
		func(enumNum protoreflect.EnumNumber) {
			entryTypeDailyWorkPeriod = ddv1.EntryTypeDailyWorkPeriod(enumNum)
		}, func(unrecognized int32) {
			entryTypeDailyWorkPeriod = ddv1.EntryTypeDailyWorkPeriod_ENTRY_TYPE_DAILY_WORK_PERIOD_UNRECOGNIZED
		})

	return entryTypeDailyWorkPeriod, nil
}

// appendEntryTypeDailyWorkPeriod appends entry type daily work period as a single byte.
//
// The data type `EntryTypeDailyWorkPeriod` is specified in the Data Dictionary, Section 2.66.
//
// ASN.1 Definition:
//
//	EntryTypeDailyWorkPeriod ::= INTEGER {
//	    begin(0), end(1), relatedToGNSS(2), relatedToITS(3)
//	} (0..255)
//
// Binary Layout (1 byte):
//   - Entry Type Daily Work Period (1 byte): Raw integer value (0-3)
func appendEntryTypeDailyWorkPeriod(dst []byte, entryTypeDailyWorkPeriod ddv1.EntryTypeDailyWorkPeriod) []byte {
	// Get the protocol value for the enum
	protocolValue := GetEntryTypeDailyWorkPeriodProtocolValue(entryTypeDailyWorkPeriod, 0)
	return append(dst, byte(protocolValue))
}
