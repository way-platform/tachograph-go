// Code generated based on regulation Appendix 7. DO NOT EDIT.

package tachounit

// VuTag represents a vehicle unit TV tag identifier.
type VuTag uint16

// Defines the known vehicle unit TV tags based on SID 76 + TREP combinations.
const (
	// VU_DownloadInterfaceVersion is the VU Download Interface Version (SID 0x76, TREP 0x00)
	VU_DownloadInterfaceVersion VuTag = 0x7600
	// VU_OverviewFirstGen is the VU Overview First Generation (SID 0x76, TREP 0x01)
	VU_OverviewFirstGen VuTag = 0x7601
	// VU_ActivitiesFirstGen is the VU Activities First Generation (SID 0x76, TREP 0x02)
	VU_ActivitiesFirstGen VuTag = 0x7602
	// VU_EventsAndFaultsFirstGen is the VU Events and Faults First Generation (SID 0x76, TREP 0x03)
	VU_EventsAndFaultsFirstGen VuTag = 0x7603
	// VU_DetailedSpeedFirstGen is the VU Detailed Speed First Generation (SID 0x76, TREP 0x04)
	VU_DetailedSpeedFirstGen VuTag = 0x7604
	// VU_TechnicalDataFirstGen is the VU Technical Data First Generation (SID 0x76, TREP 0x05)
	VU_TechnicalDataFirstGen VuTag = 0x7605
	// VU_OverviewSecondGen is the VU Overview Second Generation (SID 0x76, TREP 0x21)
	VU_OverviewSecondGen VuTag = 0x7621
	// VU_ActivitiesSecondGen is the VU Activities Second Generation (SID 0x76, TREP 0x22)
	VU_ActivitiesSecondGen VuTag = 0x7622
	// VU_EventsAndFaultsSecondGen is the VU Events and Faults Second Generation (SID 0x76, TREP 0x23)
	VU_EventsAndFaultsSecondGen VuTag = 0x7623
	// VU_DetailedSpeedSecondGen is the VU Detailed Speed Second Generation (SID 0x76, TREP 0x24)
	VU_DetailedSpeedSecondGen VuTag = 0x7624
	// VU_TechnicalDataSecondGen is the VU Technical Data Second Generation (SID 0x76, TREP 0x25)
	VU_TechnicalDataSecondGen VuTag = 0x7625
	// VU_OverviewSecondGenV2 is the VU Overview Second Generation V2 (SID 0x76, TREP 0x31)
	VU_OverviewSecondGenV2 VuTag = 0x7631
	// VU_ActivitiesSecondGenV2 is the VU Activities Second Generation V2 (SID 0x76, TREP 0x32)
	VU_ActivitiesSecondGenV2 VuTag = 0x7632
	// VU_EventsAndFaultsSecondGenV2 is the VU Events and Faults Second Generation V2 (SID 0x76, TREP 0x33)
	VU_EventsAndFaultsSecondGenV2 VuTag = 0x7633
	// VU_TechnicalDataSecondGenV2 is the VU Technical Data Second Generation V2 (SID 0x76, TREP 0x35)
	VU_TechnicalDataSecondGenV2 VuTag = 0x7635
)

// IsValid checks if the VuTag is a known valid tag.
func (vt VuTag) IsValid() bool {
	switch vt {
	case VU_DownloadInterfaceVersion,
		VU_OverviewFirstGen,
		VU_ActivitiesFirstGen,
		VU_EventsAndFaultsFirstGen,
		VU_DetailedSpeedFirstGen,
		VU_TechnicalDataFirstGen,
		VU_OverviewSecondGen,
		VU_ActivitiesSecondGen,
		VU_EventsAndFaultsSecondGen,
		VU_DetailedSpeedSecondGen,
		VU_TechnicalDataSecondGen,
		VU_OverviewSecondGenV2,
		VU_ActivitiesSecondGenV2,
		VU_EventsAndFaultsSecondGenV2,
		VU_TechnicalDataSecondGenV2:
		return true
	}
	return false
}
