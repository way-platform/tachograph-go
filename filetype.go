package tachograph

import (
	"encoding/binary"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
	vupb "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
)

// inferFileType determines the type of a tachograph file based on its content.
func inferFileType(data []byte) tachographv1.File_Type {
	if len(data) < 2 {
		return tachographv1.File_TYPE_UNSPECIFIED
	}
	// Check for the card file marker.
	// According to Appendix 7, Section 3.4.2 of the tachograph regulation,
	// a downloaded card file is a concatenation of Elementary Files (EFs).
	// Each EF is preceded by a 2-byte tag and a 2-byte length.
	// Section 3.3.2 mandates that the first file downloaded is always EF_ICC,
	// which has the File Identifier (tag) 0x0002.
	opts := cardv1.ElementaryFileType_EF_ICC.Descriptor().Values().ByNumber(1).Options()
	efIccTag := proto.GetExtension(opts, cardv1.E_FileId).(int32)
	firstTag := binary.BigEndian.Uint16(data[0:2])
	if firstTag == uint16(efIccTag) {
		return tachographv1.File_DRIVER_CARD
	}
	// Check for VU file markers
	// VU files can use either TV format (0x76XX) or TRTP format (raw TREP values)
	// First check for TV format with 2-byte tags starting with 0x76xx
	if data[0] == 0x76 {
		// Check if the second byte corresponds to a valid TREP value
		secondByte := data[1]
		// Check against known VU transfer types
		values := vupb.TransferType_TRANSFER_TYPE_UNSPECIFIED.Descriptor().Values()
		for i := 0; i < values.Len(); i++ {
			valueDesc := values.Get(i)
			opts := valueDesc.Options()
			if proto.HasExtension(opts, vupb.E_TrepValue) {
				trepValue := proto.GetExtension(opts, vupb.E_TrepValue).(int32)
				if uint8(trepValue) == secondByte {
					return tachographv1.File_VEHICLE_UNIT
				}
			}
		}
	}

	// Also check for TRTP format where the first byte is directly a TREP value
	firstByte := data[0]
	// Check against known VU transfer types
	values := vupb.TransferType_TRANSFER_TYPE_UNSPECIFIED.Descriptor().Values()
	for i := 0; i < values.Len(); i++ {
		valueDesc := values.Get(i)
		opts := valueDesc.Options()
		if proto.HasExtension(opts, vupb.E_TrepValue) {
			trepValue := proto.GetExtension(opts, vupb.E_TrepValue).(int32)
			if uint8(trepValue) == firstByte {
				return tachographv1.File_VEHICLE_UNIT
			}
		}
	}
	// Also check for manufacturer-specific TREP values (0x11-0x1F)
	if firstByte >= 0x11 && firstByte <= 0x1F {
		return tachographv1.File_VEHICLE_UNIT
	}
	return tachographv1.File_TYPE_UNSPECIFIED
}
