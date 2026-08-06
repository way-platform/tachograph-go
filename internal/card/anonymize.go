package card

import (
	"github.com/way-platform/tachograph-go/internal/dd"
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AnonymizeOptions configures the anonymization of card files.
type AnonymizeOptions struct {
	// PreserveDistanceAndTrips controls whether distance and trip data are preserved.
	PreserveDistanceAndTrips bool

	// PreserveTimestamps controls whether timestamps are preserved.
	PreserveTimestamps bool
}

// AnonymizeDriverCardFile creates an anonymized copy of a driver card file.
func (opts AnonymizeOptions) AnonymizeDriverCardFile(file *cardv1.DriverCardFile) (*cardv1.DriverCardFile, error) {
	if file == nil {
		return nil, nil
	}

	// Clone the file to avoid mutating the input
	result := proto.Clone(file).(*cardv1.DriverCardFile)

	// Anonymize common EFs (Master File)
	if icc := result.GetIcc(); icc != nil {
		result.SetIcc(opts.anonymizeIcc(icc))
	}
	if ic := result.GetIc(); ic != nil {
		result.SetIc(opts.anonymizeIc(ic))
	}

	// Anonymize Gen1 DF (Tachograph)
	if tachograph := result.GetTachograph(); tachograph != nil {
		// Clear certificates (TLV format: set to nil to omit blocks entirely)
		tachograph.SetCardCertificate(nil)
		tachograph.SetCaCertificate(nil)

		if appId := tachograph.GetApplicationIdentification(); appId != nil {
			tachograph.SetApplicationIdentification(opts.anonymizeApplicationIdentification(appId))
		}
		if identification := tachograph.GetIdentification(); identification != nil {
			tachograph.SetIdentification(opts.anonymizeDriverCardIdentification(identification))
		}
		if drivingLicenceInfo := tachograph.GetDrivingLicenceInfo(); drivingLicenceInfo != nil {
			tachograph.SetDrivingLicenceInfo(opts.anonymizeDrivingLicenceInfo(drivingLicenceInfo))
		}
		if eventsData := tachograph.GetEventsData(); eventsData != nil {
			tachograph.SetEventsData(opts.anonymizeEventsData(eventsData))
		}
		if faultsData := tachograph.GetFaultsData(); faultsData != nil {
			tachograph.SetFaultsData(opts.anonymizeFaultsData(faultsData))
		}
		if driverActivityData := tachograph.GetDriverActivityData(); driverActivityData != nil {
			tachograph.SetDriverActivityData(opts.anonymizeDriverActivityData(driverActivityData))
		}
		if vehiclesUsed := tachograph.GetVehiclesUsed(); vehiclesUsed != nil {
			tachograph.SetVehiclesUsed(opts.anonymizeVehiclesUsed(vehiclesUsed))
		}
		if places := tachograph.GetPlaces(); places != nil {
			tachograph.SetPlaces(opts.anonymizePlaces(places))
		}
		if currentUsage := tachograph.GetCurrentUsage(); currentUsage != nil {
			tachograph.SetCurrentUsage(opts.anonymizeCurrentUsage(currentUsage))
		}
		if controlActivityData := tachograph.GetControlActivityData(); controlActivityData != nil {
			tachograph.SetControlActivityData(opts.anonymizeControlActivityData(controlActivityData))
		}
		if specificConditions := tachograph.GetSpecificConditions(); specificConditions != nil {
			tachograph.SetSpecificConditions(opts.anonymizeSpecificConditions(specificConditions))
		}
	}

	// Anonymize Gen2 DF (Tachograph_G2)
	if tachographG2 := result.GetTachographG2(); tachographG2 != nil {
		// Clear certificates (TLV format: set to nil to omit blocks entirely)
		tachographG2.SetCardMaCertificate(nil)
		tachographG2.SetCardSignCertificate(nil)
		tachographG2.SetCaCertificate(nil)
		tachographG2.SetLinkCertificate(nil)

		// Anonymize Gen2 versions of shared EFs
		if appIdG2 := tachographG2.GetApplicationIdentification(); appIdG2 != nil {
			tachographG2.SetApplicationIdentification(opts.anonymizeApplicationIdentificationG2(appIdG2))
		}
		if identification := tachographG2.GetIdentification(); identification != nil {
			tachographG2.SetIdentification(opts.anonymizeDriverCardIdentification(identification))
		}
		if drivingLicenceInfo := tachographG2.GetDrivingLicenceInfo(); drivingLicenceInfo != nil {
			tachographG2.SetDrivingLicenceInfo(opts.anonymizeDrivingLicenceInfo(drivingLicenceInfo))
		}
		if eventsData := tachographG2.GetEventsData(); eventsData != nil {
			tachographG2.SetEventsData(opts.anonymizeEventsData(eventsData))
		}
		if faultsData := tachographG2.GetFaultsData(); faultsData != nil {
			tachographG2.SetFaultsData(opts.anonymizeFaultsData(faultsData))
		}
		if driverActivityData := tachographG2.GetDriverActivityData(); driverActivityData != nil {
			tachographG2.SetDriverActivityData(opts.anonymizeDriverActivityData(driverActivityData))
		}
		if vehiclesUsedG2 := tachographG2.GetVehiclesUsed(); vehiclesUsedG2 != nil {
			tachographG2.SetVehiclesUsed(opts.anonymizeVehiclesUsedG2(vehiclesUsedG2))
		}
		if placesG2 := tachographG2.GetPlaces(); placesG2 != nil {
			tachographG2.SetPlaces(opts.anonymizePlacesG2(placesG2))
		}
		if currentUsage := tachographG2.GetCurrentUsage(); currentUsage != nil {
			tachographG2.SetCurrentUsage(opts.anonymizeCurrentUsage(currentUsage))
		}
		if controlActivityData := tachographG2.GetControlActivityData(); controlActivityData != nil {
			tachographG2.SetControlActivityData(opts.anonymizeControlActivityData(controlActivityData))
		}
		if specificConditionsG2 := tachographG2.GetSpecificConditions(); specificConditionsG2 != nil {
			tachographG2.SetSpecificConditions(opts.anonymizeSpecificConditionsG2(specificConditionsG2))
		}

		// Anonymize Gen2-exclusive EFs
		if vehicleUnitsUsed := tachographG2.GetVehicleUnitsUsed(); vehicleUnitsUsed != nil {
			tachographG2.SetVehicleUnitsUsed(opts.anonymizeVehicleUnitsUsed(vehicleUnitsUsed))
		}
		if gnssPlaces := tachographG2.GetGnssPlaces(); gnssPlaces != nil {
			tachographG2.SetGnssPlaces(opts.anonymizeGnssPlaces(gnssPlaces))
		}

		if companyActivityData := tachographG2.GetCompanyActivityData(); companyActivityData != nil {
			tachographG2.SetCompanyActivityData(opts.anonymizeCompanyActivityData(companyActivityData))
		}
		if placesAuth := tachographG2.GetPlacesAuthentication(); placesAuth != nil {
			tachographG2.SetPlacesAuthentication(opts.anonymizePlacesAuthentication(placesAuth))
		}
		if gnssPlacesAuth := tachographG2.GetGnssPlacesAuthentication(); gnssPlacesAuth != nil {
			tachographG2.SetGnssPlacesAuthentication(opts.anonymizeGnssPlacesAuthentication(gnssPlacesAuth))
		}
		if loadTypeEntries := tachographG2.GetLoadTypeEntries(); loadTypeEntries != nil {
			tachographG2.SetLoadTypeEntries(opts.anonymizeLoadTypeEntries(loadTypeEntries))
		}
		if loadUnloadOps := tachographG2.GetLoadUnloadOperations(); loadUnloadOps != nil {
			tachographG2.SetLoadUnloadOperations(opts.anonymizeLoadUnloadOperations(loadUnloadOps))
		}
	}

	return result, nil
}

func (opts AnonymizeOptions) anonymizePlacesAuthentication(data *cardv1.PlacesAuthentication) *cardv1.PlacesAuthentication {
	if data == nil {
		return nil
	}
	result := &cardv1.PlacesAuthentication{}
	result.SetNewestRecordIndex(data.GetNewestRecordIndex())
	baseTimestamp := &timestamppb.Timestamp{Seconds: 1577836800}
	originalRecords := data.GetRecords()
	anonymizedRecords := make([]*cardv1.PlacesAuthentication_Record, len(originalRecords))
	for i, rec := range originalRecords {
		anon := &cardv1.PlacesAuthentication_Record{}
		if opts.PreserveTimestamps {
			anon.SetEntryTime(rec.GetEntryTime())
		} else {
			anon.SetEntryTime(baseTimestamp)
		}
		anon.SetAuthenticationStatus(rec.GetAuthenticationStatus())
		anonymizedRecords[i] = anon
	}
	result.SetRecords(anonymizedRecords)
	return result
}

func (opts AnonymizeOptions) anonymizeGnssPlacesAuthentication(data *cardv1.GnssPlacesAuthentication) *cardv1.GnssPlacesAuthentication {
	if data == nil {
		return nil
	}
	result := &cardv1.GnssPlacesAuthentication{}
	result.SetNewestRecordIndex(data.GetNewestRecordIndex())
	baseTimestamp := &timestamppb.Timestamp{Seconds: 1577836800}
	originalRecords := data.GetRecords()
	anonymizedRecords := make([]*cardv1.GnssPlacesAuthentication_Record, len(originalRecords))
	for i, rec := range originalRecords {
		anon := &cardv1.GnssPlacesAuthentication_Record{}
		if opts.PreserveTimestamps {
			anon.SetTimestamp(rec.GetTimestamp())
		} else {
			anon.SetTimestamp(baseTimestamp)
		}
		anon.SetAuthenticationStatus(rec.GetAuthenticationStatus())
		anonymizedRecords[i] = anon
	}
	result.SetRecords(anonymizedRecords)
	return result
}

func (opts AnonymizeOptions) anonymizeLoadTypeEntries(data *cardv1.LoadTypeEntries) *cardv1.LoadTypeEntries {
	if data == nil {
		return nil
	}
	result := &cardv1.LoadTypeEntries{}
	result.SetNewestRecordIndex(data.GetNewestRecordIndex())
	baseTimestamp := &timestamppb.Timestamp{Seconds: 1577836800}
	originalRecords := data.GetRecords()
	anonymizedRecords := make([]*cardv1.LoadTypeEntries_Record, len(originalRecords))
	for i, rec := range originalRecords {
		anon := &cardv1.LoadTypeEntries_Record{}
		if opts.PreserveTimestamps {
			anon.SetTimestamp(rec.GetTimestamp())
		} else {
			anon.SetTimestamp(baseTimestamp)
		}
		anon.SetLoadTypeEntered(rec.GetLoadTypeEntered())
		anonymizedRecords[i] = anon
	}
	result.SetRecords(anonymizedRecords)
	return result
}

func (opts AnonymizeOptions) anonymizeLoadUnloadOperations(data *cardv1.LoadUnloadOperations) *cardv1.LoadUnloadOperations {
	if data == nil {
		return nil
	}
	result := &cardv1.LoadUnloadOperations{}
	result.SetNewestRecordIndex(data.GetNewestRecordIndex())
	baseTimestamp := &timestamppb.Timestamp{Seconds: 1577836800}
	originalRecords := data.GetRecords()
	anonymizedRecords := make([]*cardv1.LoadUnloadOperations_Record, len(originalRecords))
	for i, rec := range originalRecords {
		anon := &cardv1.LoadUnloadOperations_Record{}
		if opts.PreserveTimestamps {
			anon.SetTimestamp(rec.GetTimestamp())
		} else {
			anon.SetTimestamp(baseTimestamp)
		}
		anon.SetOperationType(rec.GetOperationType())
		anon.SetGnssPlaceAuthRecord(rec.GetGnssPlaceAuthRecord())
		anon.SetVehicleOdometerKm(rec.GetVehicleOdometerKm())
		anonymizedRecords[i] = anon
	}
	result.SetRecords(anonymizedRecords)
	return result
}

// anonymizeCompanyActivityData creates an anonymized copy of CompanyActivityData,
// replacing card numbers, vehicle registrations, and timestamps with static test values.
func (opts AnonymizeOptions) anonymizeCompanyActivityData(data *cardv1.CompanyActivityData) *cardv1.CompanyActivityData {
	if data == nil {
		return nil
	}

	ddOpts := dd.AnonymizeOptions{
		PreserveDistanceAndTrips: opts.PreserveDistanceAndTrips,
		PreserveTimestamps:       opts.PreserveTimestamps,
	}

	result := &cardv1.CompanyActivityData{}
	result.SetNewestRecordIndex(data.GetNewestRecordIndex())

	// Static base timestamp: 2020-01-01 00:00:00 UTC
	baseTimestamp := &timestamppb.Timestamp{Seconds: 1577836800}

	originalRecords := data.GetRecords()
	anonymizedRecords := make([]*cardv1.CompanyActivityData_Record, len(originalRecords))
	for i, rec := range originalRecords {
		anon := &cardv1.CompanyActivityData_Record{}

		anon.SetCompanyActivityType(rec.GetCompanyActivityType())

		// Use static test timestamp: 2020-01-01 00:00:00 UTC (epoch: 1577836800)
		if opts.PreserveTimestamps {
			anon.SetCompanyActivityTime(rec.GetCompanyActivityTime())
			anon.SetDownloadPeriodBegin(rec.GetDownloadPeriodBegin())
			anon.SetDownloadPeriodEnd(rec.GetDownloadPeriodEnd())
		} else {
			anon.SetCompanyActivityTime(baseTimestamp)
			anon.SetDownloadPeriodBegin(baseTimestamp)
			anon.SetDownloadPeriodEnd(baseTimestamp)
		}

		// Anonymize card number and vehicle registration
		anon.SetCardNumberInformation(ddOpts.AnonymizeFullCardNumberAndGeneration(rec.GetCardNumberInformation()))
		anon.SetVehicleRegistrationInformation(ddOpts.AnonymizeVehicleRegistrationIdentification(rec.GetVehicleRegistrationInformation()))

		anonymizedRecords[i] = anon
	}
	result.SetRecords(anonymizedRecords)

	return result
}
