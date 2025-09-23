package tachograph

import (
	"bytes"
	"encoding/binary"
	"errors"

	"google.golang.org/protobuf/proto"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
)

// Unmarshal parses a .DDD file's byte data into a protobuf File message.
func Unmarshal(data []byte) (*tachographv1.File, error) {
	fileType := InferFileType(data)
	switch fileType {
	case CardFileType:
		return unmarshalCard(data)
	case UnitFileType:
		return nil, errors.New("vehicle unit unmarshaling not yet implemented")
	}
	return nil, errors.New("unknown or unsupported file type")
}

func unmarshalCard(data []byte) (*tachographv1.File, error) {
	file := &tachographv1.File{}
	file.SetType(tachographv1.File_DRIVER_CARD) // Assume Driver card for now
	file.SetDriverCard(&cardv1.DriverCardFile{})

	r := bytes.NewReader(data)

	for r.Len() > 0 {
		// Read Tag - 3 bytes (FID + appendix) according to DDD format spec
		tagBytes := make([]byte, 3)
		if _, err := r.Read(tagBytes); err != nil {
			return nil, err
		}
		// Extract FID (first 2 bytes) and appendix (last byte)
		fid := binary.BigEndian.Uint16(tagBytes[0:2])
		appendix := tagBytes[2]
		
		// Read Length - 2 bytes
		var length uint16
		if err := binary.Read(r, binary.BigEndian, &length); err != nil {
			return nil, err
		}

		value := make([]byte, length)
		if _, err := r.Read(value); err != nil {
			return nil, err
		}

		// Skip signatures (appendix '01' and '03') - we only process data (appendix '00' and '02')
		if appendix == 0x01 || appendix == 0x03 {
			continue // Skip signature TLV objects
		}
		
		// Find file type from FID and dispatch
		fileType := findFileTypeByTag(int32(fid))
		if fileType == cardv1.ElementaryFileType_ELEMENTARY_FILE_UNSPECIFIED {
			continue // Skip unknown tags
		}

		driverCard := file.GetDriverCard()
		switch fileType {
		case cardv1.ElementaryFileType_EF_ICC:
			icc := &cardv1.IccIdentification{}
			if err := UnmarshalIcc(value, icc); err != nil {
				return nil, err
			}
			driverCard.SetIcc(icc)
		case cardv1.ElementaryFileType_EF_IDENTIFICATION:
			identification := &cardv1.CardIdentification{}
			holderIdentification := &cardv1.DriverCardHolderIdentification{}
			if err := UnmarshalIdentification(value, identification, holderIdentification); err != nil {
				return nil, err
			}
			driverCard.SetIdentification(identification)
			driverCard.SetHolderIdentification(holderIdentification)
		case cardv1.ElementaryFileType_EF_DRIVING_LICENCE_INFO:
			drivingLicenceInfo := &cardv1.DrivingLicenceInfo{}
			if err := UnmarshalDrivingLicenceInfo(value, drivingLicenceInfo); err != nil {
				return nil, err
			}
			driverCard.SetDrivingLicenceInfo(drivingLicenceInfo)
		case cardv1.ElementaryFileType_EF_EVENTS_DATA:
			eventsData := &cardv1.EventData{}
			if err := UnmarshalEventsData(value, eventsData); err != nil {
				return nil, err
			}
			driverCard.SetEventsData(eventsData)
		case cardv1.ElementaryFileType_EF_FAULTS_DATA:
			faultsData := &cardv1.FaultData{}
			if err := UnmarshalFaultsData(value, faultsData); err != nil {
				return nil, err
			}
			driverCard.SetFaultsData(faultsData)
			// ... other cases to be added
		}
	}

	return file, nil
}

func findFileTypeByTag(tag int32) cardv1.ElementaryFileType {
	values := cardv1.ElementaryFileType_ELEMENTARY_FILE_UNSPECIFIED.Descriptor().Values()
	for i := 0; i < values.Len(); i++ {
		valueDesc := values.Get(i)
		opts := valueDesc.Options()
		if proto.HasExtension(opts, cardv1.E_FileId) {
			if proto.GetExtension(opts, cardv1.E_FileId).(int32) == tag {
				return cardv1.ElementaryFileType(valueDesc.Number())
			}
		}
	}
	return cardv1.ElementaryFileType_ELEMENTARY_FILE_UNSPECIFIED
}
