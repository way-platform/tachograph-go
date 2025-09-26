package tachograph

import (
	"fmt"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

func unmarshalDriverCardFile(input *cardv1.RawCardFile) (*cardv1.DriverCardFile, error) {
	var output cardv1.DriverCardFile
	for i := 0; i < len(input.GetRecords()); i++ {
		record := input.GetRecords()[i]
		if record.GetContentType() != cardv1.ContentType_DATA {
			return nil, fmt.Errorf("record %d has unexpected content type", i)
		}
		if !record.HasFile() {
			return nil, fmt.Errorf("record %d has no file type", i)
		}
		var signature []byte
		if i+1 < len(input.GetRecords()) {
			nextRecord := input.GetRecords()[i+1]
			if nextRecord.GetFile() == record.GetFile() && nextRecord.GetContentType() == cardv1.ContentType_SIGNATURE {
				signature = nextRecord.GetValue()
				i++
			}
		}
		switch record.GetFile() {
		case cardv1.ElementaryFileType_EF_ICC:
			icc, err := unmarshalIcc(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				return nil, fmt.Errorf("unexpected signature for EF_ICC")
			}
			output.SetIcc(icc)

		case cardv1.ElementaryFileType_EF_IC:
			ic, err := unmarshalCardIc(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				return nil, fmt.Errorf("unexpected signature for EF_IC")
			}
			output.SetIc(ic)

		case cardv1.ElementaryFileType_EF_IDENTIFICATION:
			identification, err := unmarshalIdentification(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				identification.SetSignature(signature)
			}
			output.SetIdentification(identification)

		case cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION:
			appId, err := unmarshalCardApplicationIdentification(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				appId.SetSignature(signature)
			}
			output.SetApplicationIdentification(appId)

		case cardv1.ElementaryFileType_EF_DRIVING_LICENCE_INFO:
			drivingLicenceInfo, err := unmarshalDrivingLicenceInfo(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				drivingLicenceInfo.SetSignature(signature)
			}
			output.SetDrivingLicenceInfo(drivingLicenceInfo)

		case cardv1.ElementaryFileType_EF_EVENTS_DATA:
			eventsData, err := unmarshalEventsData(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				eventsData.SetSignature(signature)
			}
			output.SetEventsData(eventsData)

		case cardv1.ElementaryFileType_EF_FAULTS_DATA:
			faultsData, err := unmarshalFaultsData(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				faultsData.SetSignature(signature)
			}
			output.SetFaultsData(faultsData)

		case cardv1.ElementaryFileType_EF_DRIVER_ACTIVITY_DATA:
			activityData, err := unmarshalDriverActivityData(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				activityData.SetSignature(signature)
			}
			output.SetDriverActivityData(activityData)

		case cardv1.ElementaryFileType_EF_VEHICLES_USED:
			vehiclesUsed, err := unmarshalCardVehiclesUsed(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				vehiclesUsed.SetSignature(signature)
			}
			output.SetVehiclesUsed(vehiclesUsed)

		case cardv1.ElementaryFileType_EF_PLACES:
			places, err := unmarshalCardPlaces(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				places.SetSignature(signature)
			}
			output.SetPlaces(places)

		case cardv1.ElementaryFileType_EF_CURRENT_USAGE:
			currentUsage, err := unmarshalCardCurrentUsage(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				currentUsage.SetSignature(signature)
			}
			output.SetCurrentUsage(currentUsage)

		case cardv1.ElementaryFileType_EF_CONTROL_ACTIVITY_DATA:
			controlActivity, err := unmarshalCardControlActivityData(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				controlActivity.SetSignature(signature)
			}
			output.SetControlActivityData(controlActivity)

		case cardv1.ElementaryFileType_EF_SPECIFIC_CONDITIONS:
			specificConditions, err := unmarshalCardSpecificConditions(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				specificConditions.SetSignature(signature)
			}
			output.SetSpecificConditions(specificConditions)

		case cardv1.ElementaryFileType_EF_CARD_DOWNLOAD_DRIVER:
			lastDownload, err := unmarshalCardLastDownload(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				lastDownload.SetSignature(signature)
			}
			output.SetLastCardDownload(lastDownload)

		case cardv1.ElementaryFileType_EF_VEHICLE_UNITS_USED:
			vehicleUnits, err := unmarshalCardVehicleUnitsUsed(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				vehicleUnits.SetSignature(signature)
			}
			output.SetVehicleUnitsUsed(vehicleUnits)

		case cardv1.ElementaryFileType_EF_GNSS_PLACES:
			gnssPlaces, err := unmarshalCardGnssPlaces(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				gnssPlaces.SetSignature(signature)
			}
			output.SetGnssPlaces(gnssPlaces)

		case cardv1.ElementaryFileType_EF_APPLICATION_IDENTIFICATION_V2:
			appIdV2, err := unmarshalCardApplicationIdentificationV2(record.GetValue())
			if err != nil {
				return nil, err
			}
			if signature != nil {
				appIdV2.SetSignature(signature)
			}
			output.SetApplicationIdentificationV2(appIdV2)

		case cardv1.ElementaryFileType_EF_CARD_CERTIFICATE:
			if output.GetCertificates() == nil {
				output.SetCertificates(&cardv1.Certificates{})
			}
			output.GetCertificates().SetCardCertificate(record.GetValue())
			if signature != nil {
				return nil, fmt.Errorf("unexpected signature for EF_CARD_CERTIFICATE")
			}

		case cardv1.ElementaryFileType_EF_CA_CERTIFICATE:
			if output.GetCertificates() == nil {
				output.SetCertificates(&cardv1.Certificates{})
			}
			output.GetCertificates().SetCaCertificate(record.GetValue())
			if signature != nil {
				return nil, fmt.Errorf("unexpected signature for EF_CA_CERTIFICATE")
			}
		}
	}
	return &output, nil
}
