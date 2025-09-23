package tachograph

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
	"unsafe"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
)

// Marshal serializes a protobuf File message into the binary DDD file format.
func Marshal(file *tachographv1.File) ([]byte, error) {
	// Start with a reasonably sized buffer to avoid reallocations.
	buf := make([]byte, 0, 8192)
	var err error

	switch file.GetType() {
	case tachographv1.File_DRIVER_CARD, tachographv1.File_WORKSHOP_CARD, tachographv1.File_CONTROL_CARD, tachographv1.File_COMPANY_CARD:
		buf, err = appendCard(buf, file)
	case tachographv1.File_VEHICLE_UNIT:
		buf, err = appendVU(buf, file)
	default:
		err = errors.New("unsupported file type for marshaling")
	}

	if err != nil {
		return nil, err
	}
	return buf, nil
}

// appendCard orchestrates writing a card file.
func appendCard(dst []byte, file *tachographv1.File) ([]byte, error) {
	var err error
	switch file.GetType() {
	case tachographv1.File_DRIVER_CARD:
		dst, err = appendDriverCard(dst, file.GetDriverCard())
	default:
		return nil, errors.New("unsupported card type for marshaling")
	}
	if err != nil {
		return nil, err
	}
	return dst, nil
}

// appendDriverCard orchestrates the writing of a driver card file.
func appendDriverCard(dst []byte, card *cardv1.DriverCardFile) ([]byte, error) {
	var err error

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_ICC, card.GetIcc(), AppendIcc)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_IC, card.GetIc(), AppendCardIc)
	if err != nil {
		return nil, err
	}

	// EF_Identification is a composite file, so we handle it specially.
	valBuf := make([]byte, 0, 143)
	valBuf, err = AppendCardIdentification(valBuf, card.GetIdentification())
	if err != nil {
		return nil, err
	}
	valBuf, err = AppendDriverCardHolderIdentification(valBuf, card.GetHolderIdentification())
	if err != nil {
		return nil, err
	}

	opts := cardv1.ElementaryFileType_EF_IDENTIFICATION.Descriptor().Values().ByNumber(protoreflect.EnumNumber(cardv1.ElementaryFileType_EF_IDENTIFICATION)).Options()
	tag := proto.GetExtension(opts, cardv1.E_FileId).(int32)
	dst = binary.BigEndian.AppendUint16(dst, uint16(tag))
	dst = binary.BigEndian.AppendUint16(dst, uint16(len(valBuf)))
	dst = append(dst, valBuf...)

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_DRIVING_LICENCE_INFO, card.GetDrivingLicenceInfo(), AppendDrivingLicenceInfo)
	if err != nil {
		return nil, err
	}

	// --- Special handling for EF_Events_Data ---
	eventsValBuf := make([]byte, 0, 1728) // Max size for Gen1
	eventsPerType := int(card.GetApplicationIdentification().GetEventsPerTypeCount())
	allEvents := card.GetEventsData().GetRecords()

	eventsByType := make(map[int32][]*cardv1.EventData_Record)
	for _, e := range allEvents {
		protocolValue := GetEventFaultTypeProtocolValue(e.GetEventType(), e.GetUnrecognizedEventType())
		eventsByType[protocolValue] = append(eventsByType[protocolValue], e)
	}

	// The 6 event groups in a Gen1 card file structure, ordered by type code.
	eventGroupTypeCodes := []int32{0x01, 0x02, 0x03, 0x04, 0x05, 0x07} // Example codes

	for _, eventTypeCode := range eventGroupTypeCodes {
		groupEvents := eventsByType[eventTypeCode]
		for j := 0; j < eventsPerType; j++ {
			if j < len(groupEvents) {
				eventsValBuf, err = AppendEventRecord(eventsValBuf, groupEvents[j])
				if err != nil {
					return nil, err
				}
			} else {
				// Pad with an empty 24-byte record
				eventsValBuf = append(eventsValBuf, make([]byte, 24)...)
			}
		}
	}

	eventsOpts := cardv1.ElementaryFileType_EF_EVENTS_DATA.Descriptor().Values().ByNumber(protoreflect.EnumNumber(cardv1.ElementaryFileType_EF_EVENTS_DATA)).Options()
	eventsTag := proto.GetExtension(eventsOpts, cardv1.E_FileId).(int32)
	dst = binary.BigEndian.AppendUint16(dst, uint16(eventsTag))
	dst = binary.BigEndian.AppendUint16(dst, uint16(len(eventsValBuf)))
	dst = append(dst, eventsValBuf...)

	// --- Special handling for EF_Faults_Data ---
	faultsValBuf := make([]byte, 0, 1152) // Max size for Gen1
	faultsPerType := int(card.GetApplicationIdentification().GetFaultsPerTypeCount())
	allFaults := card.GetFaultsData().GetRecords()

	faultsByType := make(map[bool][]*cardv1.FaultData_Record)
	for _, f := range allFaults {
		isEquipmentFault := (f.GetFaultType() >= 0x30 && f.GetFaultType() <= 0x3F)
		faultsByType[isEquipmentFault] = append(faultsByType[isEquipmentFault], f)
	}

	// Order: Equipment faults (true), then Card faults (false)
	faultGroupOrder := []bool{true, false}

	for _, isEquipmentFault := range faultGroupOrder {
		groupFaults := faultsByType[isEquipmentFault]
		for j := 0; j < faultsPerType; j++ {
			if j < len(groupFaults) {
				faultsValBuf, err = AppendFaultRecord(faultsValBuf, groupFaults[j])
				if err != nil {
					return nil, err
				}
			} else {
				// Pad with an empty 24-byte record
				faultsValBuf = append(faultsValBuf, make([]byte, 24)...)
			}
		}
	}

	faultsOpts := cardv1.ElementaryFileType_EF_FAULTS_DATA.Descriptor().Values().ByNumber(protoreflect.EnumNumber(cardv1.ElementaryFileType_EF_FAULTS_DATA)).Options()
	faultsTag := proto.GetExtension(faultsOpts, cardv1.E_FileId).(int32)
	dst = binary.BigEndian.AppendUint16(dst, uint16(faultsTag))
	dst = binary.BigEndian.AppendUint16(dst, uint16(len(faultsValBuf)))
	dst = append(dst, faultsValBuf...)

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_DRIVER_ACTIVITY_DATA, card.GetDriverActivityData(), AppendDriverActivity)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_VEHICLES_USED, card.GetVehiclesUsed(), AppendVehiclesUsed)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_PLACES, card.GetPlaces(), AppendPlaces)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_CURRENT_USAGE, card.GetCurrentUsage(), AppendCurrentUsage)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_CONTROL_ACTIVITY_DATA, card.GetControlActivityData(), AppendCardControlActivityData)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_SPECIFIC_CONDITIONS, card.GetSpecificConditions(), AppendCardSpecificConditions)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_CARD_DOWNLOAD_DRIVER, card.GetLastCardDownload(), AppendCardLastDownload)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION, card.GetApplicationIdentification(), AppendCardApplicationIdentification)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_VEHICLE_UNITS_USED, card.GetVehicleUnitsUsed(), AppendCardVehicleUnitsUsed)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_GNSS_PLACES, card.GetGnssPlaces(), AppendCardGnssPlaces)
	if err != nil {
		return nil, err
	}

	dst, err = appendTlv(dst, cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION_V2, card.GetApplicationIdentificationV2(), AppendCardApplicationIdentificationV2)
	if err != nil {
		return nil, err
	}

	// Append proprietary EFs if any exist for this card
	cardPtr := uintptr(unsafe.Pointer(card))
	if proprietaryEFs := GetProprietaryEFs(cardPtr); proprietaryEFs != nil {
		dst = proprietaryEFs.AppendProprietaryEFs(dst)
		// Clean up after marshalling
		ClearProprietaryEFs(cardPtr)
	}

	return dst, nil
}

// appendTlv is a generic helper for writing TLV records with zero value-allocation.
func appendTlv[T proto.Message](
	dst []byte,
	fileType cardv1.ElementaryFileType,
	msg T,
	appenderFunc func([]byte, T) ([]byte, error),
) ([]byte, error) {
	// Use reflection to check if the message is nil
	msgValue := reflect.ValueOf(msg)
	if !msgValue.IsValid() || (msgValue.Kind() == reflect.Ptr && msgValue.IsNil()) {
		return dst, nil // Don't write anything if the message is nil
	}

	opts := fileType.Descriptor().Values().ByNumber(protoreflect.EnumNumber(fileType)).Options()
	tag := proto.GetExtension(opts, cardv1.E_FileId).(int32)

	lenPos := len(dst)
	dst = append(dst, 0, 0, 0, 0, 0) // Placeholder for Tag (3 bytes) and Length (2 bytes)
	valPos := len(dst)

	var err error
	dst, err = appenderFunc(dst, msg)
	if err != nil {
		return nil, err
	}

	valLen := len(dst) - valPos

	// Write data tag (FID + appendix 0x00)
	binary.BigEndian.PutUint16(dst[lenPos:], uint16(tag))
	dst[lenPos+2] = 0x00 // appendix for data
	binary.BigEndian.PutUint16(dst[lenPos+3:], uint16(valLen))

	// Add signature block (FID + appendix 0x01, 128 bytes of zeros for now)
	dst = binary.BigEndian.AppendUint16(dst, uint16(tag))
	dst = append(dst, 0x01)                       // appendix for signature
	dst = binary.BigEndian.AppendUint16(dst, 128) // Signature length
	// Add 128 bytes of signature data (zeros for now)
	signature := make([]byte, 128)
	dst = append(dst, signature...)

	return dst, nil
}

// appendVU orchestrates writing a VU file in TV format
func appendVU(dst []byte, file *tachographv1.File) ([]byte, error) {
	vuFile := file.GetVehicleUnit()
	if vuFile == nil {
		return dst, nil
	}

	buf := bytes.NewBuffer(dst)

	// Process each transfer in the VU file
	for _, transfer := range vuFile.GetTransfers() {
		// Get the tag for this transfer type
		tag := getTrepValueForTransferType(transfer.GetType())
		if tag == 0 {
			continue // Skip unknown transfer types
		}

		// Append the 2-byte tag (0x76XX format)
		vuTag := uint16(0x7600 | (uint16(tag) & 0xFF))
		appendVuTag(buf, vuTag)

		// Append the transfer data based on type
		switch transfer.GetType() {
		case vuv1.TransferType_DOWNLOAD_INTERFACE_VERSION:
			AppendDownloadInterfaceVersion(buf, transfer.GetDownloadInterfaceVersion())
		case vuv1.TransferType_OVERVIEW_GEN1, vuv1.TransferType_OVERVIEW_GEN2_V1, vuv1.TransferType_OVERVIEW_GEN2_V2:
			AppendOverview(buf, transfer.GetOverview())
		case vuv1.TransferType_ACTIVITIES_GEN1, vuv1.TransferType_ACTIVITIES_GEN2_V1, vuv1.TransferType_ACTIVITIES_GEN2_V2:
			AppendVuActivities(buf, transfer.GetActivities())
		case vuv1.TransferType_EVENTS_AND_FAULTS_GEN1, vuv1.TransferType_EVENTS_AND_FAULTS_GEN2_V1:
			AppendVuEventsAndFaults(buf, transfer.GetEventsAndFaults())
		case vuv1.TransferType_DETAILED_SPEED_GEN1, vuv1.TransferType_DETAILED_SPEED_GEN2:
			AppendVuDetailedSpeed(buf, transfer.GetDetailedSpeed())
		case vuv1.TransferType_TECHNICAL_DATA_GEN1, vuv1.TransferType_TECHNICAL_DATA_GEN2_V1:
			AppendVuTechnicalData(buf, transfer.GetTechnicalData())
		default:
			// Skip unknown transfer types
		}
	}

	return buf.Bytes(), nil
}

// getTrepValueForTransferType returns the TREP value for a given transfer type
func getTrepValueForTransferType(transferType vuv1.TransferType) uint8 {
	values := vuv1.TransferType_TRANSFER_TYPE_UNSPECIFIED.Descriptor().Values()
	for i := 0; i < values.Len(); i++ {
		valueDesc := values.Get(i)
		if vuv1.TransferType(valueDesc.Number()) == transferType {
			opts := valueDesc.Options()
			if proto.HasExtension(opts, vuv1.E_TrepValue) {
				return uint8(proto.GetExtension(opts, vuv1.E_TrepValue).(int32))
			}
		}
	}
	return 0
}
