package tachounit

// IsValid checks if the DownloadType is a known value.
func (dt DownloadType) IsValid() bool {
	switch dt {
	case
		TRTP_DownloadInterfaceVersion,
		TRTP_Overview,
		TRTP_ActivitiesOfASpecifiedDate,
		TRTP_EventsAndFaults,
		TRTP_DetailedSpeed,
		TRTP_TechnicalData,
		TRTP_Overview_0x15,
		TRTP_ActivitiesOfASpecifiedDate_0x16,
		TRTP_EventsAndFaults_0x17,
		TRTP_DetailedSpeed_0x18,
		TRTP_TechnicalData_0x19,
		TRTP_Overview_0x1F,
		TRTP_ActivitiesOfASpecifiedDate_0x20,
		TRTP_EventsAndFaults_0x21,
		TRTP_TechnicalData_0x23:
		return true
	}
	return false
}
