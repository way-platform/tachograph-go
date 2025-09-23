package tachograph

import (
	"encoding/binary"

	cardpb "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
	vupb "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// FileType represents the type of a tachograph file.
type FileType string

const (
	// UnknownFileType is the default file type.
	UnknownFileType FileType = "unknown"
	// CardFileType indicates a driver card file.
	CardFileType FileType = "card"
	// UnitFileType indicates a vehicle unit file.
	UnitFileType FileType = "unit"
)

// InferFileType determines the type of a tachograph file based on its content.
func InferFileType(data []byte) FileType {
	if len(data) < 2 {
		return UnknownFileType
	}

	// Check for the card file marker.
	// According to Appendix 7, Section 3.4.2 of the tachograph regulation,
	// a downloaded card file is a concatenation of Elementary Files (EFs).
	// Each EF is preceded by a 2-byte tag and a 2-byte length.
	// Section 3.3.2 mandates that the first file downloaded is always EF_ICC,
	// which has the File Identifier (tag) 0x0002.
	opts := cardpb.ElementaryFileType_EF_ICC.Descriptor().Values().ByNumber(1).Options()
	efIccTag := proto.GetExtension(opts, cardpb.E_FileId).(int32)

	firstTag := binary.BigEndian.Uint16(data[0:2])
	if firstTag == uint16(efIccTag) {
		return CardFileType
	}

	// Check for VU file markers
	// VU files use TV format with 2-byte tags starting with 0x76xx
	// Check if this looks like a VU tag by examining the first byte
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
					return UnitFileType
				}
			}
		}
	}

	return UnknownFileType
}

func InferRawCardFileType(input *cardv1.RawCardFile) tachographv1.File_Type {
	if input == nil || len(input.GetRecords()) == 0 {
		return tachographv1.File_TYPE_UNSPECIFIED
	}
	// Test each card type in priority order using dual-cursor matching
	cardTypes := []tachographv1.File_Type{
		tachographv1.File_DRIVER_CARD,
		tachographv1.File_WORKSHOP_CARD,
		tachographv1.File_CONTROL_CARD,
		tachographv1.File_COMPANY_CARD,
	}
	for _, cardType := range cardTypes {
		if matchesCardTypeSequentially(input.GetRecords(), cardType) {
			return cardType
		}
	}
	return tachographv1.File_TYPE_UNSPECIFIED
}

// matchesCardTypeSequentially uses dual-cursor approach to match raw records against protobuf field structure
func matchesCardTypeSequentially(records []*cardv1.RawCardFile_Record, cardType tachographv1.File_Type) bool {
	msgType := getMessageTypeForCardType(cardType)
	if msgType == nil {
		return false
	}

	fields := msgType.Descriptor().Fields()
	recordCursor := 0 // Current position in raw records
	fieldCursor := 0  // Current position in protobuf fields

	// Walk through both sequences simultaneously
	for recordCursor < len(records) && fieldCursor < fields.Len() {
		record := records[recordCursor]
		field := fields.Get(fieldCursor)

		// Get expected EF from current field
		expectedEF := getExpectedEF(field)
		if expectedEF == cardv1.ElementaryFileType_ELEMENTARY_FILE_UNSPECIFIED {
			fieldCursor++
			continue // Skip non-EF fields
		}

		actualEF := record.GetFile()

		if actualEF == expectedEF {
			// Perfect match - advance both cursors
			recordCursor++
			fieldCursor++
		} else if isFieldConditional(field) {
			// Expected EF is conditional and missing - advance field cursor only
			fieldCursor++
		} else if isFieldOptional(field) {
			// Expected EF is optional and missing - advance field cursor only
			fieldCursor++
		} else {
			// Required field missing - this card type doesn't match
			return false
		}
	}

	// Check if we consumed all required fields
	return validateRemainingFields(fields, fieldCursor)
}

// getMessageTypeForCardType returns the protobuf message type for a given card type
func getMessageTypeForCardType(cardType tachographv1.File_Type) protoreflect.MessageType {
	switch cardType {
	case tachographv1.File_DRIVER_CARD:
		return (&cardv1.DriverCardFile{}).ProtoReflect().Type()
	case tachographv1.File_WORKSHOP_CARD:
		return (&cardv1.WorkshopCardFile{}).ProtoReflect().Type()
	case tachographv1.File_CONTROL_CARD:
		return (&cardv1.ControlCardFile{}).ProtoReflect().Type()
	case tachographv1.File_COMPANY_CARD:
		return (&cardv1.CompanyCardFile{}).ProtoReflect().Type()
	default:
		return nil
	}
}

// getExpectedEF maps protobuf field names to ElementaryFileType
func getExpectedEF(field protoreflect.FieldDescriptor) cardv1.ElementaryFileType {
	// Map field names to their corresponding ElementaryFileType
	fieldName := string(field.Name())
	switch fieldName {
	case "icc":
		return cardv1.ElementaryFileType_EF_ICC
	case "ic":
		return cardv1.ElementaryFileType_EF_IC
	case "application_identification":
		return cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION
	case "identification":
		return cardv1.ElementaryFileType_EF_IDENTIFICATION
	case "holder_identification":
		return cardv1.ElementaryFileType_EF_IDENTIFICATION // Same EF, different part
	case "driving_licence_info":
		return cardv1.ElementaryFileType_EF_DRIVING_LICENCE_INFO
	case "events_data":
		return cardv1.ElementaryFileType_EF_EVENTS_DATA
	case "faults_data":
		return cardv1.ElementaryFileType_EF_FAULTS_DATA
	case "driver_activity_data":
		return cardv1.ElementaryFileType_EF_DRIVER_ACTIVITY_DATA
	case "vehicles_used":
		return cardv1.ElementaryFileType_EF_VEHICLES_USED
	case "places":
		return cardv1.ElementaryFileType_EF_PLACES
	case "current_usage":
		return cardv1.ElementaryFileType_EF_CURRENT_USAGE
	case "control_activity_data":
		return cardv1.ElementaryFileType_EF_CONTROL_ACTIVITY_DATA
	case "calibration":
		return cardv1.ElementaryFileType_EF_CALIBRATION
	case "sensor_installation_data":
		return cardv1.ElementaryFileType_EF_SENSOR_INSTALLATION_DATA
	case "controller_activity_data":
		return cardv1.ElementaryFileType_EF_CONTROLLER_ACTIVITY_DATA
	case "company_activity_data":
		return cardv1.ElementaryFileType_EF_COMPANY_ACTIVITY_DATA
	case "specific_conditions":
		return cardv1.ElementaryFileType_EF_SPECIFIC_CONDITIONS
	case "vehicle_units_used":
		return cardv1.ElementaryFileType_EF_VEHICLE_UNITS_USED
	case "gnss_places":
		return cardv1.ElementaryFileType_EF_GNSS_PLACES
	case "application_identification_v2":
		return cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION_V2
	case "places_authentication":
		return cardv1.ElementaryFileType_EF_PLACES_AUTHENTICATION
	case "gnss_places_authentication":
		return cardv1.ElementaryFileType_EF_GNSS_PLACES_AUTHENTICATION
	case "border_crossings":
		return cardv1.ElementaryFileType_EF_BORDER_CROSSINGS
	case "load_unload_operations":
		return cardv1.ElementaryFileType_EF_LOAD_UNLOAD_OPERATIONS
	case "load_type_entries":
		return cardv1.ElementaryFileType_EF_LOAD_TYPE_ENTRIES
	case "calibration_add_data":
		return cardv1.ElementaryFileType_EF_CALIBRATION_ADD_DATA
	case "vu_configuration":
		return cardv1.ElementaryFileType_EF_VU_CONFIGURATION
	default:
		return cardv1.ElementaryFileType_ELEMENTARY_FILE_UNSPECIFIED
	}
}

// isFieldConditional checks if a protobuf field represents a conditional EF
func isFieldConditional(field protoreflect.FieldDescriptor) bool {
	// Define conditional fields based on regulation knowledge
	fieldName := string(field.Name())
	conditionalFields := map[string]bool{
		"application_identification_v2": true, // Gen2 only
		"places_authentication":         true, // Gen2 only
		"gnss_places_authentication":    true, // Gen2 only
		"border_crossings":              true, // Gen2 only
		"load_unload_operations":        true, // Gen2 only
		"load_type_entries":             true, // Gen2 only
		"calibration_add_data":          true, // Workshop Gen2 only
		"vu_configuration":              true, // Control Gen2 only
		"gnss_places":                   true, // Gen2 only
		"vehicle_units_used":            true, // Gen2 only
	}
	return conditionalFields[fieldName]
}

// isFieldOptional checks if a protobuf field is optional
func isFieldOptional(field protoreflect.FieldDescriptor) bool {
	// In proto3, fields without presence are typically optional
	// Fields with presence (like message types) are more important
	return !field.HasPresence() || field.Cardinality() == protoreflect.Optional
}

// validateRemainingFields ensures all remaining fields are optional or conditional
func validateRemainingFields(fields protoreflect.FieldDescriptors, fieldCursor int) bool {
	for i := fieldCursor; i < fields.Len(); i++ {
		field := fields.Get(i)
		expectedEF := getExpectedEF(field)
		if expectedEF != cardv1.ElementaryFileType_ELEMENTARY_FILE_UNSPECIFIED {
			if !isFieldOptional(field) && !isFieldConditional(field) {
				return false // Required field not processed
			}
		}
	}
	return true
}
