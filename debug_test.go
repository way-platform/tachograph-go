package tachograph

import (
	"os"
	"testing"
)

func TestDebugUnmarshal(t *testing.T) {
	filePath := "testdata/card/proprietary-Nuutti_Nestori_Sahala_2025-09-11_08-38-24.DDD"

	originalData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	t.Logf("Original file size: %d bytes", len(originalData))
	t.Logf("First 32 bytes: %X", originalData[:32])

	file, err := UnmarshalFile(originalData)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	t.Logf("File type: %v", file.GetType())

	if driverCard := file.GetDriverCard(); driverCard != nil {
		t.Logf("Driver card fields present:")
		if driverCard.GetIcc() != nil {
			t.Logf("  - ICC: present")
		}
		if driverCard.GetIc() != nil {
			t.Logf("  - IC: present")
		}
		if driverCard.GetApplicationIdentification() != nil {
			t.Logf("  - Application Identification: present")
		}
		if driverCard.GetIdentification() != nil {
			t.Logf("  - Identification: present")
		}
		if driverCard.GetHolderIdentification() != nil {
			t.Logf("  - Holder Identification: present")
		}
		if driverCard.GetDrivingLicenceInfo() != nil {
			t.Logf("  - Driving Licence Info: present")
		}
		if driverCard.GetEventsData() != nil {
			t.Logf("  - Events Data: present (%d records)", len(driverCard.GetEventsData().GetRecords()))
		}
		if driverCard.GetFaultsData() != nil {
			t.Logf("  - Faults Data: present (%d records)", len(driverCard.GetFaultsData().GetRecords()))
		}
		if driverCard.GetDriverActivityData() != nil {
			t.Logf("  - Driver Activity Data: present")
		}
		if driverCard.GetVehiclesUsed() != nil {
			t.Logf("  - Vehicles Used: present")
		}
		if driverCard.GetPlaces() != nil {
			t.Logf("  - Places: present")
		}
		if driverCard.GetCurrentUsage() != nil {
			t.Logf("  - Current Usage: present")
		}
		if driverCard.GetControlActivityData() != nil {
			t.Logf("  - Control Activity Data: present")
		}
		if driverCard.GetSpecificConditions() != nil {
			t.Logf("  - Specific Conditions: present")
		}
		if driverCard.GetLastCardDownload() != nil {
			t.Logf("  - Last Card Download: present")
		}
		if driverCard.GetVehicleUnitsUsed() != nil {
			t.Logf("  - Vehicle Units Used: present")
		}
		if driverCard.GetGnssPlaces() != nil {
			t.Logf("  - GNSS Places: present")
		}
		if driverCard.GetApplicationIdentificationV2() != nil {
			t.Logf("  - Application Identification V2: present")
		}
		if driverCard.GetCertificates() != nil {
			certs := driverCard.GetCertificates()
			t.Logf("  - Certificates: present")
			if len(certs.GetCardCertificate()) > 0 {
				t.Logf("    - Card Certificate: %d bytes", len(certs.GetCardCertificate()))
			}
			if len(certs.GetCaCertificate()) > 0 {
				t.Logf("    - CA Certificate: %d bytes", len(certs.GetCaCertificate()))
			}
		}
	}

	marshalledData, err := Marshal(file)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	t.Logf("Marshalled file size: %d bytes", len(marshalledData))
	t.Logf("First 32 bytes: %X", marshalledData[:min(32, len(marshalledData))])

	t.Logf("Size difference: %d bytes", len(originalData)-len(marshalledData))
}
